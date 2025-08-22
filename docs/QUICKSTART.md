# GoNest Quick Start Guide

This guide will walk you through creating a new GoNest project from scratch, including folder setup, dependency initialization, and your first application.

## 🚀 Installation

To get started, you can either use our **CLI tool** (recommended), **automated setup scripts**, or manually set up your project. All approaches will produce the same outcome.

### 🚀 CLI Tool (Recommended)

Our CLI tool provides the fastest and most reliable way to create a new GoNest project:

```bash
# 1. Install GoNest CLI globally (Recommended)
$ git clone https://github.com/ulims/GoNest.git
$ cd GoNest
$ go install ./cmd/gonest

# 2. Verify installation
$ gonest --help

# 3. Create a new project
$ gonest new my-project-name

# 4. Create with specific template and strict mode
$ gonest new my-api --template=api --strict
```

The CLI tool automatically:
- ✅ Create the recommended project structure
- ✅ Initialize Go module and Git repository
- ✅ Install all GoNest dependencies
- ✅ Generate configuration files
- ✅ Set up Docker and build automation
- ✅ Create comprehensive documentation
- ✅ Support multiple project templates
- ✅ Generate components (modules, controllers, services)

> **HINT**  
> The CLI tool is the most reliable way to get started. It handles all dependencies and creates a production-ready project structure.

### 🚀 Automated Setup Scripts

If you prefer to use our setup scripts directly:

#### Linux/macOS
```bash
# Clone the GoNest repository
$ git clone https://github.com/ulims/GoNest.git
$ cd GoNest

# Make the script executable and run it
$ chmod +x scripts/setup-project.sh
$ ./scripts/setup-project.sh my-project-name
```

#### Windows
```cmd
# Clone the GoNest repository
$ git clone https://github.com/ulims/GoNest.git
$ cd GoNest

# Run the batch script
$ scripts\setup-project.bat
```

### Alternatives

#### Manual Setup
If you prefer to set up manually:

```bash
# Create a new directory for your project
$ mkdir my-gonest-app
$ cd my-gonest-app

# Initialize Go module
$ go mod init my-gonest-app

# Add GoNest dependency
$ go get github.com/ulims/GoNest
```

#### CLI Tool
GoNest includes a powerful CLI tool for project scaffolding and component generation:

```bash
# 1. Install GoNest CLI globally (Recommended)
$ git clone https://github.com/ulims/GoNest.git
$ cd GoNest
$ go install ./cmd/gonest

# 2. Verify installation
$ gonest --help

# 3. Create a new project
$ gonest new my-project-name

# 4. Create with specific template and strict mode
$ gonest new my-api --template=api --strict

# 5. Generate components in existing projects
$ gonest generate module user
$ gonest generate controller user
$ gonest generate service user
```

> **NOTE**  
> Unlike NestJS, GoNest focuses on providing powerful setup scripts and CLI tools rather than separate starter repositories. This approach gives you more control and keeps everything in one place.

## 🚀 Prerequisites

Before starting, ensure you have:

- **Go 1.21+** installed on your system
- **Git** for version control
- **A code editor** (VS Code, GoLand, Vim, etc.)
- **Basic Go knowledge** (packages, modules, structs)

### Verify Go Installation

```bash
go version
# Should output: go version go1.21.x windows/amd64 (or similar)
```

## 📁 Step 1: Create Project Structure

### 1.1 Create Project Directory

```bash
# Create a new directory for your project
mkdir my-gonest-app
cd my-gonest-app
```

### 1.2 Initialize Git Repository (Optional but Recommended)

```bash
# Initialize git repository
git init

# Create .gitignore file
echo "# GoNest Project
# Binaries
bin/
dist/
*.exe
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment files
.env
.env.local
.env.*.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Logs
*.log

# Build artifacts
build/
tmp/" > .gitignore
```

## 🔧 Step 2: Initialize Go Module

