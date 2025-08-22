package gonest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// TestApp represents a test application
type TestApp struct {
	app        *Application
	server     *httptest.Server
	httpClient *http.Client
	logger     *logrus.Logger
}

// NewTestApp creates a new test application
func NewTestApp(t *testing.T) *TestApp {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	config := DefaultConfig()
	config.Port = "0" // Use random port for testing

	app := NewApplication().
		Config(config).
		Logger(logger).
		Build()

	testApp := &TestApp{
		app:        app,
		httpClient: &http.Client{},
		logger:     logger,
	}

	return testApp
}

// WithModule adds a module to the test app
func (ta *TestApp) WithModule(module *Module) *TestApp {
	ta.app.Module(module)
	return ta
}

// WithConfig sets configuration for the test app
func (ta *TestApp) WithConfig(config *Config) *TestApp {
	ta.app.Config = config
	return ta
}

// Start starts the test application
func (ta *TestApp) Start(t *testing.T) *TestApp {
	// Initialize the application
	if err := ta.app.Initialize(); err != nil {
		t.Fatalf("Failed to initialize test app: %v", err)
	}

	// Create test server
	ta.server = httptest.NewServer(ta.app.Echo)

	return ta
}

// Stop stops the test application
func (ta *TestApp) Stop() {
	if ta.server != nil {
		ta.server.Close()
	}
	if ta.app != nil {
		ta.app.Stop()
	}
}

// GetBaseURL returns the base URL of the test server
func (ta *TestApp) GetBaseURL() string {
	if ta.server == nil {
		return ""
	}
	return ta.server.URL
}

// GetApp returns the underlying application
func (ta *TestApp) GetApp() *Application {
	return ta.app
}

// GetEcho returns the underlying Echo instance
func (ta *TestApp) GetEcho() *echo.Echo {
	return ta.app.Echo
}

// TestRequest represents a test HTTP request
type TestRequest struct {
	testApp *TestApp
	method  string
	path    string
	headers map[string]string
	body    interface{}
	query   map[string]string
	cookies []*http.Cookie
}

// NewTestRequest creates a new test request
func (ta *TestApp) NewRequest(method, path string) *TestRequest {
	return &TestRequest{
		testApp: ta,
		method:  method,
		path:    path,
		headers: make(map[string]string),
		query:   make(map[string]string),
		cookies: make([]*http.Cookie, 0),
	}
}

// GET creates a GET request
func (ta *TestApp) GET(path string) *TestRequest {
	return ta.NewRequest("GET", path)
}

// POST creates a POST request
func (ta *TestApp) POST(path string) *TestRequest {
	return ta.NewRequest("POST", path)
}

// PUT creates a PUT request
func (ta *TestApp) PUT(path string) *TestRequest {
	return ta.NewRequest("PUT", path)
}

// DELETE creates a DELETE request
func (ta *TestApp) DELETE(path string) *TestRequest {
	return ta.NewRequest("DELETE", path)
}

// PATCH creates a PATCH request
func (ta *TestApp) PATCH(path string) *TestRequest {
	return ta.NewRequest("PATCH", path)
}

// WithHeader adds a header to the request
func (tr *TestRequest) WithHeader(key, value string) *TestRequest {
	tr.headers[key] = value
	return tr
}

// WithHeaders adds multiple headers to the request
func (tr *TestRequest) WithHeaders(headers map[string]string) *TestRequest {
	for key, value := range headers {
		tr.headers[key] = value
	}
	return tr
}

// WithAuth sets Authorization header
func (tr *TestRequest) WithAuth(token string) *TestRequest {
	tr.headers["Authorization"] = "Bearer " + token
	return tr
}

// WithJSON sets JSON body and Content-Type header
func (tr *TestRequest) WithJSON(body interface{}) *TestRequest {
	tr.body = body
	tr.headers["Content-Type"] = "application/json"
	return tr
}

// WithBody sets request body
func (tr *TestRequest) WithBody(body interface{}) *TestRequest {
	tr.body = body
	return tr
}

// WithQuery adds query parameters
func (tr *TestRequest) WithQuery(key, value string) *TestRequest {
	tr.query[key] = value
	return tr
}

// WithQueryParams adds multiple query parameters
func (tr *TestRequest) WithQueryParams(params map[string]string) *TestRequest {
	for key, value := range params {
		tr.query[key] = value
	}
	return tr
}

