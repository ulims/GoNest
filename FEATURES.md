# GoNest Framework - Advanced Features

GoNest is a powerful Go framework inspired by NestJS, providing enterprise-grade features for building scalable web applications and APIs.

## 🚀 Key Features

### 1. WebSocket Support
Real-time communication with Gateway pattern similar to NestJS.

```go
// WebSocket Gateway
type ChatGateway struct {
    *gonest.BaseWebSocketGateway
}

func (cg *ChatGateway) OnConnection(client *gonest.WebSocketClient) error {
    client.JoinRoom("general")
    return client.Emit("welcome", map[string]interface{}{
        "message": "Welcome to the chat!",
    })
}

func (cg *ChatGateway) OnMessage(client *gonest.WebSocketClient, message *gonest.WebSocketMessage) error {
    // Handle different message types
    switch message.Event {
    case "chat_message":
        // Broadcast to room
        return nil
    }
    return nil
}

// Register gateway
wsServer := gonest.NewWebSocketServer(config, logger)
wsServer.RegisterGateway(chatGateway)
app.GetEcho().GET("/ws/chat", wsServer.HandleConnection("/chat"))
```

**Features:**
- Gateway pattern with lifecycle hooks
- Room-based messaging
- Authentication middleware
- Auto-reconnection support
- Message filtering and validation

### 2. JWT Authentication & Authorization
Comprehensive authentication system with passport-like strategies.

```go
// JWT Configuration
authConfig := gonest.DefaultJWTConfig()
authConfig.SecretKey = "your-secret-key"
authConfig.TokenExpiry = 24 * time.Hour
authService := gonest.NewAuthService(authConfig, logger)

// Local Strategy
passportService := gonest.NewPassportService(logger)
localStrategy := gonest.NewLocalStrategy(func(username, password string) (*gonest.AuthUser, error) {
    // Validate credentials
    return &gonest.AuthUser{
        ID:       "user-id",
        Username: username,
        Roles:    []string{"user"},
    }, nil
})
passportService.Use(localStrategy)

// Protect routes
apiGroup.Use(authService.JWTMiddleware())

// Role-based authorization
adminGroup.Use(authService.RequireRoles("admin"))
```

**Features:**
- JWT token generation and validation
- Multiple authentication strategies
- Role-based access control (RBAC)
- Guard system for route protection
- Token refresh mechanism
- Custom authentication providers

### 3. Advanced Caching System
Multi-level caching with decorators and interceptors.

```go
// Memory Cache Provider
cacheProvider := gonest.NewMemoryCache(logger)
cacheService := gonest.NewCacheService(cacheProvider, logger)

// Cache data
cacheService.Set(ctx, "user:123", user, 10*time.Minute)

// Get or set pattern
var user User
err := cacheService.GetOrSet(ctx, "user:123", &user, func() (interface{}, error) {
    return getUserFromDB("123")
}, 10*time.Minute)

// HTTP Response Caching
cacheInterceptor := gonest.NewCacheInterceptor(cacheService, logger)
cacheConfig := &gonest.CacheConfig{
    TTL: 5 * time.Minute,
    KeyBuilder: func(c echo.Context) string {
        return fmt.Sprintf("api:%s:%s", c.Request().Method, c.Request().URL.Path)
    },
}
app.Use(cacheInterceptor.Middleware(cacheConfig))
```

**Features:**
- Multiple cache providers (Memory, Redis-compatible)
- HTTP response caching middleware
- Cache decorators for functions
- Conditional caching
- Cache invalidation strategies
- TTL support

### 4. Event System
Powerful event-driven architecture for decoupled communication.

```go
// Create event system
eventBus := gonest.NewEventBus(logger)
eventService := gonest.NewEventService(eventBus, "app", logger)

// Emit events
eventService.Emit(ctx, gonest.EventUserCreated, gonest.UserEventData{
    UserID:   user.ID,
    Username: user.Username,
    Action:   "create",
})

// Listen to events
eventService.On(gonest.EventUserCreated, func(ctx context.Context, event *gonest.Event) error {
    // Send welcome email, update analytics, etc.
    logger.Infof("User created: %s", event.Data.(gonest.UserEventData).Username)
    return nil
}, gonest.EventListenerConfig{
    Async:    true,
    Priority: 1,
    Retry:    gonest.DefaultRetryConfig(),
})
```

**Features:**
- Event emitter and listener pattern
- Asynchronous event processing
- Event priority and filtering
- Retry mechanisms with backoff
- Pattern-based event matching
- Namespace isolation

### 5. Configuration Management
Environment-based configuration with multiple providers.

```go
// Configuration structure
type AppConfig struct {
    Server struct {
        Port int    `config:"server.port,default=8080"`
        Host string `config:"server.host,default=localhost"`
    }
    Database struct {
        URI string `config:"database.uri,required"`
    }
}

// Load configuration
configService, err := gonest.LoadEnvironmentConfig("development", logger)

// Inject configuration
var config AppConfig
gonest.InjectConfig(configService, &config)

// Add custom providers
configService.AddProvider(gonest.NewFileConfigProvider("config.yaml"))
configService.AddProvider(gonest.NewEnvironmentConfigProvider("APP"))
```

**Features:**
- Multiple configuration providers (ENV, YAML, JSON)
- Environment-specific configs
- Configuration validation
- Hot reload support
- Nested configuration structures
- Default values and required fields

### 6. Rate Limiting
Advanced rate limiting with multiple strategies.