### 2.1 Initialize Go Module

```bash
# Initialize Go module
go mod init my-gonest-app
```

This creates a `go.mod` file that manages your project's dependencies.

### 2.2 Create Basic Project Structure

```bash
# Create directory structure
mkdir -p cmd/server
mkdir -p internal/modules
mkdir -p internal/config
mkdir -p internal/shared
mkdir -p pkg
mkdir -p docs
mkdir -p scripts
mkdir -p tests
```

Your project structure should now look like:

```
my-gonest-app/
├── cmd/
│   └── server/
├── internal/
│   ├── modules/
│   ├── config/
│   └── shared/
├── pkg/
├── docs/
├── scripts/
├── tests/
├── go.mod
└── .gitignore
```

## 📦 Step 3: Install GoNest Dependencies

### 3.1 Add GoNest Framework

```bash
# Add GoNest framework
go get github.com/ulims/GoNest
```

### 3.2 Add Additional Dependencies

```bash
# Add Echo framework (HTTP server)
go get github.com/labstack/echo/v4

# Add Logrus for structured logging
go get github.com/sirupsen/logrus

# Add Validator for request validation
go get github.com/go-playground/validator/v10

# Add JWT for authentication
go get github.com/golang-jwt/jwt/v5

# Add WebSocket support
go get github.com/gorilla/websocket

# Add Redis client (optional, for caching)
go get github.com/redis/go-redis/v9

# Add MongoDB driver (optional, for database)
go get go.mongodb.org/mongo-driver/mongo
```

### 3.3 Verify Dependencies

```bash
# Download and tidy dependencies
go mod tidy

# Check module status
go mod verify
```

## 🏗️ Step 4: Create Your First Application

### 4.1 Create Main Application File

Create `cmd/server/main.go`:

```go
package main

import (
    "context"
    "github.com/sirupsen/logrus"
    gonest "github.com/ulims/GoNest"
)

func main() {
    // Initialize logger
    logger := logrus.New()
    logger.SetLevel(logrus.InfoLevel)
    logger.SetFormatter(&logrus.JSONFormatter{})
    
    // Create application with configuration
    app := gonest.NewApplication().
        Config(&gonest.Config{
            Port:        "8080",
            Host:        "localhost",
            Environment: "development",
            LogLevel:    "info",
        }).
        Logger(logger).
        Build()
    
    // Register application lifecycle hooks
    app.LifecycleManager.RegisterHook(
        gonest.EventApplicationStart,
        gonest.LifecycleHookFunc(func(ctx context.Context) error {
            logger.Info("🚀 Application starting up...")
            logger.Info("📁 GoNest application initialized successfully")
            return nil
        }),
        gonest.PriorityHigh,
    )
    
    app.LifecycleManager.RegisterHook(
        gonest.EventApplicationStop,
        gonest.LifecycleHookFunc(func(ctx context.Context) error {
            logger.Info("🛑 Application shutting down...")
            return nil
        }),
        gonest.PriorityHigh,
    )
    
    // Start the application
    logger.Info("Starting GoNest application...")
    if err := app.Start(); err != nil {
        logger.Fatal("Failed to start application:", err)
    }
}
```

### 4.2 Create Configuration File

Create `internal/config/config.go`:

```go
package config

import (
    "os"
    "strconv"
)

// Config holds application configuration
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
    Redis    RedisConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
    Port        string
    Host        string
    Environment string
    LogLevel    string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
    Host     string
    Port     int
    Name     string
    Username string
    Password string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
    Secret     string
    Expiration int64
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
    return &Config{
        Server: ServerConfig{
            Port:        getEnv("PORT", "8080"),
            Host:        getEnv("HOST", "localhost"),
            Environment: getEnv("ENV", "development"),
            LogLevel:    getEnv("LOG_LEVEL", "info"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 27017),
            Name:     getEnv("DB_NAME", "gonest"),
            Username: getEnv("DB_USERNAME", ""),
            Password: getEnv("DB_PASSWORD", ""),
        },
        JWT: JWTConfig{
            Secret:     getEnv("JWT_SECRET", "your-secret-key"),
            Expiration: getEnvAsInt64("JWT_EXPIRATION", 86400), // 24 hours
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvAsInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvAsInt("REDIS_DB", 0),
        },
    }
}

// Helper functions
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
            return intValue
        }
    }
    return defaultValue
}
```

