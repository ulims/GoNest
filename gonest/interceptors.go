package gonest

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Interceptor interface for request/response transformation
type Interceptor interface {
	Intercept(ctx echo.Context, next echo.HandlerFunc) error
}

// InterceptorFunc is a function type that implements Interceptor interface
type InterceptorFunc func(ctx echo.Context, next echo.HandlerFunc) error

// Intercept implements Interceptor interface for InterceptorFunc
func (f InterceptorFunc) Intercept(ctx echo.Context, next echo.HandlerFunc) error {
	return f(ctx, next)
}

// InterceptorMetadata contains interceptor information
type InterceptorMetadata struct {
	Interceptor Interceptor
	Priority    int
	Metadata    map[string]interface{}
}

// InterceptorRegistry manages all interceptors
type InterceptorRegistry struct {
	interceptors map[string]*InterceptorMetadata
}

// NewInterceptorRegistry creates a new interceptor registry
func NewInterceptorRegistry() *InterceptorRegistry {
	return &InterceptorRegistry{
		interceptors: make(map[string]*InterceptorMetadata),
	}
}

// Register registers an interceptor
func (ir *InterceptorRegistry) Register(name string, interceptor Interceptor, priority int) {
	ir.interceptors[name] = &InterceptorMetadata{
		Interceptor: interceptor,
		Priority:    priority,
		Metadata:    make(map[string]interface{}),
	}
}

// Get retrieves an interceptor by name
func (ir *InterceptorRegistry) Get(name string) (Interceptor, bool) {
	if metadata, exists := ir.interceptors[name]; exists {
		return metadata.Interceptor, true
	}
	return nil, false
}

// GetAll returns all registered interceptors
func (ir *InterceptorRegistry) GetAll() map[string]*InterceptorMetadata {
	return ir.interceptors
}

// Interceptor decorators
type InterceptorDecorator struct {
	Interceptors []string
}

// UseInterceptors decorator for applying interceptors to routes
func UseInterceptors(interceptors ...string) InterceptorDecorator {
	return InterceptorDecorator{Interceptors: interceptors}
}

// LoggingInterceptor logs request/response information
type LoggingInterceptor struct {
	Logger *logrus.Logger
}

// NewLoggingInterceptor creates a new logging interceptor
func NewLoggingInterceptor(logger *logrus.Logger) *LoggingInterceptor {
	return &LoggingInterceptor{Logger: logger}
}

// Intercept logs request and response information
func (li *LoggingInterceptor) Intercept(ctx echo.Context, next echo.HandlerFunc) error {
	start := time.Now()

	// Log request
	li.Logger.WithFields(logrus.Fields{
		"method": ctx.Request().Method,
		"path":   ctx.Request().URL.Path,
		"ip":     ctx.RealIP(),
	}).Info("Request started")

	// Call next handler
	err := next(ctx)

	// Log response
	duration := time.Since(start)
	status := ctx.Response().Status

	li.Logger.WithFields(logrus.Fields{
		"method":   ctx.Request().Method,
		"path":     ctx.Request().URL.Path,
		"status":   status,
		"duration": duration,
	}).Info("Request completed")

	return err
}

// TransformInterceptor transforms request/response data
type TransformInterceptor struct {
	TransformRequest  func(interface{}) (interface{}, error)
	TransformResponse func(interface{}) (interface{}, error)
}

// NewTransformInterceptor creates a new transform interceptor
func NewTransformInterceptor(
	transformRequest func(interface{}) (interface{}, error),
	transformResponse func(interface{}) (interface{}, error),
) *TransformInterceptor {
	return &TransformInterceptor{
		TransformRequest:  transformRequest,
		TransformResponse: transformResponse,
	}
}

// Intercept transforms request and response data
func (ti *TransformInterceptor) Intercept(ctx echo.Context, next echo.HandlerFunc) error {
	// Transform request if needed
	if ti.TransformRequest != nil {
		// Implementation would transform request data
	}

	// Call next handler
	err := next(ctx)
	if err != nil {
		return err
	}

	// Transform response if needed
	if ti.TransformResponse != nil {
		// Implementation would transform response data
	}

	return nil
}

// MetricsInterceptor collects metrics
type MetricsInterceptor struct {
	Metrics MetricsService
}

// NewMetricsInterceptor creates a new metrics interceptor
func NewMetricsInterceptor(metrics MetricsService) *MetricsInterceptor {
	return &MetricsInterceptor{Metrics: metrics}
}

// Intercept collects metrics for the request
func (mi *MetricsInterceptor) Intercept(ctx echo.Context, next echo.HandlerFunc) error {
	start := time.Now()

	// Call next handler
	err := next(ctx)

	// Record metrics
	duration := time.Since(start)
	status := ctx.Response().Status

	mi.Metrics.RecordRequest(ctx.Request().Method, ctx.Request().URL.Path, status, duration)

	return err
}

// InterceptorMiddleware creates middleware from interceptors
func InterceptorMiddleware(interceptors ...Interceptor) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create a chain of interceptors
			handler := next
			for i := len(interceptors) - 1; i >= 0; i-- {
				interceptor := interceptors[i]
				prevHandler := handler
				handler = func(ctx echo.Context) error {
					return interceptor.Intercept(ctx, prevHandler)
				}
			}
			return handler(c)
		}
	}
}

// MetricsService interface for metrics collection
type MetricsService interface {
	RecordRequest(method, path string, status int, duration time.Duration)
	RecordError(method, path string, error string)
}
