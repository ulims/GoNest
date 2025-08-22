package gonest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Controller represents a NestJS-like controller
type Controller struct {
	Path       string
	Handlers   map[string]*Handler
	Middleware []echo.MiddlewareFunc
}

// Handler represents an HTTP handler with metadata
type Handler struct {
	Method      string
	Path        string
	HandlerFunc echo.HandlerFunc
	Middleware  []echo.MiddlewareFunc
}

// ControllerDecoratorFunc is a function type that can be used to decorate controllers
type ControllerDecoratorFunc func(*Controller)

// Route decorator for defining routes
type Route struct {
	Method string
	Path   string
}

// Controller decorators
func ControllerDecorator(path string) ControllerDecoratorFunc {
	return func(c *Controller) {
		c.Path = path
	}
}

func Get(path string) Route {
	return Route{Method: http.MethodGet, Path: path}
}

func Post(path string) Route {
	return Route{Method: http.MethodPost, Path: path}
}

func Put(path string) Route {
	return Route{Method: http.MethodPut, Path: path}
}

func Delete(path string) Route {
	return Route{Method: http.MethodDelete, Path: path}
}

func Patch(path string) Route {
	return Route{Method: http.MethodPatch, Path: path}
}

// ControllerBuilder provides a fluent interface for building controllers
type ControllerBuilder struct {
	controller *Controller
}

// NewController creates a new controller
func NewController() *ControllerBuilder {
	return &ControllerBuilder{
		controller: &Controller{
			Handlers:   make(map[string]*Handler),
			Middleware: make([]echo.MiddlewareFunc, 0),
		},
	}
}

// Path sets the base path for the controller
func (cb *ControllerBuilder) Path(path string) *ControllerBuilder {
	cb.controller.Path = path
	return cb
}

// Middleware adds middleware to the controller
func (cb *ControllerBuilder) Middleware(middleware ...echo.MiddlewareFunc) *ControllerBuilder {
	cb.controller.Middleware = append(cb.controller.Middleware, middleware...)
	return cb
}

// Get adds a GET route handler
func (cb *ControllerBuilder) Get(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *ControllerBuilder {
	cb.addHandler(http.MethodGet, path, handler, middleware...)
	return cb
}

// Post adds a POST route handler
func (cb *ControllerBuilder) Post(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *ControllerBuilder {
	cb.addHandler(http.MethodPost, path, handler, middleware...)
	return cb
}

// Put adds a PUT route handler
func (cb *ControllerBuilder) Put(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *ControllerBuilder {
	cb.addHandler(http.MethodPut, path, handler, middleware...)
	return cb
}

// Delete adds a DELETE route handler
func (cb *ControllerBuilder) Delete(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *ControllerBuilder {
	cb.addHandler(http.MethodDelete, path, handler, middleware...)
	return cb
}

// Patch adds a PATCH route handler
func (cb *ControllerBuilder) Patch(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *ControllerBuilder {
	cb.addHandler(http.MethodPatch, path, handler, middleware...)
	return cb
}

// addHandler adds a handler to the controller
func (cb *ControllerBuilder) addHandler(method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	key := method + ":" + path
	cb.controller.Handlers[key] = &Handler{
		Method:      method,
		Path:        path,
		HandlerFunc: handler,
		Middleware:  middleware,
	}
}

// Build returns the built controller
func (cb *ControllerBuilder) Build() *Controller {
	return cb.controller
}

// ControllerRegistry manages all controllers
type ControllerRegistry struct {
	controllers []*Controller
}

// NewControllerRegistry creates a new controller registry
func NewControllerRegistry() *ControllerRegistry {
	return &ControllerRegistry{
		controllers: make([]*Controller, 0),
	}
}

// Register registers a controller
func (cr *ControllerRegistry) Register(controller *Controller) {
	cr.controllers = append(cr.controllers, controller)
}

// RegisterAll registers all controllers from a module
func (cr *ControllerRegistry) RegisterAll(controllers []interface{}) {
	for _, controller := range controllers {
		if ctrl, ok := controller.(*Controller); ok {
			cr.Register(ctrl)
		}
	}
}

// SetupRoutes sets up all controller routes on an Echo instance
func (cr *ControllerRegistry) SetupRoutes(e *echo.Echo) {
	for _, controller := range cr.controllers {
		group := e.Group(controller.Path)

		// Apply controller-level middleware
		if len(controller.Middleware) > 0 {
			group.Use(controller.Middleware...)
		}

		// Register all handlers
		for _, handler := range controller.Handlers {
			group.Add(handler.Method, handler.Path, handler.HandlerFunc, handler.Middleware...)
		}
	}
}

// GetControllers returns all registered controllers
func (cr *ControllerRegistry) GetControllers() []*Controller {
	return cr.controllers
}

// CreateRoute creates a route decorator for method and path
func CreateRoute(method, path string) Route {
	return Route{Method: method, Path: path}
}

// Body decorator for request body binding
func Body(model interface{}) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := c.Bind(model); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return nil
	}
}

// Param decorator for path parameter binding
func Param(name string) echo.HandlerFunc {
	return func(c echo.Context) error {
		param := c.Param(name)
		if param == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing required parameter: "+name)
		}
		return nil
	}
}

// Query decorator for query parameter binding
func Query(name string) echo.HandlerFunc {
	return func(c echo.Context) error {
		query := c.QueryParam(name)
		if query == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Missing required query parameter: "+name)
		}
		return nil
	}
}