## 🎯 Step 5: Create Your First Module

### 5.1 Create User Module

Create `internal/modules/user/user_module.go`:

```go
package user

import (
    "github.com/sirupsen/logrus"
    gonest "github.com/ulims/GoNest"
)

// UserModule represents the user module
type UserModule struct {
    *gonest.Module
}

// NewUserModule creates a new user module
func NewUserModule(logger *logrus.Logger) *UserModule {
    // Create services
    userService := NewUserService(logger)
    
    // Create controllers
    userController := NewUserController(userService)
    
    // Create and return module
    module := gonest.NewModule("UserModule").
        Controller(userController).
        Service(userService).
        Build()
    
    return &UserModule{
        Module: module,
    }
}
```

### 5.2 Create User Service

Create `internal/modules/user/user_service.go`:

```go
package user

import (
    "errors"
    "sync"
    "time"
    "github.com/sirupsen/logrus"
)

// User represents a user entity
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

// UserService handles user business logic
type UserService struct {
    users  map[string]*User
    logger *logrus.Logger
    mutex  sync.RWMutex
}

// NewUserService creates a new user service
func NewUserService(logger *logrus.Logger) *UserService {
    return &UserService{
        users:  make(map[string]*User),
        logger: logger,
    }
}

// CreateUser creates a new user
func (s *UserService) CreateUser(username, email, password, firstName, lastName string) (*User, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    // Check if username already exists
    if s.usernameExists(username) {
        return nil, errors.New("username already exists")
    }
    
    // Check if email already exists
    if s.emailExists(email) {
        return nil, errors.New("email already exists")
    }
    
    // Create new user
    user := &User{
        ID:        generateID(),
        Username:  username,
        Email:     email,
        Password:  password, // In real app, hash the password
        FirstName: firstName,
        LastName:  lastName,
        Role:      "user",
        Status:    "active",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Store user
    s.users[user.ID] = user
    
    s.logger.Infof("Created user: %s (%s)", user.Username, user.ID)
    return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*User, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    user, exists := s.users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    
    return user, nil
}

// ListUsers retrieves all users
func (s *UserService) ListUsers() ([]*User, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    users := make([]*User, 0, len(s.users))
    for _, user := range s.users {
        users = append(users, user)
    }
    
    return users, nil
}

// Helper methods
func (s *UserService) usernameExists(username string) bool {
    for _, user := range s.users {
        if user.Username == username {
            return true
        }
    }
    return false
}

func (s *UserService) emailExists(email string) bool {
    for _, user := range s.users {
        if user.Email == email {
            return true
        }
    }
    return false
}

// generateID generates a unique ID (simplified for example)
func generateID() string {
    return time.Now().Format("20060102150405")
}
```

### 5.3 Create User Controller

Create `internal/modules/user/user_controller.go`:

```go
package user

import (
    "net/http"
    "github.com/labstack/echo/v4"
    gonest "github.com/ulims/GoNest"
)

// UserController handles HTTP requests for user operations
type UserController struct {
    userService *UserService
}

// NewUserController creates a new user controller
func NewUserController(userService *UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

// CreateUser handles user creation
func (c *UserController) CreateUser(ctx echo.Context) error {
    var req struct {
        Username  string `json:"username" validate:"required,min=3,max=50"`
        Email     string `json:"email" validate:"required,email"`
        Password  string `json:"password" validate:"required,min=8"`
        FirstName string `json:"first_name" validate:"required,min=2,max=50"`
        LastName  string `json:"last_name" validate:"required,min=2,max=50"`
    }
    
    if err := ctx.Bind(&req); err != nil {
        return gonest.BadRequestException("Invalid request body")
    }
    
    // Validate request
    if err := gonest.ValidateStruct(req, nil); err != nil {
        return gonest.BadRequestException(err.Error())
    }
    
    // Create user
    user, err := c.userService.CreateUser(
        req.Username,
        req.Email,
        req.Password,
        req.FirstName,
        req.LastName,
    )
    if err != nil {
        return gonest.BadRequestException(err.Error())
    }
    
    return ctx.JSON(http.StatusCreated, user)
}

// GetUser handles user retrieval by ID
func (c *UserController) GetUser(ctx echo.Context) error {
    id := ctx.Param("id")
    if id == "" {
        return gonest.BadRequestException("User ID is required")
    }
    
    user, err := c.userService.GetUser(id)
    if err != nil {
        return gonest.NotFoundException("User not found")
    }
    
    return ctx.JSON(http.StatusOK, user)
}

// ListUsers handles user listing
func (c *UserController) ListUsers(ctx echo.Context) error {
    users, err := c.userService.ListUsers()
    if err != nil {
        return gonest.BadRequestException("Failed to retrieve users")
    }
    
    return ctx.JSON(http.StatusOK, users)
}
```

## 🔗 Step 6: Integrate Module with Application

### 6.1 Update Main Application

Update `cmd/server/main.go` to include the user module:

```go
package main

import (
    "context"
    "github.com/sirupsen/logrus"
    gonest "github.com/ulims/GoNest"
    "my-gonest-app/internal/modules/user"
)

func main() {
    // Initialize logger
    logger := logrus.New()
    logger.SetLevel(logrus.InfoLevel)
    logger.SetFormatter(&logrus.JSONFormatter{})
    
    // Create application with configuration
    app := gonest.NewApplication().
        Config(&gonest.Config{
            Port:        "8080",
            Host:        "localhost",
            Environment: "development",
            LogLevel:    "info",
        }).
        Logger(logger).
        Build()
    
    // Create and register user module
    userModule := user.NewUserModule(logger)
    app.ModuleRegistry.Register(userModule.Module)
    
    // Register application lifecycle hooks
    app.LifecycleManager.RegisterHook(
        gonest.EventApplicationStart,
        gonest.LifecycleHookFunc(func(ctx context.Context) error {
            logger.Info("🚀 Application starting up...")
            logger.Info("📁 User module registered successfully")
            return nil
        }),
        gonest.PriorityHigh,
    )
    
    app.LifecycleManager.RegisterHook(
        gonest.EventApplicationStop,
        gonest.LifecycleHookFunc(func(ctx context.Context) error {
            logger.Info("🛑 Application shutting down...")
            return nil
        }),
        gonest.PriorityHigh,
    )
    
    // Start the application
    logger.Info("Starting GoNest application...")
    if err := app.Start(); err != nil {
        logger.Fatal("Failed to start application:", err)
    }
}
```

## 🚀 Step 7: Run Your Application

### 7.1 Build the Application

```bash
# Build the application
go build -o bin/server cmd/server/main.go
```

### 7.2 Run the Application

```bash
# Run the application
./bin/server
```

Or run directly:

```bash
# Run directly with go run
go run cmd/server/main.go
```

### 7.3 Test Your Application

Your application should now be running at `http://localhost:8080` with these endpoints:

- `POST /users` - Create a new user
- `GET /users/:id` - Get user by ID
- `GET /users` - List all users

## 🧪 Step 8: Test Your API

### 8.1 Create a User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### 8.2 Get a User

```bash
# Replace USER_ID with the actual ID from the create response
curl http://localhost:8080/users/USER_ID
```

### 8.3 List All Users

```bash
curl http://localhost:8080/users
```

## 📁 Final Project Structure