// WithCookie adds a cookie to the request
func (tr *TestRequest) WithCookie(cookie *http.Cookie) *TestRequest {
	tr.cookies = append(tr.cookies, cookie)
	return tr
}

// Send sends the test request and returns response
func (tr *TestRequest) Send(t *testing.T) *TestResponse {
	// Build URL with query parameters
	requestURL := tr.testApp.GetBaseURL() + tr.path
	if len(tr.query) > 0 {
		params := url.Values{}
		for key, value := range tr.query {
			params.Add(key, value)
		}
		requestURL += "?" + params.Encode()
	}

	// Prepare body
	var bodyReader io.Reader
	if tr.body != nil {
		switch body := tr.body.(type) {
		case string:
			bodyReader = strings.NewReader(body)
		case []byte:
			bodyReader = bytes.NewReader(body)
		default:
			// JSON marshal
			jsonData, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}
			bodyReader = bytes.NewReader(jsonData)
		}
	}

	// Create request
	req, err := http.NewRequest(tr.method, requestURL, bodyReader)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Add headers
	for key, value := range tr.headers {
		req.Header.Set(key, value)
	}

	// Add cookies
	for _, cookie := range tr.cookies {
		req.AddCookie(cookie)
	}

	// Send request
	resp, err := tr.testApp.httpClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	return &TestResponse{
		response: resp,
		t:        t,
	}
}

// TestResponse represents a test HTTP response
type TestResponse struct {
	response *http.Response
	t        *testing.T
	bodyRead bool
	bodyData []byte
}

// Status returns the response status code
func (tr *TestResponse) Status() int {
	return tr.response.StatusCode
}

// Header returns a response header
func (tr *TestResponse) Header(key string) string {
	return tr.response.Header.Get(key)
}

// Headers returns all response headers
func (tr *TestResponse) Headers() http.Header {
	return tr.response.Header
}

// Body returns the response body as bytes
func (tr *TestResponse) Body() []byte {
	if !tr.bodyRead {
		defer tr.response.Body.Close()
		data, err := io.ReadAll(tr.response.Body)
		if err != nil {
			tr.t.Fatalf("Failed to read response body: %v", err)
		}
		tr.bodyData = data
		tr.bodyRead = true
	}
	return tr.bodyData
}

// BodyString returns the response body as string
func (tr *TestResponse) BodyString() string {
	return string(tr.Body())
}

// JSON unmarshals response body as JSON
func (tr *TestResponse) JSON(dest interface{}) *TestResponse {
	if err := json.Unmarshal(tr.Body(), dest); err != nil {
		tr.t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}
	return tr
}