```go
// Different rate limiting strategies
generalLimit := gonest.PerMinute(100).WithStrategy(gonest.FixedWindowStrategy)
burstLimit := gonest.PerSecond(10).WithStrategy(gonest.TokenBucketStrategy)
userLimit := gonest.PerHour(1000).WithKeyGenerator(gonest.UserKeyGenerator)

// Apply rate limiting
app.Use(generalLimit.Middleware(logger))

// Custom key generators
customKeyGen := func(c echo.Context) string {
    user, _ := gonest.GetCurrentUser(c)
    return fmt.Sprintf("user:%s:api", user.ID)
}
```

**Features:**
- Multiple rate limiting algorithms (Fixed Window, Sliding Window, Token Bucket, Leaky Bucket)
- Flexible key generation strategies
- Per-user, per-IP, per-route limiting
- Burst limiting support
- Custom storage backends
- Rate limit headers

### 7. MongoDB Integration
Mongoose-like ODM for MongoDB with schema validation.

```go
// Define schema
userSchema := gonest.NewSchema()
userSchema.Field("username", gonest.String).Required().MinLength(3).MaxLength(50)
userSchema.Field("email", gonest.String).Required().Unique().Index()
userSchema.Field("age", gonest.Number).Min(13).Max(120)
userSchema.Timestamps().Collection("users")

// Create model
userModel := mongoService.Model("User", userSchema)

// Document with lifecycle hooks
type User struct {
    gonest.MongoDBBaseModel
    Username string `bson:"username" json:"username"`
    Email    string `bson:"email" json:"email"`
}

func (u *User) BeforeSave() error {
    // Validation, transformation, etc.
    return nil
}

// Query builder
query := gonest.NewMongoDBQuery().
    Where("age", gonest.Gte, 18).
    Sort("createdAt", -1).
    Limit(10)
```

**Features:**
- Schema definition with validation
- Model lifecycle hooks
- Query builder with method chaining
- Index management
- Population and relationships
- Aggregation pipeline support

### 8. Testing Utilities
Comprehensive testing framework for unit and integration tests.

```go
func TestUserController(t *testing.T) {
    // Create test application
    testApp := gonest.NewTestApp(t).
        WithModule(userModule).
        Start(t)
    defer testApp.Stop()

    // Test HTTP endpoints
    response := testApp.POST("/api/users").
        WithJSON(map[string]interface{}{
            "username": "testuser",
            "email":    "test@example.com",
            "age":      25,
        }).
        Send(t)

    response.ExpectStatus(http.StatusCreated).
        ExpectJSONField("username", "testuser").
        ExpectJSONField("email", "test@example.com")

    // Test with authentication
    token := "jwt-token"
    response = testApp.GET("/api/users/profile").
        WithAuth(token).
        Send(t)

    response.ExpectOK().
        ExpectJSONFieldExists("id")
}
```

**Features:**
- HTTP testing utilities
- Request/response builders
- Assertion helpers
- Mock services
- Test fixtures
- Database test utilities

## 🎯 Architectural Patterns

### Module System
```go
userModule := gonest.NewModule("UserModule").
    Controllers(NewUserController()).
    Services(NewUserService()).
    Providers(map[string]interface{}{
        "UserRepository": userRepository,
        "EmailService":   emailService,
    }).
    Imports(authModule, cacheModule)
```

### Dependency Injection
```go
type UserService struct {
    userRepo    UserRepository    `inject:"UserRepository"`
    emailService EmailService     `inject:"EmailService"`
    logger      *logrus.Logger    `inject:"Logger"`
}
```

### Middleware Pipeline
```go
app.Use(rateLimitMiddleware)
app.Use(authMiddleware)
app.Use(cacheMiddleware)
app.Use(validationMiddleware)
```

### Guards and Interceptors
```go
// Guards for authorization
adminGuard := gonest.NewRoleGuard("admin")
userGroup.Use(gonest.UseGuards(adminGuard))

// Interceptors for cross-cutting concerns
loggingInterceptor := gonest.NewLoggingInterceptor(logger)
app.Use(loggingInterceptor.Middleware())
```

## 📦 Complete Example

See `examples/advanced/main.go` for a comprehensive example that demonstrates:
- WebSocket chat system
- JWT authentication
- Rate limiting
- Caching
- Event system
- MongoDB integration
- Configuration management

## 🔧 Configuration

### Environment Variables
```bash
GONEST_SERVER_PORT=8080
GONEST_DATABASE_MONGODB_URI=mongodb://localhost:27017
GONEST_AUTH_JWT_SECRET=your-secret-key
GONEST_CACHE_TTL=10m
```

### Configuration Files
```yaml
# config.yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "gonest_app"

auth:
  jwt:
    secret: "your-secret-key"
    expiry: "24h"

cache:
  ttl: "10m"
```

## 🚀 Getting Started

1. Install dependencies:
```bash
go mod tidy
```

2. Run the advanced example:
```bash
go run examples/advanced/main.go
```

3. Test endpoints:
```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# Create user (requires auth)
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","age":30}'

# WebSocket connection
ws://localhost:8080/ws/chat
```

## 📊 Performance Features

- **Zero-allocation routing** with Echo v4
- **Connection pooling** for databases
- **Response compression** and caching
- **Graceful shutdown** handling
- **Memory-efficient** WebSocket management
- **Optimized JSON** serialization

## 🔒 Security Features

- **JWT token validation**
- **Rate limiting** protection
- **CORS** configuration
- **Input validation** and sanitization
- **SQL injection** protection
- **XSS prevention**

GoNest provides all the features you need to build production-ready applications that can compete with NestJS and other modern frameworks!
