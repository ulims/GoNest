# GoNest CLI Tool Guide

This guide is specifically designed for **external developers** who want to use the GoNest CLI tool to create and manage GoNest projects.

## 🚀 **Quick Start for External Developers**

### **Prerequisites**
- Go 1.21+ installed
- Git installed
- Basic knowledge of Go modules

### **Step-by-Step Installation**

#### **1. One-Line Installation (Recommended)**
```bash
# Linux/macOS
$ curl -sSL https://raw.githubusercontent.com/ulims/GoNest/master/install-gonest.sh | bash

# Windows
$ curl -o install-gonest.bat https://raw.githubusercontent.com/ulims/GoNest/master/install-gonest.bat
$ install-gonest.bat
```

#### **2. Manual Installation**
```bash
# Install GoNest CLI globally
$ git clone https://github.com/ulims/GoNest.git
$ cd GoNest
$ go install ./cmd/gonest

# Verify installation
$ gonest --help
```

> **Note**: The one-line installer handles everything automatically. For manual installation, we need to clone the repository first because Go's `go install` command works with local modules. This is the standard approach for Go CLI tools and provides the same developer experience as `npm install -g @nestjs/cli`.

#### **2. Alternative: Build from Source**
```bash
# Clone the GoNest repository
$ git clone https://github.com/ulims/GoNest.git
$ cd GoNest

# Build the CLI tool
$ go build -o gonest.exe cmd/gonest/main.go  # Windows
$ go build -o gonest cmd/gonest/main.go       # Linux/macOS

# Verify installation
$ ./gonest.exe --help  # Windows
$ ./gonest --help       # Linux/macOS
```

## 🎯 **Creating Your First Project**

### **Basic Project**
```bash
# Create a basic GoNest project
$ gonest new my-awesome-app

# Navigate to your new project
$ cd my-awesome-app

# Install dependencies
$ go mod tidy
```

### **API Project with Strict Mode**
```bash
# Create an API-focused project with enhanced security
$ gonest new my-api --template=api --strict

# Navigate to your new project
$ cd my-api

# Install dependencies
$ go mod tidy
```

### **Full-Stack Project**
```bash
# Create a full-stack project with web templates
$ gonest new my-webapp --template=fullstack

# Navigate to your new project
$ cd my-webapp

# Install dependencies
$ go mod tidy
```

### **Microservice Project**
```bash
# Create a microservice project with gRPC support
$ gonest new my-service --template=microservice

# Navigate to your new project
$ cd my-service

# Install dependencies
$ go mod tidy
```

## 🔧 **Generating Components**

Once you have a project created, you can generate additional components:

### **Generate a Module**
```bash
# Navigate to your GoNest project root
$ cd my-awesome-app

# Generate a user module
$ gonest generate module user
```

### **Generate a Controller**
```bash
# Generate a user controller
$ gonest generate controller user
```

### **Generate a Service**
```bash
# Generate a user service
$ gonest generate service user
```

### **Generate DTOs and Entities**
```bash
# Generate user DTOs
$ gonest generate dto user

# Generate user entities
$ gonest generate entity user
```

## 🏗️ **Project Management Commands**

### **Build Your Application**
```bash
# Build the application
$ gonest build
```

### **Run Your Application**
```bash
# Run the application in development mode
$ gonest run
```

### **Test Your Application**
```bash
# Run all tests
$ gonest test
```

## 📋 **Available Templates**

| Template | Description | Best For |
|----------|-------------|----------|
| `basic` | Standard GoNest structure | General applications, learning |
| `api` | API-focused with Swagger docs | REST APIs, microservices |
| `fullstack` | Web app with HTML templates | Full-stack web applications |
| `microservice` | gRPC + protobuf support | Microservice architecture |

## 🎯 **Command Reference**

### **New Project Commands**
```bash
# Basic syntax
gonest new <project-name> [flags]

# Examples
gonest new my-app                    # Basic project
gonest new my-api --template=api     # API project
gonest new my-web --template=fullstack # Full-stack project
gonest new my-service --template=microservice # Microservice
gonest new my-app --strict          # With strict mode
gonest new my-app --force           # Overwrite existing
```

### **Generate Commands**
```bash
# Basic syntax
gonest generate <type> <name> [flags]

# Examples
gonest generate module user          # Generate user module
gonest generate controller user      # Generate user controller
gonest generate service user         # Generate user service
gonest generate dto user             # Generate user DTOs
gonest generate entity user          # Generate user entities
```

### **Project Management Commands**
```bash
gonest build                         # Build the application
gonest run                           # Run the application
gonest test                          # Run tests
gonest --help                        # Show help
```

## 🔒 **Strict Mode Features**

When you use `--strict` flag, your project gets enhanced security and validation:

- **Enhanced Input Validation**: Strict request validation
- **Security Headers**: CORS, XSS protection, etc.
- **Rate Limiting**: Built-in rate limiting
- **Request Logging**: Comprehensive request/response logging
- **Error Handling**: Enhanced error handling and logging

## 📁 **Generated Project Structure**

After running `./gonest.exe new my-app`, you'll get:

```
my-app/
├── cmd/
│   └── server/
│       └── main.go          # Main application entry point
├── internal/
│   ├── modules/             # Your business modules
│   ├── config/              # Configuration management
│   └── shared/              # Shared utilities
├── pkg/                     # Public packages
├── docs/                    # Documentation
├── go.mod                   # Go module file
├── .env                     # Environment variables
├── .gitignore              # Git ignore file
├── Makefile                # Build automation
└── README.md               # Project documentation
```

## 🚨 **Troubleshooting**

### **Common Issues**

#### **1. "command not found" error**
```bash
# If using global installation
$ gonest --help
# Should show help if installed correctly

# If using local build, make sure you're in the GoNest directory
$ pwd
# Should show: /path/to/GoNest

# Make sure gonest.exe exists
$ ls -la gonest.exe
```

#### **2. "go: module not found" error**
```bash
# Navigate to your project directory
$ cd my-app

# Install dependencies
$ go mod tidy
```

#### **3. Permission denied on Linux/macOS**
```bash
# Make the CLI tool executable
$ chmod +x gonest
```

### **Getting Help**
```bash
# Show general help
$ gonest --help

# Show command-specific help
$ gonest new --help
$ gonest generate --help
```

## 🔄 **Updating the CLI Tool**

To get the latest version of the CLI tool:

```bash
# Navigate to GoNest directory
$ cd GoNest

# Pull latest changes
$ git pull origin main

# Rebuild the CLI tool
$ go build -o gonest.exe cmd/gonest/main.go
```

## 📚 **Next Steps**

After creating your project:

1. **Explore the generated code** in `cmd/server/main.go`
2. **Customize configuration** in `internal/config/config.go`
3. **Add your modules** using `./gonest.exe generate module <name>`
4. **Run your application** with `./gonest.exe run`
5. **Check the documentation** in the `docs/` folder

## 🆘 **Need Help?**

- **Documentation**: Check the [main README](../README.md)
- **Architecture**: Review the [architecture guide](../ARCHITECTURE.md)
- **Examples**: Look at the [examples](../examples/) directory
- **Issues**: Open an [issue on GitHub](https://github.com/ulims/GoNest/issues)

---

**Happy coding with GoNest! 🎉**