// Cookie returns a cookie from the response
func (tr *TestResponse) Cookie(name string) *http.Cookie {
	for _, cookie := range tr.response.Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

// Cookies returns all cookies from the response
func (tr *TestResponse) Cookies() []*http.Cookie {
	return tr.response.Cookies()
}

// Assertion methods

// ExpectStatus asserts the response status code
func (tr *TestResponse) ExpectStatus(expectedStatus int) *TestResponse {
	if tr.Status() != expectedStatus {
		tr.t.Errorf("Expected status %d, got %d. Body: %s", expectedStatus, tr.Status(), tr.BodyString())
	}
	return tr
}

// ExpectOK asserts status 200
func (tr *TestResponse) ExpectOK() *TestResponse {
	return tr.ExpectStatus(http.StatusOK)
}

// ExpectCreated asserts status 201
func (tr *TestResponse) ExpectCreated() *TestResponse {
	return tr.ExpectStatus(http.StatusCreated)
}

// ExpectBadRequest asserts status 400
func (tr *TestResponse) ExpectBadRequest() *TestResponse {
	return tr.ExpectStatus(http.StatusBadRequest)
}

// ExpectUnauthorized asserts status 401
func (tr *TestResponse) ExpectUnauthorized() *TestResponse {
	return tr.ExpectStatus(http.StatusUnauthorized)
}

// ExpectForbidden asserts status 403
func (tr *TestResponse) ExpectForbidden() *TestResponse {
	return tr.ExpectStatus(http.StatusForbidden)
}

// ExpectNotFound asserts status 404
func (tr *TestResponse) ExpectNotFound() *TestResponse {
	return tr.ExpectStatus(http.StatusNotFound)
}

// ExpectInternalServerError asserts status 500
func (tr *TestResponse) ExpectInternalServerError() *TestResponse {
	return tr.ExpectStatus(http.StatusInternalServerError)
}

// ExpectHeader asserts a response header
func (tr *TestResponse) ExpectHeader(key, expectedValue string) *TestResponse {
	actualValue := tr.Header(key)
	if actualValue != expectedValue {
		tr.t.Errorf("Expected header %s to be '%s', got '%s'", key, expectedValue, actualValue)
	}
	return tr
}

// ExpectHeaderContains asserts a response header contains a value
func (tr *TestResponse) ExpectHeaderContains(key, expectedSubstring string) *TestResponse {
	actualValue := tr.Header(key)
	if !strings.Contains(actualValue, expectedSubstring) {
		tr.t.Errorf("Expected header %s to contain '%s', got '%s'", key, expectedSubstring, actualValue)
	}
	return tr
}

// ExpectBodyContains asserts the response body contains a string
func (tr *TestResponse) ExpectBodyContains(expectedSubstring string) *TestResponse {
	body := tr.BodyString()
	if !strings.Contains(body, expectedSubstring) {
		tr.t.Errorf("Expected body to contain '%s', got: %s", expectedSubstring, body)
	}
	return tr
}

// ExpectJSONField asserts a JSON field value
func (tr *TestResponse) ExpectJSONField(field string, expectedValue interface{}) *TestResponse {
	var data map[string]interface{}
	tr.JSON(&data)

	actualValue := getNestedField(data, field)
	if !reflect.DeepEqual(actualValue, expectedValue) {
		tr.t.Errorf("Expected JSON field %s to be %v, got %v", field, expectedValue, actualValue)
	}
	return tr
}

// ExpectJSONFieldExists asserts a JSON field exists
func (tr *TestResponse) ExpectJSONFieldExists(field string) *TestResponse {
	var data map[string]interface{}
	tr.JSON(&data)

	if getNestedField(data, field) == nil {
		tr.t.Errorf("Expected JSON field %s to exist", field)
	}
	return tr
}

// ExpectJSONArray asserts the response is a JSON array
func (tr *TestResponse) ExpectJSONArray() *TestResponse {
	var data []interface{}
	if err := json.Unmarshal(tr.Body(), &data); err != nil {
		tr.t.Errorf("Expected JSON array, but failed to unmarshal: %v", err)
	}
	return tr
}

// ExpectJSONArrayLength asserts the JSON array length
func (tr *TestResponse) ExpectJSONArrayLength(expectedLength int) *TestResponse {
	var data []interface{}
	tr.JSON(&data)

	if len(data) != expectedLength {
		tr.t.Errorf("Expected JSON array length %d, got %d", expectedLength, len(data))
	}
	return tr
}

// Mock utilities

// MockService interface for creating service mocks
type MockService interface {
	Reset()
}

// ServiceMocker helps create service mocks for testing
type ServiceMocker struct {
	mocks map[string]interface{}
}

// NewServiceMocker creates a new service mocker
func NewServiceMocker() *ServiceMocker {
	return &ServiceMocker{
		mocks: make(map[string]interface{}),
	}
}

// Mock registers a mock service
func (sm *ServiceMocker) Mock(name string, mock interface{}) {
	sm.mocks[name] = mock
}

// GetMock retrieves a mock service
func (sm *ServiceMocker) GetMock(name string) interface{} {
	return sm.mocks[name]
}

// ResetAll resets all mocks
func (sm *ServiceMocker) ResetAll() {
	for _, mock := range sm.mocks {
		if resetable, ok := mock.(MockService); ok {
			resetable.Reset()
		}
	}
}

// Test fixtures

// TestFixture interface for test data setup
type TestFixture interface {
	Setup(t *testing.T) error
	Teardown(t *testing.T) error
}

// DatabaseFixture provides database test data setup
type DatabaseFixture struct {
	setupFunc    func(*testing.T) error
	teardownFunc func(*testing.T) error
}

// NewDatabaseFixture creates a new database fixture
func NewDatabaseFixture(setup, teardown func(*testing.T) error) *DatabaseFixture {
	return &DatabaseFixture{
		setupFunc:    setup,
		teardownFunc: teardown,
	}
}

// Setup sets up test data
func (df *DatabaseFixture) Setup(t *testing.T) error {
	if df.setupFunc != nil {
		return df.setupFunc(t)
	}
	return nil
}

// Teardown cleans up test data
func (df *DatabaseFixture) Teardown(t *testing.T) error {
	if df.teardownFunc != nil {
		return df.teardownFunc(t)
	}
	return nil
}

// Test runner with fixtures

// TestRunner manages test execution with fixtures
type TestRunner struct {
	fixtures []TestFixture
	mocker   *ServiceMocker
}

// NewTestRunner creates a new test runner
func NewTestRunner() *TestRunner {
	return &TestRunner{
		fixtures: make([]TestFixture, 0),
		mocker:   NewServiceMocker(),
	}
}

// WithFixture adds a fixture to the test runner
func (tr *TestRunner) WithFixture(fixture TestFixture) *TestRunner {
	tr.fixtures = append(tr.fixtures, fixture)
	return tr
}

// WithMock adds a mock to the test runner
func (tr *TestRunner) WithMock(name string, mock interface{}) *TestRunner {
	tr.mocker.Mock(name, mock)
	return tr
}

// Run runs a test with all fixtures and mocks
func (tr *TestRunner) Run(t *testing.T, testFunc func(*testing.T)) {
	// Setup fixtures
	for _, fixture := range tr.fixtures {
		if err := fixture.Setup(t); err != nil {
			t.Fatalf("Failed to setup fixture: %v", err)
		}
	}

	// Defer teardown
	defer func() {
		for i := len(tr.fixtures) - 1; i >= 0; i-- {
			if err := tr.fixtures[i].Teardown(t); err != nil {
				t.Errorf("Failed to teardown fixture: %v", err)
			}
		}
		tr.mocker.ResetAll()
	}()

	// Run test
	testFunc(t)
}

// Utility functions

// getNestedField retrieves a nested field from a map
func getNestedField(data map[string]interface{}, field string) interface{} {
	parts := strings.Split(field, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			return current[part]
		}

		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}

// TestHelper provides additional testing utilities
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// AssertEqual asserts two values are equal
func (th *TestHelper) AssertEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		msg := fmt.Sprintf("Expected %v, got %v", expected, actual)
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			}
		}
		th.t.Error(msg)
	}
}

