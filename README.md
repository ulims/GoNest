# GoNest Framework

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ulims/GoNest)](https://goreportcard.com/report/github.com/ulims/GoNest)

A powerful, enterprise-grade Go web framework inspired by NestJS, designed for building scalable, maintainable applications with modern architectural patterns.

## 🚀 Features

- **🏗️ Modular Architecture**: NestJS-style module system with dependency injection
- **🔄 Lifecycle Management**: Comprehensive application and module lifecycle hooks
- **🛡️ Built-in Security**: Guards, interceptors, and authentication systems
- **📊 Database Integration**: MongoDB support with Mongoose-like ODM
- **🌐 WebSocket Support**: Real-time communication capabilities
- **⚡ High Performance**: Built on Echo framework for optimal performance
- **🧪 Testing Utilities**: Built-in testing framework and utilities
- **📝 Validation**: Request/response validation with struct tags
- **🎯 CLI Tools**: Code generation and project management tools

## 📁 Project Structure

When you create a new GoNest application, you'll get a well-organized project structure that follows Go and NestJS best practices:

```
my-gonest-app/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── modules/              # Feature modules (business domains)
│   │   ├── user/            # User module example
│   │   │   ├── user_module.go     # Module definition and DI setup
│   │   │   ├── user_service.go    # Business logic layer
│   │   │   ├── user_controller.go # HTTP request handlers
│   │   │   └── user_dto.go        # Data transfer objects
│   │   └── auth/            # Authentication module
│   ├── config/              # Configuration management
│   │   └── config.go
│   ├── middleware/          # Custom middleware
│   └── shared/              # Shared utilities and types
├── pkg/                     # Public packages (if needed)
├── scripts/                 # Build and deployment scripts
├── docs/                    # Project documentation
├── tests/                   # Integration and e2e tests
├── .env                     # Environment variables
├── .gitignore              # Git ignore rules
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── Dockerfile              # Container configuration
├── docker-compose.yml      # Multi-service setup
└── README.md              # Project documentation
```

### Module Structure

Each feature module follows a consistent, flat structure inspired by NestJS:

```
internal/modules/user/
├── user_module.go          # Module registration and dependency injection
├── user_controller.go      # HTTP request handling
├── user_service.go         # Business logic implementation
├── user_dto.go            # Request/response data structures
├── user_entity.go         # Domain entities/models
└── user_repository.go     # Data access layer (if needed)
```

## 🚀 Installation

To get started, you can either use our **automated setup scripts** (recommended), or manually set up your project. Both approaches will produce the same outcome.

### 🚀 Automated Setup (Recommended)

Our setup scripts provide the fastest and most reliable way to create a new GoNest project:

#### Linux/macOS
```bash
# Clone the GoNest repository
$ git clone https://github.com/ulims/GoNest.git
$ cd GoNest

# Make the script executable and run it
$ chmod +x scripts/setup-project.sh
$ ./scripts/setup-project.sh my-awesome-app
```

#### Windows
```cmd
# Clone the GoNest repository
$ git clone https://github.com/ulims/GoNest.git
$ cd GoNest

# Run the batch script
$ scripts\setup-project.bat
```

The setup scripts automatically:
- ✅ Create the recommended project structure
- ✅ Initialize Go module and Git repository
- ✅ Install all GoNest dependencies
- ✅ Generate configuration files
- ✅ Set up Docker and build automation
- ✅ Create comprehensive documentation

> **HINT**  
> The setup scripts are the most reliable way to get started. They handle all dependencies and create a production-ready project structure.

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





### Basic Application Structure

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
    
    // Create application
    app := gonest.NewApplication().
        Config(&gonest.Config{
            Port:        "8080",
            Host:        "localhost",
            Environment: "development",
        }).
        Logger(logger).
        Build()
    
    // Start application
    if err := app.Start(); err != nil {
        logger.Fatal("Failed to start application:", err)
    }
}
```

### Run Your Application

```bash
go run main.go
```

Your application will be available at `http://localhost:8080`

## 📚 Documentation

