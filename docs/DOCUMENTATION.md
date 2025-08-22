# GoNest Framework Documentation

Welcome to the comprehensive documentation for GoNest, a powerful, enterprise-grade Go web framework inspired by NestJS. This documentation will guide you through all aspects of the framework, from basic concepts to advanced features.

## 📚 Table of Contents

1. [🚀 Quick Start](QUICKSTART.md) - Get up and running in minutes
2. [🏗️ Architecture Guide](ARCHITECTURE.md) - Deep dive into architectural patterns
3. [📋 Features Overview](#features) - Complete feature reference
4. [🎯 Examples](examples/) - Working examples and tutorials
5. [🔧 API Reference](#api-reference) - Complete API documentation
6. [🧪 Testing Guide](#testing-guide) - Testing strategies and utilities
7. [🚀 Deployment Guide](#deployment-guide) - Production deployment
8. [🔒 Security Guide](#security-guide) - Security best practices

## 🚀 Getting Started

### Quick Start

For developers who want to jump right in, start with our [Quick Start Guide](QUICKSTART.md) which covers:

- Creating a new project from scratch
- Setting up the project structure
- Installing GoNest dependencies
- Building your first application
- Creating modules, services, and controllers
- Testing your API endpoints

### Architecture Example

The `examples/architecture/` demonstrates the recommended **NestJS-style modular structure**:

```
examples/architecture/
├── main.go                 # Application entry point
├── main_module.go          # Root module that imports feature modules
├── modules/                # Feature modules directory
│   └── user/              # User feature module
│       ├── user_module.go    # Module definition and registration
│       ├── user_service.go   # Business logic layer
│       └── user_controller.go # HTTP request handling
└── README.md              # Module documentation
```

**Key Benefits of the Flat Module Structure:**
- **Simpler imports**: No deep nesting
- **Easier navigation**: All related files in one directory
- **Better Go conventions**: Follows Go package organization
- **Reduced complexity**: Fewer directory levels to manage

## 🏗️ Core Concepts

### 1. **Application (`gonest.Application`)**
The root container that orchestrates the entire application:

```go
app := gonest.NewApplication().
    Config(&gonest.Config{
        Port:        "8080",
        Host:        "localhost",
        Environment: "development",
        LogLevel:    "info",
    }).
    Logger(logger).
    Build()
```

**Responsibilities:**
- Module registration and management
- Lifecycle management
- Configuration management
- HTTP server setup
- Graceful shutdown

### 2. **Module System (`gonest.Module`)**
Modules are the building blocks of GoNest applications:

```go
module := gonest.NewModule("UserModule").
    Controller(userController).
    Service(userService).
    Build()
```

**Key Features:**
- **Self-contained**: Each module manages its own dependencies
- **Importable**: Modules can import other modules
- **Configurable**: Module-specific settings and options
- **Lifecycle-aware**: Module initialization and cleanup hooks

### 3. **Controller Layer**
Controllers handle HTTP requests and responses:

```go
type UserController struct {
    userService *UserService
}

func (c *UserController) CreateUser(ctx echo.Context) error {
    // Handle HTTP request
    // Validate input
    // Call service
    // Return response
}
```

**Responsibilities:**
- HTTP request parsing
- Input validation
- Service coordination
- Response formatting
- Error handling

### 4. **Service Layer**
Services contain business logic and data operations:

```go
type UserService struct {
    users  map[string]*User
    logger *logrus.Logger
    mutex  sync.RWMutex
}

func (s *UserService) CreateUser(username, email, password, firstName, lastName string) (*User, error) {
    // Business logic implementation
}
```

**Responsibilities:**
- Business logic implementation
- Data validation
- External service integration
- Transaction management
- Error handling

### 5. **Model Layer**
Models represent domain entities and business rules:

```go
type User struct {
    ID        string    `json:"id" validate:"required"`
    Username  string    `json:"username" validate:"required,min=3,max=50"`
    Email     string    `json:"email" validate:"required,email"`
    Password  string    `json:"password,omitempty" validate:"required,min=8"`
    FirstName string    `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string    `json:"last_name" validate:"required,min=2,max=50"`
    Role      string    `json:"role" validate:"required,oneof=user admin moderator"`
    Status    string    `json:"status" validate:"required,oneof=active inactive suspended"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Features:**
- Validation tags for data integrity
- JSON serialization support
- Business rule encapsulation
- Immutable design patterns

## 📁 Recommended Project Structure

GoNest follows a **flat module structure** for simplicity and clarity:

```
my-gonest-app/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── modules/                 # Feature modules
│   │   ├── user/               # User module
│   │   │   ├── user_module.go    # Module definition
│   │   │   ├── user_service.go   # Business logic
│   │   │   └── user_controller.go # HTTP handlers
│   │   ├── auth/               # Auth module
│   │   │   ├── auth_module.go    # Module definition
│   │   │   ├── auth_service.go   # Business logic
│   │   │   └── auth_controller.go # HTTP handlers
│   │   └── product/            # Product module
│   │       ├── product_module.go   # Module definition
│   │       ├── product_service.go  # Business logic
│   │       └── product_controller.go # HTTP handlers
│   ├── config/
│   │   └── config.go           # Configuration management
│   └── shared/
│       ├── middleware/          # Shared middleware
│       ├── utils/              # Utility functions
│       └── constants/          # Application constants
├── pkg/                        # Public packages
├── docs/                       # Documentation
├── scripts/                    # Build and deployment scripts
├── tests/                      # Integration tests
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── .env                        # Environment variables
└── .gitignore                 # Git ignore file
```

### Module Structure
Each module follows this pattern:

```
modules/{feature}/
├── {feature}_module.go     # Module definition
├── {feature}_service.go    # Business logic
└── {feature}_controller.go # HTTP handlers
```

## 🔄 Data Flow

### Request Flow
```
HTTP Request → Controller → Service → Model
     ↑           ↓          ↓        ↓
Response ← Controller ← Service ← Model
```

### Detailed Flow
1. **HTTP Request**: Echo framework receives the request
2. **Route Matching**: Controller method is identified
3. **Middleware Execution**: Guards, interceptors, pipes run
4. **Controller Method**: Request is processed
5. **Service Call**: Business logic is executed
6. **Data Operation**: Models are manipulated
7. **Response**: JSON response is formatted and returned

## 🎯 Layer Responsibilities

### Controller Layer
- **Input Validation**: Request data validation
- **Authentication**: User identity verification
- **Authorization**: Permission checking
- **Request Routing**: Endpoint mapping
- **Response Formatting**: HTTP response creation

### Service Layer
- **Business Logic**: Core application rules
- **Data Validation**: Business rule validation
- **External Integration**: API calls, database operations
- **Transaction Management**: Data consistency
- **Error Handling**: Business logic errors

### Model Layer
- **Data Structure**: Entity representation
- **Validation Rules**: Data integrity constraints
- **Business Methods**: Entity-specific operations
- **Serialization**: JSON/XML conversion
- **Relationships**: Entity associations

## 🔐 DTO Pattern

Data Transfer Objects (DTOs) are used for input/output validation:

```go
// Request DTO (anonymous struct for simplicity)
var req struct {
    Username  string `json:"username" validate:"required,min=3,max=50"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"first_name" validate:"required,min=2,max=50"`
    LastName  string `json:"last_name" validate:"required,min=2,max=50"`
}

// Validation
if err := gonest.ValidateStruct(req, nil); err != nil {
    return gonest.BadRequestException(err.Error())
}
```

**Benefits:**
- **Input Validation**: Ensures data integrity
- **API Contracts**: Clear input/output specifications
- **Security**: Prevents malicious input
- **Documentation**: Self-documenting API structure

## 🛡️ Guards & Interceptors

### Guards
Guards run before controller methods and can:

```go
type AuthGuard struct{}

func (g *AuthGuard) CanActivate(ctx echo.Context) bool {
    // Check authentication
    // Verify permissions
    // Return true/false
}
```

### Interceptors
Interceptors wrap controller methods for:

```go
type LoggingInterceptor struct{}

func (i *LoggingInterceptor) Intercept(ctx echo.Context, next echo.HandlerFunc) error {
    start := time.Now()
    err := next(ctx)
    duration := time.Since(start)
    
    // Log request details
    return err
}
```

## 🔄 Lifecycle Management

GoNest provides comprehensive lifecycle hooks:

```go
// Application lifecycle
app.LifecycleManager.RegisterHook(
    gonest.EventApplicationStart,
    gonest.LifecycleHookFunc(func(ctx context.Context) error {
        logger.Info("🚀 Application starting up...")
        return nil
    }),
    gonest.PriorityHigh,
)

// Module lifecycle
type UserModule struct {
    *gonest.Module
}

func (m *UserModule) OnModuleInit() error {
    // Module initialization logic
    return nil
}

func (m *UserModule) OnModuleDestroy() error {
    // Module cleanup logic
    return nil
}
```

## 📈 Scaling Patterns

### Horizontal Scaling
- **Stateless Services**: No shared state between instances
- **Load Balancing**: Multiple application instances
- **Database Sharding**: Distributed data storage
- **Microservices**: Independent service deployment

### Vertical Scaling
- **Resource Optimization**: CPU and memory optimization
- **Connection Pooling**: Database connection management
- **Caching**: In-memory data storage
- **Async Processing**: Background job processing

## 🧪 Testing Strategy

### Unit Testing
```go
func TestUserService_CreateUser(t *testing.T) {
    service := NewUserService(logger)
    user, err := service.CreateUser("testuser", "test@example.com", "password", "John", "Doe")
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "testuser", user.Username)
}
```

### Integration Testing
```go
func TestUserController_CreateUser(t *testing.T) {
    // Setup test application
    // Make HTTP request
    // Verify response
}
```

### Test Utilities
GoNest provides `TestApp` for testing:

```go
testApp := gonest.NewTestApp().
    Module(userModule).
    Build()

// Test with real HTTP requests
```

## 🔒 Security Considerations

### Authentication
- **JWT Tokens**: Secure token-based authentication
- **Password Hashing**: Secure password storage
- **Session Management**: User session handling

### Authorization
- **Role-Based Access Control**: User permission management
- **Resource-Level Security**: Fine-grained access control
- **API Security**: Rate limiting and throttling

### Data Protection
- **Input Validation**: Prevent injection attacks
- **Output Encoding**: Prevent XSS attacks
- **HTTPS**: Secure communication

## 📊 Monitoring & Observability

### Logging
- **Structured Logging**: JSON-formatted logs
- **Log Levels**: Debug, Info, Warn, Error
- **Context Information**: Request tracing

### Metrics
- **Performance Metrics**: Response times, throughput
- **Business Metrics**: User actions, system usage
- **Health Checks**: System status monitoring

### Tracing
- **Request Tracing**: End-to-end request tracking
- **Performance Profiling**: Bottleneck identification
- **Error Tracking**: Exception monitoring

## 🚀 Deployment Considerations

### Containerization
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Environment Configuration
```go
type Config struct {
    Port        string `env:"PORT" envDefault:"8080"`
    Host        string `env:"HOST" envDefault:"localhost"`
    Environment string `env:"ENV" envDefault:"development"`
    LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
}
```

### Health Checks
```go
func (c *HealthController) Health(ctx echo.Context) error {
    return ctx.JSON(http.StatusOK, map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now(),
        "version": "1.0.0",
    })
}
```

## 🎯 Best Practices

### 1. **Module Design**
- Keep modules focused and single-purpose
- Minimize inter-module dependencies
- Use clear naming conventions
- Document module responsibilities

### 2. **Service Layer**
- Implement business logic, not HTTP concerns
- Use dependency injection for external services
- Handle errors gracefully
- Implement proper logging

### 3. **Controller Design**
- Keep controllers thin
- Focus on HTTP concerns only
- Use consistent error responses
- Implement proper validation

### 4. **Error Handling**
- Use appropriate HTTP status codes
- Provide meaningful error messages
- Log errors with context
- Implement error recovery strategies

### 5. **Performance**
- Use connection pooling
- Implement caching strategies
- Optimize database queries
- Use async processing where appropriate

## 🔮 Future Enhancements

### Planned Features
- **GraphQL Support**: Alternative to REST APIs
- **gRPC Integration**: High-performance RPC
- **Event Sourcing**: Event-driven architecture
- **CQRS**: Command Query Responsibility Segregation

### Community Contributions
- **Plugin System**: Third-party extensions
- **Template Engine**: Code generation
- **Migration Tools**: Database schema management
- **Monitoring Dashboard**: Built-in observability

## 📚 Additional Resources

### Examples
- **[Architecture Example](examples/architecture/)** - NestJS-style modular structure
- **[Advanced Example](examples/advanced/)** - Authentication, validation, and dependency injection
- **[MongoDB Example](examples/mongodb/)** - Database integration

### Guides
- **[Quick Start Guide](QUICKSTART.md)** - Step-by-step project setup
- **[Architecture Guide](ARCHITECTURE.md)** - Detailed architectural patterns
- **[Features Overview](#features)** - All available features

## 📚 Conclusion

The GoNest architecture provides a solid foundation for building enterprise-grade applications. By following NestJS patterns and Go best practices, developers can create:

- **Maintainable Code**: Clear structure and separation of concerns
- **Scalable Applications**: Support for horizontal and vertical scaling
- **Testable Systems**: Easy mocking and testing capabilities
- **Secure Applications**: Built-in security features and best practices
- **Performant Services**: Optimized for high throughput and low latency

This architecture enables teams to build complex applications while maintaining code quality and developer productivity. The modular design makes it easy to add new features, refactor existing code, and onboard new team members.

---

*For more information, examples, and community support, visit the GoNest framework documentation and GitHub repository.*

