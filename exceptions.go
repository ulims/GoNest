package gonest

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// ExceptionFilter interface for error handling
type ExceptionFilter interface {
	Catch(exception interface{}, ctx echo.Context) error
}

// ExceptionFilterFunc is a function type that implements ExceptionFilter interface
type ExceptionFilterFunc func(exception interface{}, ctx echo.Context) error

// Catch implements ExceptionFilter interface for ExceptionFilterFunc
func (f ExceptionFilterFunc) Catch(exception interface{}, ctx echo.Context) error {
	return f(exception, ctx)
}

// ExceptionFilterMetadata contains filter information
type ExceptionFilterMetadata struct {
	Filter   ExceptionFilter
	Priority int
	Metadata map[string]interface{}
}

// ExceptionFilterRegistry manages all exception filters
type ExceptionFilterRegistry struct {
	filters map[string]*ExceptionFilterMetadata
}

// NewExceptionFilterRegistry creates a new exception filter registry
func NewExceptionFilterRegistry() *ExceptionFilterRegistry {
	return &ExceptionFilterRegistry{
		filters: make(map[string]*ExceptionFilterMetadata),
	}
}

// Register registers an exception filter
func (efr *ExceptionFilterRegistry) Register(name string, filter ExceptionFilter, priority int) {
	efr.filters[name] = &ExceptionFilterMetadata{
		Filter:   filter,
		Priority: priority,
		Metadata: make(map[string]interface{}),
	}
}

// Get retrieves an exception filter by name
func (efr *ExceptionFilterRegistry) Get(name string) (ExceptionFilter, bool) {
	if metadata, exists := efr.filters[name]; exists {
		return metadata.Filter, true
	}
	return nil, false
}

// GetAll returns all registered exception filters
func (efr *ExceptionFilterRegistry) GetAll() map[string]*ExceptionFilterMetadata {
	return efr.filters
}

// ExceptionFilter decorators
type ExceptionFilterDecorator struct {
	Filters []string
}

// UseExceptionFilters decorator for applying exception filters
func UseExceptionFilters(filters ...string) ExceptionFilterDecorator {
	return ExceptionFilterDecorator{Filters: filters}
}

// HTTPException represents an HTTP error
type HTTPException struct {
	Status  int
	Message string
	Code    string
	Details interface{}
}

// NewHTTPException creates a new HTTP exception
func NewHTTPException(status int, message string) *HTTPException {
	return &HTTPException{
		Status:  status,
		Message: message,
	}
}

// Error implements error interface
func (he *HTTPException) Error() string {
	return he.Message
}

// WithCode adds an error code to the exception
func (he *HTTPException) WithCode(code string) *HTTPException {
	he.Code = code
	return he
}

// WithDetails adds details to the exception
func (he *HTTPException) WithDetails(details interface{}) *HTTPException {
	he.Details = details
	return he
}

// BadRequestException creates a 400 Bad Request exception
func BadRequestException(message string) *HTTPException {
	return NewHTTPException(http.StatusBadRequest, message)
}

// UnauthorizedException creates a 401 Unauthorized exception
func UnauthorizedException(message string) *HTTPException {
	return NewHTTPException(http.StatusUnauthorized, message)
}

// ForbiddenException creates a 403 Forbidden exception
func ForbiddenException(message string) *HTTPException {
	return NewHTTPException(http.StatusForbidden, message)
}

// NotFoundException creates a 404 Not Found exception
func NotFoundException(message string) *HTTPException {
	return NewHTTPException(http.StatusNotFound, message)
}

// ConflictException creates a 409 Conflict exception
func ConflictException(message string) *HTTPException {
	return NewHTTPException(http.StatusConflict, message)
}

// InternalServerErrorException creates a 500 Internal Server Error exception
func InternalServerErrorException(message string) *HTTPException {
	return NewHTTPException(http.StatusInternalServerError, message)
}

// ValidationException creates a validation error exception
type ValidationException struct {
	Errors map[string]string
}

// NewValidationException creates a new validation exception
func NewValidationException(errors map[string]string) *ValidationException {
	return &ValidationException{Errors: errors}
}

// Error implements error interface
func (ve *ValidationException) Error() string {
	return "Validation failed"
}

// HTTPExceptionFilter handles HTTP exceptions
type HTTPExceptionFilter struct {
	Logger *logrus.Logger
}

// NewHTTPExceptionFilter creates a new HTTP exception filter
func NewHTTPExceptionFilter(logger *logrus.Logger) *HTTPExceptionFilter {
	return &HTTPExceptionFilter{Logger: logger}
}

// Catch handles HTTP exceptions
func (hef *HTTPExceptionFilter) Catch(exception interface{}, ctx echo.Context) error {
	if httpException, ok := exception.(*HTTPException); ok {
		hef.Logger.WithFields(logrus.Fields{
			"status":  httpException.Status,
			"message": httpException.Message,
			"path":    ctx.Request().URL.Path,
			"method":  ctx.Request().Method,
		}).Error("HTTP Exception caught")

		response := map[string]interface{}{
			"error":  httpException.Message,
			"status": httpException.Status,
			"path":   ctx.Request().URL.Path,
			"method": ctx.Request().Method,
		}

		if httpException.Code != "" {
			response["code"] = httpException.Code
		}

		if httpException.Details != nil {
			response["details"] = httpException.Details
		}

		return ctx.JSON(httpException.Status, response)
	}

	return nil
}

