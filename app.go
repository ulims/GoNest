package gonest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// Application represents the main GoNest application
type Application struct {
	Echo                    *echo.Echo
	ModuleRegistry          *ModuleRegistry
	ServiceRegistry         *ServiceRegistry
	ControllerRegistry      *ControllerRegistry
	GuardRegistry           *GuardRegistry
	InterceptorRegistry     *InterceptorRegistry
	PipeRegistry            *PipeRegistry
	ExceptionFilterRegistry *ExceptionFilterRegistry
	WebSocketGateway        *WebSocketGateway
	DatabaseService         *DatabaseService
	MongoDBService          *MongoDBService
	LifecycleManager        *LifecycleManager
	Logger                  *logrus.Logger
	Config                  *Config
	Context                 context.Context
	Cancel                  context.CancelFunc
}

func (app *Application) Module(userModule *Module) {
	panic("unimplemented")
}

// Config represents application configuration
type Config struct {
	Port         string
	Host         string
	Environment  string
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Database     *DatabaseConfig
	MongoDB      *MongoDBConfig
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Port:         "8080",
		Host:         "localhost",
		Environment:  "development",
		LogLevel:     "info",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Database:     DefaultDatabaseConfig(),
		MongoDB:      DefaultMongoDBConfig(),
	}
}

// ApplicationBuilder provides a fluent interface for building applications
type ApplicationBuilder struct {
	app *Application
}

// NewApplication creates a new application
func NewApplication() *ApplicationBuilder {
	ctx, cancel := context.WithCancel(context.Background())

	return &ApplicationBuilder{
		app: &Application{
			Echo:                    echo.New(),
			ModuleRegistry:          NewModuleRegistry(),
			ServiceRegistry:         NewServiceRegistry(),
			ControllerRegistry:      NewControllerRegistry(),
			GuardRegistry:           NewGuardRegistry(),
			InterceptorRegistry:     NewInterceptorRegistry(),
			PipeRegistry:            NewPipeRegistry(),
			ExceptionFilterRegistry: NewExceptionFilterRegistry(),
			Logger:                  logrus.New(),
			Config:                  DefaultConfig(),
			Context:                 ctx,
			Cancel:                  cancel,
		},
	}
}

// Config sets the application configuration
func (ab *ApplicationBuilder) Config(config *Config) *ApplicationBuilder {
	ab.app.Config = config
	return ab
}

// Module adds a module to the application
func (ab *ApplicationBuilder) Module(module *Module) *ApplicationBuilder {
	ab.app.ModuleRegistry.Register(module)
	return ab
}

// Middleware adds middleware to the application
func (ab *ApplicationBuilder) Middleware(middleware ...echo.MiddlewareFunc) *ApplicationBuilder {
	ab.app.Echo.Use(middleware...)
	return ab
}

// Logger sets the application logger
func (ab *ApplicationBuilder) Logger(logger *logrus.Logger) *ApplicationBuilder {
	ab.app.Logger = logger
	return ab
}

// Database sets the database service
func (ab *ApplicationBuilder) Database(database *DatabaseService) *ApplicationBuilder {
	ab.app.DatabaseService = database
	return ab
}

// MongoDB sets the MongoDB service
func (ab *ApplicationBuilder) MongoDB(mongodb *MongoDBService) *ApplicationBuilder {
	ab.app.MongoDBService = mongodb
	return ab
}

// WebSocket sets the WebSocket gateway
func (ab *ApplicationBuilder) WebSocket(gateway *WebSocketGateway) *ApplicationBuilder {
	ab.app.WebSocketGateway = gateway
	return ab
}

// Lifecycle sets the lifecycle manager
func (ab *ApplicationBuilder) Lifecycle(manager *LifecycleManager) *ApplicationBuilder {
	ab.app.LifecycleManager = manager
	return ab
}

// Build returns the built application
func (ab *ApplicationBuilder) Build() *Application {
	return ab.app
}