// AssertNotEqual asserts two values are not equal
func (th *TestHelper) AssertNotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	if reflect.DeepEqual(expected, actual) {
		msg := fmt.Sprintf("Expected values to be different, but both were %v", expected)
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			}
		}
		th.t.Error(msg)
	}
}

// AssertNil asserts a value is nil
func (th *TestHelper) AssertNil(value interface{}, msgAndArgs ...interface{}) {
	if value != nil {
		msg := fmt.Sprintf("Expected nil, got %v", value)
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			}
		}
		th.t.Error(msg)
	}
}

// AssertNotNil asserts a value is not nil
func (th *TestHelper) AssertNotNil(value interface{}, msgAndArgs ...interface{}) {
	if value == nil {
		msg := "Expected non-nil value, got nil"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			}
		}
		th.t.Error(msg)
	}
}

// AssertTrue asserts a value is true
func (th *TestHelper) AssertTrue(value bool, msgAndArgs ...interface{}) {
	if !value {
		msg := "Expected true, got false"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			}
		}
		th.t.Error(msg)
	}
}

// AssertFalse asserts a value is false
func (th *TestHelper) AssertFalse(value bool, msgAndArgs ...interface{}) {
	if value {
		msg := "Expected false, got true"
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			}
		}
		th.t.Error(msg)
	}
}

// AssertContains asserts a string contains a substring
func (th *TestHelper) AssertContains(str, substr string, msgAndArgs ...interface{}) {
	if !strings.Contains(str, substr) {
		msg := fmt.Sprintf("Expected '%s' to contain '%s'", str, substr)
		if len(msgAndArgs) > 0 {
			if format, ok := msgAndArgs[0].(string); ok {
				msg = fmt.Sprintf(format, msgAndArgs[1:]...)
			}
		}
		th.t.Error(msg)
	}
}

// AssertPanic asserts a function panics
func (th *TestHelper) AssertPanic(fn func(), msgAndArgs ...interface{}) {
	defer func() {
		if r := recover(); r == nil {
			msg := "Expected function to panic, but it didn't"
			if len(msgAndArgs) > 0 {
				if format, ok := msgAndArgs[0].(string); ok {
					msg = fmt.Sprintf(format, msgAndArgs[1:]...)
				}
			}
			th.t.Error(msg)
		}
	}()

	fn()
}
