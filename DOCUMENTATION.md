# GoNest Framework Documentation

GoNest is a powerful, scalable Go framework inspired by NestJS that provides a complete solution for building enterprise-grade applications. It combines the performance of Go with the elegant architecture patterns of NestJS.

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Core Concepts](#core-concepts)
5. [Modules](#modules)
6. [Controllers](#controllers)
7. [Services](#services)
8. [Dependency Injection](#dependency-injection)
9. [Data Transfer Objects (DTOs)](#data-transfer-objects-dtos)
10. [Validation](#validation)
11. [Authentication & Authorization](#authentication--authorization)
12. [Guards](#guards)
13. [Interceptors](#interceptors)
14. [Exception Filters](#exception-filters)
15. [Pipes](#pipes)
16. [Lifecycle Hooks](#lifecycle-hooks)
17. [WebSockets](#websockets)
18. [Caching](#caching)
19. [Events](#events)
20. [Rate Limiting](#rate-limiting)
21. [Configuration Management](#configuration-management)
22. [Testing](#testing)
23. [MongoDB Integration](#mongodb-integration)
24. [Best Practices](#best-practices)
25. [API Reference](#api-reference)

## Overview

GoNest provides a complete framework for building scalable, maintainable Go applications with:

- **Modular Architecture**: Organize code into modules with clear boundaries
- **Dependency Injection**: Automatic dependency resolution and injection
- **Decorators**: Use decorators for metadata and configuration
- **Guards**: Route-level authentication and authorization
- **Interceptors**: Request/response transformation and logging
- **Exception Filters**: Centralized error handling
- **WebSockets**: Real-time communication support
- **Caching**: Built-in caching with multiple providers
- **Events**: Event-driven architecture
- **Rate Limiting**: Protect your APIs from abuse
- **MongoDB Integration**: Simplified database operations
- **Testing Utilities**: Comprehensive testing support

## Installation

```bash
go get github.com/your-username/gonest
```

## Quick Start

```go
package main

import (
    "github.com/labstack/echo/v4"
    "gonest"
)

// User represents a user entity
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

// UserService handles user business logic
type UserService struct {
    users map[string]*User
}

func NewUserService() *UserService {
    return &UserService{
        users: make(map[string]*User),
    }
}

func (s *UserService) CreateUser(user *User) error {
    s.users[user.ID] = user
    return nil
}

func (s *UserService) GetUser(id string) (*User, error) {
    user, exists := s.users[id]
    if !exists {
        return nil, gonest.NotFoundException("User not found")
    }
    return user, nil
}

// UserController handles HTTP requests
type UserController struct {
    userService *UserService `inject:"UserService"`
}

func (c *UserController) CreateUser(ctx echo.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        return gonest.BadRequestException("Invalid request")
    }

    if err := c.userService.CreateUser(&user); err != nil {
        return err
    }

    return ctx.JSON(201, user)
}

func (c *UserController) GetUser(ctx echo.Context) error {
    id := ctx.Param("id")
    user, err := c.userService.GetUser(id)
    if err != nil {
        return err
    }
    return ctx.JSON(200, user)
}

func main() {
    // Create module
    userModule := gonest.NewModule("UserModule").
        Controller(&UserController{}).
        Service(NewUserService()).
        Build()

    // Create application
    app := gonest.NewApplication().
        Config(&gonest.Config{Port: "8080"}).
        Logger(logrus.New()).
        Build()

    // Register module
    app.ModuleRegistry.Register(userModule)

    // Start application
    app.Start()
}
```

## Core Concepts

### Modules

Modules are the basic building blocks of GoNest applications. They encapsulate related functionality and can contain controllers, services, and other modules.

```go
type Module struct {
    Name        string
    Controllers []interface{}
    Services    []interface{}
    Modules     []*Module
    Providers   []interface{}
    Imports     []*Module
    Exports     []interface{}
}
```

### Controllers

Controllers handle incoming HTTP requests and return responses. They use decorators to define routes and HTTP methods.

```go
type UserController struct {
    userService *UserService `inject:"UserService"`
}

// @Get("/users/:id")
func (c *UserController) GetUser(ctx echo.Context) error {
    // Handle request
}
```

### Services

Services contain business logic and are injected into controllers and other services.

```go
type UserService struct {
    userRepository *UserRepository `inject:"UserRepository"`
    cacheService   *CacheService   `inject:"CacheService"`
}

func (s *UserService) CreateUser(user *User) error {
    // Business logic
}
```

## Modules

### Creating Modules

```go
userModule := gonest.NewModule("UserModule").
    Controller(&UserController{}).
    Service(NewUserService()).
    Provider(NewUserRepository()).
    Build()
```

### Module Imports

```go
authModule := gonest.NewModule("AuthModule").
    Service(NewAuthService()).
    Build()

userModule := gonest.NewModule("UserModule").
    Import(authModule).
    Controller(&UserController{}).
    Build()
```

### Module Exports

```go
userModule := gonest.NewModule("UserModule").
    Service(NewUserService()).
    Export(NewUserService()).
    Build()
```

## Controllers

### Basic Controller

```go
type UserController struct {
    userService *UserService `inject:"UserService"`
}

func (c *UserController) GetUsers(ctx echo.Context) error {
    users, err := c.userService.GetAllUsers()
    if err != nil {
        return err
    }
    return ctx.JSON(200, users)
}
```

### Controller with Decorators

```go
type UserController struct {
    userService *UserService `inject:"UserService"`
}

// @Get("/users")
// @UseGuards(AuthGuard)
// @UseInterceptors(LoggingInterceptor)
func (c *UserController) GetUsers(ctx echo.Context) error {
    // Implementation
}
```

### Route Parameters

```go
func (c *UserController) GetUser(ctx echo.Context) error {
    id := ctx.Param("id")
    user, err := c.userService.GetUser(id)
    if err != nil {
        return err
    }
    return ctx.JSON(200, user)
}
```

### Request Body

```go
func (c *UserController) CreateUser(ctx echo.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        return gonest.BadRequestException("Invalid request body")
    }
    
    if err := c.userService.CreateUser(&user); err != nil {
        return err
    }
    
    return ctx.JSON(201, user)
}
```

## Services

### Basic Service

```go
type UserService struct {
    userRepository *UserRepository `inject:"UserRepository"`
}

func (s *UserService) GetUser(id string) (*User, error) {
    return s.userRepository.FindByID(id)
}

func (s *UserService) CreateUser(user *User) error {
    return s.userRepository.Create(user)
}
```

### Service with Dependencies

```go
type UserService struct {
    userRepository *UserRepository `inject:"UserRepository"`
    cacheService   *CacheService   `inject:"CacheService"`
    eventService   *EventService   `inject:"EventService"`
}

func (s *UserService) GetUser(id string) (*User, error) {
    // Try cache first
    if cached, err := s.cacheService.Get(context.Background(), "user:"+id, &User{}); err == nil {
        return cached.(*User), nil
    }
    
    // Get from repository
    user, err := s.userRepository.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    s.cacheService.Set(context.Background(), "user:"+id, user, 5*time.Minute)
    
    return user, nil
}
```

## Dependency Injection

GoNest provides automatic dependency injection using struct tags.

### Basic Injection

```go
type UserController struct {
    userService *UserService `inject:"UserService"`
}
```

### Named Injection

```go
type UserController struct {
    primaryDB   *Database `inject:"PrimaryDatabase"`
    secondaryDB *Database `inject:"SecondaryDatabase"`
}
```

### Interface Injection

```go
type UserController struct {
    userService UserServiceInterface `inject:"UserService"`
}
```

## Data Transfer Objects (DTOs)

DTOs define the structure of request and response data.

### Creating DTOs

```go
userDTO := gonest.NewDTO().
    Field("name", reflect.TypeOf(""), map[string]string{
        "required": "true",
        "min":      "2",
        "max":      "50",
    }).
    Field("email", reflect.TypeOf(""), map[string]string{
        "required": "true",
        "email":    "true",
    }).
    Build()
```

### Using DTOs

```go
func (c *UserController) CreateUser(ctx echo.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        return gonest.BadRequestException("Invalid request")
    }
    
    // Validate using DTO
    validator := gonest.NewDTOValidator()
    if err := gonest.ValidateStruct(&user, validator); err != nil {
        return gonest.BadRequestException("Validation failed")
    }
    
    return c.userService.CreateUser(&user)
}
```

## Validation

### Struct Validation

```go
type User struct {
    Name     string `validate:"required,min=2,max=50"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"min=13,max=120"`
    Password string `validate:"required,min=6"`
}
```

### Custom Validation

```go
type User struct {
    Email string `validate:"required,email,unique_email"`
}

// Register custom validator
validator := gonest.NewDTOValidator()
validator.RegisterValidation("unique_email", func(fl validator.FieldLevel) bool {
    // Custom validation logic
    return true
})
```

## Authentication & Authorization

### JWT Authentication

```go
// Configure JWT
jwtConfig := gonest.DefaultJWTConfig()
jwtConfig.SecretKey = "your-secret-key"
jwtConfig.TokenExpiry = 24 * time.Hour

authService := gonest.NewAuthService(jwtConfig, logger)

// Use in controller
func (c *UserController) GetProfile(ctx echo.Context) error {
    user, err := gonest.GetCurrentUser(ctx)
    if err != nil {
        return gonest.UnauthorizedException("Not authenticated")
    }
    return ctx.JSON(200, user)
}
```

### Login Endpoint

```go
func (c *AuthController) Login(ctx echo.Context) error {
    var loginRequest gonest.LoginRequest
    if err := ctx.Bind(&loginRequest); err != nil {
        return gonest.BadRequestException("Invalid request")
    }
    
    token, err := c.authService.Login(loginRequest.Username, loginRequest.Password)
    if err != nil {
        return gonest.UnauthorizedException("Invalid credentials")
    }
    
    return ctx.JSON(200, map[string]string{
        "token": token,
    })
}
```

## Guards

Guards determine whether a request should be handled by the route handler.

### Authentication Guard

```go
type AuthGuard struct{}

func (ag *AuthGuard) CanActivate(ctx echo.Context) (bool, error) {
    token := ctx.Request().Header.Get("Authorization")
    if token == "" {
        return false, gonest.UnauthorizedException("No token provided")
    }
    
    // Validate token
    return true, nil
}
```

### Role Guard

```go
type RoleGuard struct {
    requiredRoles []string
}

func (rg *RoleGuard) CanActivate(ctx echo.Context) (bool, error) {
    user, err := gonest.GetCurrentUser(ctx)
    if err != nil {
        return false, err
    }
    
    for _, role := range rg.requiredRoles {
        if contains(user.Roles, role) {
            return true, nil
        }
    }
    
    return false, gonest.ForbiddenException("Insufficient permissions")
}
```

### Using Guards

```go
// @UseGuards(AuthGuard, RoleGuard{"admin"})
func (c *UserController) AdminOnly(ctx echo.Context) error {
    // Only accessible by authenticated admins
}
```

## Interceptors

Interceptors transform requests and responses, and add cross-cutting concerns.

### Logging Interceptor

```go
type LoggingInterceptor struct{}

func (li *LoggingInterceptor) Intercept(ctx echo.Context, next echo.HandlerFunc) error {
    start := time.Now()
    
    // Log request
    log.Printf("Request: %s %s", ctx.Request().Method, ctx.Request().URL.Path)
    
    // Call next handler
    err := next(ctx)
    
    // Log response
    duration := time.Since(start)
    log.Printf("Response: %d in %v", ctx.Response().Status, duration)
    
    return err
}
```

### Caching Interceptor

```go
type CacheInterceptor struct {
    cacheService *CacheService
}

func (ci *CacheInterceptor) Intercept(ctx echo.Context, next echo.HandlerFunc) error {
    if ctx.Request().Method != "GET" {
        return next(ctx)
    }
    
    cacheKey := fmt.Sprintf("%s:%s", ctx.Request().Method, ctx.Request().URL.Path)
    
    // Try to get from cache
    if cached, err := ci.cacheService.Get(ctx.Request().Context(), cacheKey, nil); err == nil {
        return ctx.JSON(200, cached)
    }
    
    // Call next handler and cache result
    return next(ctx)
}
```

## Exception Filters

Exception filters handle exceptions thrown during request processing.

### Global Exception Filter

```go
type GlobalExceptionFilter struct{}

func (gef *GlobalExceptionFilter) Catch(exception interface{}, ctx echo.Context) error {
    switch e := exception.(type) {
    case *gonest.HTTPException:
        return ctx.JSON(e.StatusCode, map[string]interface{}{
            "error":   e.Message,
            "code":    e.StatusCode,
            "path":    ctx.Request().URL.Path,
            "method":  ctx.Request().Method,
            "time":    time.Now(),
        })
    default:
        return ctx.JSON(500, map[string]interface{}{
            "error": "Internal server error",
            "time":  time.Now(),
        })
    }
}
```

### Validation Exception Filter

```go
type ValidationExceptionFilter struct{}

func (vef *ValidationExceptionFilter) Catch(exception interface{}, ctx echo.Context) error {
    if validationErr, ok := exception.(validator.ValidationErrors); ok {
        errors := make(map[string]string)
        for _, err := range validationErr {
            errors[err.Field()] = err.Tag()
        }
        
        return ctx.JSON(400, map[string]interface{}{
            "error":   "Validation failed",
            "details": errors,
        })
    }
    
    return nil
}
```

## Pipes

Pipes transform input data before it reaches the route handler.

### Transformation Pipe

```go
type TransformPipe struct{}

func (tp *TransformPipe) Transform(value interface{}) (interface{}, error) {
    if str, ok := value.(string); ok {
        return strings.ToUpper(str), nil
    }
    return value, nil
}
```

### Validation Pipe

```go
type ValidationPipe struct {
    validator *validator.Validate
}

func (vp *ValidationPipe) Transform(value interface{}) (interface{}, error) {
    if err := vp.validator.Struct(value); err != nil {
        return nil, err
    }
    return value, nil
}
```

## Lifecycle Hooks

Lifecycle hooks allow you to run code at specific points in the application lifecycle.

### Application Lifecycle

```go
type AppLifecycleHook struct{}

func (alh *AppLifecycleHook) OnApplicationStart(ctx context.Context) error {
    log.Println("Application starting...")
    return nil
}

func (alh *AppLifecycleHook) OnApplicationStop(ctx context.Context) error {
    log.Println("Application stopping...")
    return nil
}
```

### Module Lifecycle

```go
type ModuleLifecycleHook struct{}

func (mlh *ModuleLifecycleHook) OnModuleInit(ctx context.Context) error {
    log.Println("Module initializing...")
    return nil
}

func (mlh *ModuleLifecycleHook) OnModuleDestroy(ctx context.Context) error {
    log.Println("Module destroying...")
    return nil
}
```

## WebSockets

### WebSocket Gateway

```go
type ChatGateway struct {
    *gonest.BaseWebSocketGateway
}

func NewChatGateway(logger *logrus.Logger) *ChatGateway {
    return &ChatGateway{
        BaseWebSocketGateway: gonest.NewBaseWebSocketGateway("/chat", logger),
    }
}

func (cg *ChatGateway) OnConnection(client *gonest.WebSocketClient) error {
    client.JoinRoom("general")
    return client.Emit("welcome", map[string]interface{}{
        "message": "Welcome to the chat!",
    })
}

func (cg *ChatGateway) OnMessage(client *gonest.WebSocketClient, message *gonest.WebSocketMessage) error {
    switch message.Event {
    case "chat_message":
        // Broadcast to room
        return cg.BroadcastToRoom("general", message.Event, message.Data)
    }
    return nil
}
```

### WebSocket Client

```go
type WebSocketClient struct {
    ID       string
    Conn     *websocket.Conn
    Rooms    map[string]bool
    Send     chan []byte
}
```

## Caching

### Memory Cache

```go
cacheService := gonest.NewCacheService(gonest.NewMemoryCache(logger), logger)

// Set cache
cacheService.Set(context.Background(), "key", "value", 5*time.Minute)

// Get cache
value, err := cacheService.Get(context.Background(), "key", "")
```

### Cache Interceptor

```go
cacheInterceptor := &gonest.CacheInterceptor{
    CacheService: cacheService,
    Condition: func(c echo.Context) bool {
        return c.Request().Method == "GET"
    },
}
```

## Events

### Event Service

```go
eventBus := gonest.NewEventBus(logger)
eventService := gonest.NewEventService(eventBus, "app", logger)

// Emit event
eventService.Emit(context.Background(), "user.created", map[string]interface{}{
    "user_id": "123",
    "email":   "user@example.com",
})

// Listen to event
eventService.On("user.created", func(ctx context.Context, event *gonest.Event) error {
    log.Printf("User created: %v", event.Data)
    return nil
})
```

## Rate Limiting

### Fixed Window Rate Limiter

```go
rateLimiter := gonest.NewFixedWindowRateLimiter(100, time.Minute)
middleware := rateLimiter.Middleware(logger)

// Use in routes
app.Use(middleware)
```

### Token Bucket Rate Limiter

```go
rateLimiter := gonest.NewTokenBucketRateLimiter(10, 10, time.Minute)
middleware := rateLimiter.Middleware(logger)
```

### Sliding Window Rate Limiter

```go
rateLimiter := gonest.NewSlidingWindowRateLimiter(50, time.Minute)
middleware := rateLimiter.Middleware(logger)
```

## Configuration Management

### Environment Configuration

```go
configService := gonest.NewConfigService(logger)
configService.AddProvider(gonest.NewEnvironmentConfigProvider("APP_"))

var config struct {
    Port     string `config:"port,default=8080"`
    Database string `config:"database,required"`
}

if err := gonest.InjectConfig(configService, &config); err != nil {
    log.Fatal(err)
}
```

### File Configuration

```go
configService := gonest.NewConfigService(logger)
configService.AddProvider(gonest.NewFileConfigProvider("config.yaml"))
```

## Testing

### Test Application

```go
func TestUserController(t *testing.T) {
    testApp := gonest.NewTestApp().
        WithModule(userModule).
        WithConfig(&gonest.Config{Port: "0"}).
        Build()
    
    defer testApp.Stop()
    
    // Test request
    response := testApp.Request("POST", "/users").
        WithJSON(map[string]interface{}{
            "name":  "John Doe",
            "email": "john@example.com",
        }).
        ExpectStatus(201).
        ExpectJSON().
        Get()
    
    assert.Equal(t, "John Doe", response.JSON()["name"])
}
```

### HTTP Testing

```go
func TestGetUser(t *testing.T) {
    testApp := gonest.NewTestApp().Build()
    defer testApp.Stop()
    
    response := testApp.Request("GET", "/users/123").
        WithHeader("Authorization", "Bearer token").
        ExpectStatus(200).
        Get()
    
    assert.NotNil(t, response.JSON())
}
```

## MongoDB Integration

### MongoDB Service

```go
mongoConfig := gonest.DefaultMongoDBConfig()
mongoConfig.URI = "mongodb://localhost:27017"
mongoConfig.Database = "gonest"

mongoService := gonest.NewMongoDBService(mongoConfig, logger)
if err := mongoService.Connect(); err != nil {
    log.Fatal(err)
}
```

### Schema Definition

```go
userSchema := gonest.NewSchema("users").
    Field("name", gonest.String).IsRequired().MinLength(2).MaxLength(50).
    Field("email", gonest.String).IsRequired().Pattern(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).
    Field("age", gonest.Number).Min(13).Max(120).
    HasTimestamps().
    Build()
```

### Model Operations

```go
userModel := gonest.NewMongoDBModel(userSchema, mongoService)

// Create
user := &User{Name: "John", Email: "john@example.com"}
if err := userModel.Create(context.Background(), user); err != nil {
    return err
}

// Find by ID
var foundUser User
if err := userModel.FindById(context.Background(), user.ID, &foundUser); err != nil {
    return err
}

// Find with query
query := gonest.NewMongoDBQuery().
    Where("age", ">=", 18).
    Where("role", "=", "admin").
    Sort("created_at", -1).
    Limit(10)

var users []User
if err := userModel.Find(context.Background(), query, &users); err != nil {
    return err
}
```

## Best Practices

### Project Structure

```
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── modules/
│   │   ├── user/
│   │   │   ├── controller.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── dto.go
│   │   └── auth/
│   │       ├── controller.go
│   │       ├── service.go
│   │       └── guard.go
│   ├── config/
│   │   └── config.go
│   └── shared/
│       ├── middleware/
│       ├── interceptors/
│       └── exceptions/
├── pkg/
│   └── utils/
├── go.mod
└── go.sum
```

### Error Handling

```go
// Use framework exceptions
if user == nil {
    return gonest.NotFoundException("User not found")
}

if err != nil {
    return gonest.BadRequestException("Invalid request")
}

// Custom exceptions
type CustomException struct {
    gonest.HTTPException
    Code string `json:"code"`
}

func NewCustomException(message, code string) *CustomException {
    return &CustomException{
        HTTPException: *gonest.NewHTTPException(400, message),
        Code:          code,
    }
}
```

### Logging

```go
logger := logrus.New()
logger.SetLevel(logrus.InfoLevel)
logger.SetFormatter(&logrus.JSONFormatter{})

// In services
logger.WithFields(logrus.Fields{
    "user_id": user.ID,
    "action":  "user_created",
}).Info("User created successfully")
```

### Configuration

```go
type Config struct {
    Server struct {
        Port string `yaml:"port" default:"8080"`
        Host string `yaml:"host" default:"localhost"`
    } `yaml:"server"`
    
    Database struct {
        URI      string `yaml:"uri" required:"true"`
        Database string `yaml:"database" default:"gonest"`
    } `yaml:"database"`
    
    Auth struct {
        JWTSecret string        `yaml:"jwt_secret" required:"true"`
        Expiry    time.Duration `yaml:"expiry" default:"24h"`
    } `yaml:"auth"`
}
```

## API Reference

### Application

```go
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
```

### Module

```go
type Module struct {
    Name        string
    Controllers []interface{}
    Services    []interface{}
    Modules     []*Module
    Providers   []interface{}
    Imports     []*Module
    Exports     []interface{}
}
```

### Service Registry

```go
type ServiceRegistry struct {
    services map[string]*ServiceEntry
    mutex    sync.RWMutex
}

type ServiceEntry struct {
    Name     string
    Instance interface{}
    Type     reflect.Type
}
```

### Controller Registry

```go
type ControllerRegistry struct {
    controllers []*ControllerEntry
    mutex       sync.RWMutex
}

type ControllerEntry struct {
    Path    string
    Methods map[string]*MethodEntry
}
```

### Guard

```go
type Guard interface {
    CanActivate(ctx echo.Context) (bool, error)
}
```

### Interceptor

```go
type Interceptor interface {
    Intercept(ctx echo.Context, next echo.HandlerFunc) error
}
```

### Exception Filter

```go
type ExceptionFilter interface {
    Catch(exception interface{}, ctx echo.Context) error
}
```

### Pipe

```go
type Pipe interface {
    Transform(value interface{}) (interface{}, error)
}
```

### Lifecycle Hook

```go
type LifecycleHook interface {
    OnApplicationStart(ctx context.Context) error
    OnApplicationStop(ctx context.Context) error
    OnModuleInit(ctx context.Context) error
    OnModuleDestroy(ctx context.Context) error
}
```

This documentation provides a comprehensive guide to using the GoNest framework. For more examples and advanced usage patterns, refer to the examples directory in the repository.