After completing all steps, your project should look like this:

```
my-gonest-app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── modules/
│   │   └── user/
│   │       ├── user_module.go
│   │       ├── user_service.go
│   │       └── user_controller.go
│   ├── config/
│   │   └── config.go
│   └── shared/
├── pkg/
├── docs/
├── scripts/
├── tests/
├── bin/
│   └── server
├── go.mod
├── go.sum
└── .gitignore
```

## 🔧 Step 9: Environment Configuration

### 9.1 Create Environment File

Create `.env` file:

```env
# Server Configuration
PORT=8080
HOST=localhost
ENV=development
LOG_LEVEL=info

# Database Configuration
DB_HOST=localhost
DB_PORT=27017
DB_NAME=gonest
DB_USERNAME=
DB_PASSWORD=

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION=86400

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 9.2 Add Environment Loading

Install environment loading package:

```bash
go get github.com/joho/godotenv
```

Update your main.go to load environment variables:

```go
package main

import (
    "context"
    "log"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    gonest "github.com/ulims/GoNest"
    "my-gonest-app/internal/config"
    "my-gonest-app/internal/modules/user"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    // Load configuration
    cfg := config.LoadConfig()
    
    // Initialize logger
    logger := logrus.New()
    logger.SetLevel(logrus.InfoLevel)
    logger.SetFormatter(&logrus.JSONFormatter{})
    
    // Create application with configuration
    app := gonest.NewApplication().
        Config(&gonest.Config{
            Port:        cfg.Server.Port,
            Host:        cfg.Server.Host,
            Environment: cfg.Server.Environment,
            LogLevel:    cfg.Server.LogLevel,
        }).
        Logger(logger).
        Build()
    
    // Create and register user module
    userModule := user.NewUserModule(logger)
    app.ModuleRegistry.Register(userModule.Module)
    
    // Register application lifecycle hooks
    app.LifecycleManager.RegisterHook(
        gonest.EventApplicationStart,
        gonest.LifecycleHookFunc(func(ctx context.Context) error {
            logger.Info("🚀 Application starting up...")
            logger.Info("📁 User module registered successfully")
            logger.Infof("🌍 Server running on %s:%s", cfg.Server.Host, cfg.Server.Port)
            return nil
        }),
        gonest.PriorityHigh,
    )
    
    app.LifecycleManager.RegisterHook(
        gonest.EventApplicationStop,
        gonest.LifecycleHookFunc(func(ctx context.Context) error {
            logger.Info("🛑 Application shutting down...")
            return nil
        }),
        gonest.PriorityHigh,
    )
    
    // Start the application
    logger.Info("Starting GoNest application...")
    if err := app.Start(); err != nil {
        logger.Fatal("Failed to start application:", err)
    }
}
```

## 🎯 Next Steps

Now that you have a basic GoNest application running, you can:

1. **Add More Modules**: Create additional modules for different business domains
2. **Database Integration**: Add MongoDB or PostgreSQL integration
3. **Authentication**: Implement JWT-based authentication
4. **Validation**: Add more comprehensive request validation
5. **Testing**: Write unit and integration tests
6. **Documentation**: Add API documentation with Swagger
7. **Deployment**: Prepare your application for production deployment

## 🚨 Troubleshooting

### Common Issues

1. **Module Not Found**: Ensure your import paths are correct
2. **Port Already in Use**: Change the port in your configuration
3. **Dependency Issues**: Run `go mod tidy` to resolve dependency conflicts
4. **Build Errors**: Check that all files are in the correct packages

### Getting Help

- Check the [Full Documentation](DOCUMENTATION.md)
- Review the [Architecture Guide](ARCHITECTURE.md)
- Look at the [Examples](examples/) directory
- Open an [Issue](https://github.com/ulims/GoNest/issues) on GitHub

---

**Congratulations! 🎉 You've successfully created your first GoNest application. You now have a solid foundation to build upon!**