// Initialize initializes the application
func (app *Application) Initialize() error {
	// Set up logger
	level, err := logrus.ParseLevel(app.Config.LogLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	app.Logger.SetLevel(level)

	// Set up Echo
	app.Echo.Logger.SetOutput(app.Logger.Writer())

	// Add default middleware
	app.Echo.Use(middleware.Logger())
	app.Echo.Use(middleware.Recover())
	app.Echo.Use(middleware.CORS())

	// Initialize lifecycle manager if not set
	if app.LifecycleManager == nil {
		app.LifecycleManager = NewLifecycleManager(app.Logger)
	}

	// Initialize WebSocket gateway if not set
	if app.WebSocketGateway == nil {
		// WebSocket gateway setup would be done manually by the user
	}

	// Initialize database service if not set
	if app.DatabaseService == nil && app.Config.Database != nil {
		app.DatabaseService = NewDatabaseService(app.Config.Database, app.Logger)
	}

	// Initialize MongoDB service if not set
	if app.MongoDBService == nil && app.Config.MongoDB != nil {
		app.MongoDBService = NewMongoDBService(app.Config.MongoDB, app.Logger)
	}

	// Register default lifecycle hooks
	app.registerDefaultLifecycleHooks()

	// Trigger application start lifecycle event
	if err := app.LifecycleManager.TriggerEvent(app.Context, EventApplicationStart); err != nil {
		return fmt.Errorf("failed to trigger application start event: %v", err)
	}

	// Initialize modules
	if err := app.initializeModules(); err != nil {
		return fmt.Errorf("failed to initialize modules: %v", err)
	}

	// Initialize services
	if err := app.initializeServices(); err != nil {
		return fmt.Errorf("failed to initialize services: %v", err)
	}

	// Initialize controllers
	if err := app.initializeControllers(); err != nil {
		return fmt.Errorf("failed to initialize controllers: %v", err)
	}

	// Set up routes
	app.setupRoutes()

	// Set up WebSocket routes
	app.setupWebSocketRoutes()

	return nil
}

// registerDefaultLifecycleHooks registers default lifecycle hooks
func (app *Application) registerDefaultLifecycleHooks() {
	// Register database lifecycle hook if database service exists
	if app.DatabaseService != nil {
		dbHook := NewDatabaseLifecycleHook(app.DatabaseService, app.Logger)
		app.LifecycleManager.RegisterHook(EventApplicationStart, dbHook, PriorityHigh)
		app.LifecycleManager.RegisterHook(EventApplicationStop, dbHook, PriorityHigh)
	}

	// Register MongoDB lifecycle hook if MongoDB service exists
	if app.MongoDBService != nil {
		mongoHook := NewMongoDBLifecycleHook(app.MongoDBService, app.Logger)
		app.LifecycleManager.RegisterHook(EventApplicationStart, mongoHook, PriorityHigh)
		app.LifecycleManager.RegisterHook(EventApplicationStop, mongoHook, PriorityHigh)
	}

	// WebSocket lifecycle hooks would be registered here if needed
}

// initializeModules initializes all modules
func (app *Application) initializeModules() error {
	modules := app.ModuleRegistry.GetAll()

	for name, module := range modules {
		app.Logger.Infof("Initializing module: %s", name)

		// Register sub-modules
		for _, subModule := range module.Modules {
			app.ModuleRegistry.Register(subModule)
		}

		// Register imports
		for _, importModule := range module.Imports {
			app.ModuleRegistry.Register(importModule)
		}

		// Trigger module init lifecycle event
		if err := app.LifecycleManager.TriggerEvent(app.Context, EventModuleInit); err != nil {
			return fmt.Errorf("failed to trigger module init event for %s: %v", name, err)
		}
	}

	return nil
}

// initializeServices initializes all services
func (app *Application) initializeServices() error {
	services := app.ServiceRegistry.GetAll()

	for name, service := range services {
		app.Logger.Infof("Initializing service: %s", name)

		// Inject dependencies
		if err := app.ServiceRegistry.Inject(service.Instance); err != nil {
			return fmt.Errorf("failed to inject dependencies for service %s: %v", name, err)
		}
	}

	return nil
}

// initializeControllers initializes all controllers
func (app *Application) initializeControllers() error {
	controllers := app.ControllerRegistry.GetControllers()

	for _, controller := range controllers {
		app.Logger.Infof("Initializing controller: %s", controller.Path)

		// Inject dependencies into controller
		if err := app.ServiceRegistry.Inject(controller); err != nil {
			return fmt.Errorf("failed to inject dependencies for controller %s: %v", controller.Path, err)
		}
	}

	return nil
}

// setupRoutes sets up all routes
func (app *Application) setupRoutes() {
	app.ControllerRegistry.SetupRoutes(app.Echo)
}

// setupWebSocketRoutes sets up WebSocket routes
func (app *Application) setupWebSocketRoutes() {
	// WebSocket routes can be manually added by user
	// WebSocket routes can be configured manually
}

// Start starts the application
func (app *Application) Start() error {
	// Initialize the application
	if err := app.Initialize(); err != nil {
		return err
	}

	// Start server
	go func() {
		addr := fmt.Sprintf("%s:%s", app.Config.Host, app.Config.Port)
		app.Logger.Infof("Starting server on %s", addr)

		if err := app.Echo.Start(addr); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	app.Logger.Info("Shutting down server...")
	app.Cancel()

	// Trigger application stop lifecycle event
	if err := app.LifecycleManager.TriggerEvent(app.Context, EventApplicationStop); err != nil {
		app.Logger.Errorf("Failed to trigger application stop event: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Echo.Shutdown(ctx); err != nil {
		app.Logger.Fatalf("Failed to shutdown server: %v", err)
	}

	app.Logger.Info("Server stopped")
	return nil
}

// Stop stops the application
func (app *Application) Stop() error {
	app.Cancel()
	return app.Echo.Shutdown(context.Background())
}

// GetService retrieves a service by name
func (app *Application) GetService(name string) (interface{}, bool) {
	return app.ServiceRegistry.Get(name)
}

// GetServiceByType retrieves a service by type
func (app *Application) GetServiceByType(serviceType interface{}) (interface{}, bool) {
	return app.ServiceRegistry.GetByType(reflect.TypeOf(serviceType))
}

// RegisterService registers a service
func (app *Application) RegisterService(name string, service interface{}) {
	app.ServiceRegistry.Register(name, service)
}

// RegisterController registers a controller
func (app *Application) RegisterController(controller *Controller) {
	app.ControllerRegistry.Register(controller)
}

// RegisterGuard registers a guard
func (app *Application) RegisterGuard(name string, guard Guard, priority int) {
	app.GuardRegistry.Register(name, guard, priority)
}

// RegisterInterceptor registers an interceptor
func (app *Application) RegisterInterceptor(name string, interceptor Interceptor, priority int) {
	app.InterceptorRegistry.Register(name, interceptor, priority)
}

// RegisterPipe registers a pipe
func (app *Application) RegisterPipe(name string, pipe Pipe, priority int) {
	app.PipeRegistry.Register(name, pipe, priority)
}

// RegisterExceptionFilter registers an exception filter
func (app *Application) RegisterExceptionFilter(name string, filter ExceptionFilter, priority int) {
	app.ExceptionFilterRegistry.Register(name, filter, priority)
}

// Use adds middleware to the application
func (app *Application) Use(middleware ...echo.MiddlewareFunc) {
	app.Echo.Use(middleware...)
}

// Group creates a new route group
func (app *Application) Group(prefix string) *echo.Group {
	return app.Echo.Group(prefix)
}

// Static serves static files
func (app *Application) Static(prefix, root string) {
	app.Echo.Static(prefix, root)
}

// File serves a file
func (app *Application) File(path, file string) {
	app.Echo.File(path, file)
}

// GET adds a GET route
func (app *Application) GET(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	app.Echo.GET(path, handler, middleware...)
}

// POST adds a POST route
func (app *Application) POST(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	app.Echo.POST(path, handler, middleware...)
}

// PUT adds a PUT route
func (app *Application) PUT(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	app.Echo.PUT(path, handler, middleware...)
}

// DELETE adds a DELETE route
func (app *Application) DELETE(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	app.Echo.DELETE(path, handler, middleware...)
}

// PATCH adds a PATCH route
func (app *Application) PATCH(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) {
	app.Echo.PATCH(path, handler, middleware...)
}
