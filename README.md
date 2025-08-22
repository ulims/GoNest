# GoNest Framework

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-username/gonest)](https://goreportcard.com/report/github.com/your-username/gonest)
[![GoDoc](https://godoc.org/github.com/your-username/gonest?status.svg)](https://godoc.org/github.com/your-username/gonest)

A powerful, scalable Go framework inspired by NestJS that provides a complete solution for building enterprise-grade applications. GoNest combines the performance of Go with the elegant architecture patterns of NestJS.

## ✨ Features

- **🏗️ Modular Architecture**: Organize code into modules with clear boundaries
- **🔧 Dependency Injection**: Automatic dependency resolution and injection
- **🛡️ Guards & Interceptors**: Route-level authentication and request/response transformation
- **📝 Data Validation**: Comprehensive validation using struct tags and DTOs
- **🔐 Authentication**: JWT-based authentication with refresh tokens
- **⚡ WebSockets**: Real-time communication support
- **💾 Caching**: Built-in caching with multiple providers
- **📡 Events**: Event-driven architecture with async processing
- **🚦 Rate Limiting**: Protect your APIs with multiple strategies
- **🗄️ MongoDB Integration**: Simplified database operations
- **🧪 Testing**: Comprehensive testing utilities
- **⚙️ Configuration**: Environment-based configuration management

## 🚀 Quick Start

### Installation

```bash
go get github.com/ulims/gonest
```

### Basic Example

```go
package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/sirupsen/logrus"
    gonest "github.com/ulims/gonest"
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

    return ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) GetUser(ctx echo.Context) error {
    id := ctx.Param("id")
    user, err := c.userService.GetUser(id)
    if err != nil {
        return err
    }
    return ctx.JSON(http.StatusOK, user)
}

func main() {
    // Initialize logger
    logger := logrus.New()
    logger.SetLevel(logrus.InfoLevel)
    logger.SetFormatter(&logrus.JSONFormatter{})

    // Create module
    userModule := gonest.NewModule("UserModule").
        Controller(&UserController{}).
        Service(NewUserService()).
        Build()

    // Create application
    app := gonest.NewApplication().
        Config(&gonest.Config{Port: "8080"}).
        Logger(logger).
        Build()

    // Register module
    app.ModuleRegistry.Register(userModule)

    // Start application
    if err := app.Start(); err != nil {
        logger.Fatal("Failed to start application:", err)
    }
}
```

## 📚 Documentation

- **[Complete Documentation](DOCUMENTATION.md)** - Comprehensive guide with examples
- **[Features Overview](FEATURES.md)** - Detailed feature descriptions
- **[API Reference](DOCUMENTATION.md#api-reference)** - Complete API documentation

## 🛠️ CLI Tool

GoNest includes a powerful CLI tool for project scaffolding and code generation:

```bash
# Install CLI
go install github.com/ulims/gonest/cmd/gonest@latest

# Create new project
gonest new my-app
cd my-app

# Run application
gonest run

# Build application
gonest build

# Run tests
gonest test
```

## 📁 Project Structure

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
│   │   │   └── dto.go
│   │   └── auth/
│   │       ├── controller.go
│   │       └── service.go
│   ├── config/
│   └── shared/
├── pkg/
├── go.mod
└── go.sum
```

## 🎯 Examples

Check out the [examples](examples/) directory for comprehensive examples:

- **[Basic](examples/basic/)** - Simple CRUD operations
- **[Advanced](examples/advanced/)** - Authentication, validation, and dependency injection
- **[MongoDB](examples/mongodb/)** - Database integration with MongoDB

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [NestJS](https://nestjs.com/) architecture patterns
- Built on top of [Echo](https://echo.labstack.com/) web framework
- Uses [Logrus](https://github.com/sirupsen/logrus) for structured logging

## 📞 Support

- 📖 [Documentation](DOCUMENTATION.md)
- 🐛 [Issues](https://github.com/ulims/gonest/issues)
- 💬 [Discussions](https://github.com/ulims/gonest/discussions)

---

**GoNest** - Building scalable Go applications with elegance and performance.