// ValidationExceptionFilter handles validation exceptions
type ValidationExceptionFilter struct {
	Logger *logrus.Logger
}

// NewValidationExceptionFilter creates a new validation exception filter
func NewValidationExceptionFilter(logger *logrus.Logger) *ValidationExceptionFilter {
	return &ValidationExceptionFilter{Logger: logger}
}

// Catch handles validation exceptions
func (vef *ValidationExceptionFilter) Catch(exception interface{}, ctx echo.Context) error {
	if validationException, ok := exception.(*ValidationException); ok {
		vef.Logger.WithFields(logrus.Fields{
			"errors": validationException.Errors,
			"path":   ctx.Request().URL.Path,
			"method": ctx.Request().Method,
		}).Error("Validation Exception caught")

		response := map[string]interface{}{
			"error":   "Validation failed",
			"status":  http.StatusBadRequest,
			"path":    ctx.Request().URL.Path,
			"method":  ctx.Request().Method,
			"details": validationException.Errors,
		}

		return ctx.JSON(http.StatusBadRequest, response)
	}

	return nil
}

// GenericExceptionFilter handles generic exceptions
type GenericExceptionFilter struct {
	Logger *logrus.Logger
}

// NewGenericExceptionFilter creates a new generic exception filter
func NewGenericExceptionFilter(logger *logrus.Logger) *GenericExceptionFilter {
	return &GenericExceptionFilter{Logger: logger}
}

// Catch handles generic exceptions
func (gef *GenericExceptionFilter) Catch(exception interface{}, ctx echo.Context) error {
	// Log the exception
	gef.Logger.WithFields(logrus.Fields{
		"exception": fmt.Sprintf("%v", exception),
		"type":      reflect.TypeOf(exception).String(),
		"path":      ctx.Request().URL.Path,
		"method":    ctx.Request().Method,
	}).Error("Generic Exception caught")

	// Return a generic error response
	response := map[string]interface{}{
		"error":  "Internal server error",
		"status": http.StatusInternalServerError,
		"path":   ctx.Request().URL.Path,
		"method": ctx.Request().Method,
	}

	return ctx.JSON(http.StatusInternalServerError, response)
}

// ExceptionFilterChain chains multiple exception filters
type ExceptionFilterChain struct {
	filters []ExceptionFilter
}

// NewExceptionFilterChain creates a new exception filter chain
func NewExceptionFilterChain(filters ...ExceptionFilter) *ExceptionFilterChain {
	return &ExceptionFilterChain{filters: filters}
}

// Catch applies all filters in sequence
func (efc *ExceptionFilterChain) Catch(exception interface{}, ctx echo.Context) error {
	for _, filter := range efc.filters {
		if err := filter.Catch(exception, ctx); err != nil {
			return err
		}
	}
	return nil
}

// ExceptionFilterMiddleware creates middleware from exception filters
func ExceptionFilterMiddleware(filters ...ExceptionFilter) echo.MiddlewareFunc {
	chain := NewExceptionFilterChain(filters...)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					// Apply exception filters
					if err := chain.Catch(r, c); err != nil {
						// If filter fails, return generic error
						c.JSON(http.StatusInternalServerError, map[string]interface{}{
							"error":  "Internal server error",
							"status": http.StatusInternalServerError,
						})
					}
				}
			}()

			return next(c)
		}
	}
}

// GlobalExceptionHandler handles all exceptions globally
type GlobalExceptionHandler struct {
	filters []ExceptionFilter
	logger  *logrus.Logger
}

// NewGlobalExceptionHandler creates a new global exception handler
func NewGlobalExceptionHandler(logger *logrus.Logger) *GlobalExceptionHandler {
	return &GlobalExceptionHandler{
		filters: []ExceptionFilter{
			NewHTTPExceptionFilter(logger),
			NewValidationExceptionFilter(logger),
			NewGenericExceptionFilter(logger),
		},
		logger: logger,
	}
}

// Handle handles an exception using all registered filters
func (geh *GlobalExceptionHandler) Handle(exception interface{}, ctx echo.Context) error {
	chain := NewExceptionFilterChain(geh.filters...)
	return chain.Catch(exception, ctx)
}

// AddFilter adds a filter to the global handler
func (geh *GlobalExceptionHandler) AddFilter(filter ExceptionFilter) {
	geh.filters = append(geh.filters, filter)
}

// Built-in exception filter instances
var (
	HTTPExceptionFilterInstance       = NewHTTPExceptionFilter(nil)
	ValidationExceptionFilterInstance = NewValidationExceptionFilter(nil)
	GenericExceptionFilterInstance    = NewGenericExceptionFilter(nil)
)