- **[📖 Full Documentation](docs/DOCUMENTATION.md)** - Comprehensive framework guide
- **[🏗️ Architecture Guide](ARCHITECTURE.md)** - Detailed architectural patterns
- **[🚀 Quick Start Guide](docs/QUICKSTART.md)** - Step-by-step project setup
- **[📋 Features Overview](docs/DOCUMENTATION.md#features)** - All available features
- **[🧪 Examples](examples/)** - Working examples and tutorials
- **[🔧 Setup Scripts](scripts/README.md)** - Automated project initialization

## 🎯 Architecture Principles

GoNest follows proven architectural patterns inspired by NestJS:

- **🏗️ Modular Design**: Organize code into feature modules with clear boundaries
- **💉 Dependency Injection**: Automatic dependency resolution and injection
- **📱 Flat Module Structure**: Keep module files organized without deep nesting
- **🔄 Separation of Concerns**: Controllers handle HTTP, Services handle business logic
- **🧪 Testable by Design**: Easy to mock and test individual components
- **📐 Consistent Patterns**: Every module follows the same organizational structure

### Quick Start Example

```go
// Create a complete user module in minutes
userModule := gonest.NewModule("UserModule").
    Controller(userController).
    Service(userService).
    Build()

app.ModuleRegistry.Register(userModule)
```

## 🛠️ CLI Tool (Coming Soon)

GoNest will include a powerful CLI tool for project management. Currently in development:

```bash
# Future CLI commands (not yet available)
gonest new my-project      # Create new project
gonest build               # Build project
gonest run                 # Run project
gonest test                # Run tests
```

> **NOTE**  
> The CLI tool is currently in development. For now, use our [automated setup scripts](#-automated-setup-recommended) for the best experience.

## 🔧 Key Components

### Modules
```go
type UserModule struct {
    *gonest.Module
}

func NewUserModule(logger *logrus.Logger) *UserModule {
    userService := NewUserService(logger)
    userController := NewUserController(userService)
    
    module := gonest.NewModule("UserModule").
        Controller(userController).
        Service(userService).
        Build()
    
    return &UserModule{Module: module}
}
```

### Services
```go
type UserService struct {
    users  map[string]*User
    logger *logrus.Logger
}

func (s *UserService) CreateUser(username, email, password string) (*User, error) {
    // Business logic implementation
}
```

### Controllers
```go
type UserController struct {
    userService *UserService
}

func (c *UserController) CreateUser(ctx echo.Context) error {
    // HTTP request handling
}
```

## 🧪 Testing

GoNest provides built-in testing utilities:

```go
func TestUserService(t *testing.T) {
    testApp := gonest.NewTestApp().
        Module(userModule).
        Build()
    
    // Test with real HTTP requests
    response := testApp.Request("POST", "/users").
        WithJSON(map[string]interface{}{
            "username": "testuser",
            "email":    "test@example.com",
        }).
        ExpectStatus(201).
        Get()
    
    assert.NotNil(t, response.JSON())
}
```

## 📈 Performance

- **High Throughput**: Built on Echo framework for optimal performance
- **Low Memory Usage**: Efficient memory management and garbage collection
- **Fast Startup**: Minimal initialization overhead
- **Scalable**: Support for horizontal and vertical scaling

## 🔒 Security

- **Input Validation**: Comprehensive request validation
- **Authentication**: JWT-based authentication system
- **Authorization**: Role-based access control
- **Rate Limiting**: Built-in rate limiting strategies
- **CORS Support**: Configurable cross-origin resource sharing

## 🌟 Why GoNest?

- **🚀 NestJS Familiarity**: If you know NestJS, you'll feel at home
- **⚡ Go Performance**: Leverage Go's speed and efficiency
- **🏗️ Enterprise Ready**: Built for production applications
- **🔧 Developer Experience**: Excellent tooling and documentation
- **📚 Rich Ecosystem**: Comprehensive feature set out of the box
- **🧪 Testing First**: Built-in testing utilities and patterns
- **⚡ Quick Setup**: Automated project initialization scripts

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Areas for Contribution

- 🐛 Bug fixes and improvements
- ✨ New features and enhancements
- 📚 Documentation improvements
- 🧪 Test coverage expansion
- 🔧 Performance optimizations

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **NestJS Team**: For the excellent architectural inspiration
- **Echo Framework**: For the high-performance HTTP foundation
- **Go Community**: For the amazing ecosystem and tools

## 📞 Support

- **📖 Documentation**: [Full Documentation](docs/DOCUMENTATION.md)
- **🐛 Issues**: [GitHub Issues](https://github.com/ulims/GoNest/issues)
- **💬 Discussions**: [GitHub Discussions](https://github.com/ulims/GoNest/discussions)
- **📧 Email**: your-email@example.com

---

**Build amazing applications with GoNest - The Go framework that brings NestJS elegance to Go! 🚀**
