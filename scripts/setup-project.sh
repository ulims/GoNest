#!/bin/bash

# GoNest Project Setup Script
# This script creates a new GoNest project with the recommended structure

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.21+ first."
        exit 1
    fi
    
    if ! command_exists git; then
        print_warning "Git is not installed. Git repository will not be initialized."
        GIT_AVAILABLE=false
    else
        GIT_AVAILABLE=true
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    REQUIRED_VERSION="1.21"
    
    if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
        print_error "Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION+"
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}

# Get project details
get_project_details() {
    print_status "Getting project details..."
    
    if [ -z "$1" ]; then
        read -p "Enter project name: " PROJECT_NAME
    else
        PROJECT_NAME=$1
    fi
    
    if [ -z "$PROJECT_NAME" ]; then
        print_error "Project name cannot be empty"
        exit 1
    fi
    
    # Convert to lowercase and replace spaces with hyphens
    PROJECT_NAME=$(echo "$PROJECT_NAME" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')
    
    read -p "Enter Go module name (default: $PROJECT_NAME): " MODULE_NAME
    if [ -z "$MODULE_NAME" ]; then
        MODULE_NAME=$PROJECT_NAME
    fi
    
    read -p "Enter project description: " PROJECT_DESCRIPTION
    if [ -z "$PROJECT_DESCRIPTION" ]; then
        PROJECT_DESCRIPTION="A GoNest application"
    fi
    
    read -p "Enter author name: " AUTHOR_NAME
    if [ -z "$AUTHOR_NAME" ]; then
        AUTHOR_NAME="Your Name"
    fi
    
    read -p "Enter author email: " AUTHOR_EMAIL
    if [ -z "$AUTHOR_EMAIL" ]; then
        AUTHOR_EMAIL="your.email@example.com"
    fi
    
    print_success "Project details collected"
}

# Create project structure
create_project_structure() {
    print_status "Creating project structure..."
    
    # Create main project directory
    mkdir -p "$PROJECT_NAME"
    cd "$PROJECT_NAME"
    
    # Create directory structure
    mkdir -p cmd/server
    mkdir -p internal/modules
    mkdir -p internal/config
    mkdir -p internal/shared/middleware
    mkdir -p internal/shared/utils
    mkdir -p internal/shared/constants
    mkdir -p pkg
    mkdir -p docs
    mkdir -p scripts
    mkdir -p tests
    mkdir -p bin
    mkdir -p build
    mkdir -p deployments
    
    print_success "Project structure created"
}

# Initialize Go module
init_go_module() {
    print_status "Initializing Go module..."
    
    go mod init "$MODULE_NAME"
    
    print_success "Go module initialized"
}

# Initialize Git repository
init_git_repository() {
    if [ "$GIT_AVAILABLE" = true ]; then
        print_status "Initializing Git repository..."
        
        git init
        
        # Create .gitignore
        cat > .gitignore << 'EOF'
# GoNest Project
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
tmp/

# Coverage reports
coverage/

# Database files
*.db
*.sqlite

# Backup files
*.bak
*.backup

# Archive files
*.tar.gz
*.zip
EOF
        
        print_success "Git repository initialized"
    else
        print_warning "Skipping Git initialization (Git not available)"
    fi
}

# Create main application file
create_main_app() {
    print_status "Creating main application file..."
    
    cat > cmd/server/main.go << 'EOF'
package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	gonest "github.com/ulims/GoNest/gonest"
	"MODULE_NAME/internal/config"
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

	// Register application lifecycle hooks
	app.LifecycleManager.RegisterHook(
		gonest.EventApplicationStart,
		gonest.LifecycleHookFunc(func(ctx context.Context) error {
			logger.Info("🚀 Application starting up...")
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
EOF

	# Replace MODULE_NAME placeholder
	if [[ "$OSTYPE" == "darwin"* ]]; then
		# macOS
		sed -i '' "s/MODULE_NAME/$MODULE_NAME/g" cmd/server/main.go
	else
		# Linux/Windows
		sed -i "s/MODULE_NAME/$MODULE_NAME/g" cmd/server/main.go
	fi

	print_success "Main application file created"
}

# Create configuration file
create_config_file() {
    print_status "Creating configuration file..."
    
    cat > internal/config/config.go << 'EOF'
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
EOF

	print_success "Configuration file created"
}

# Create environment file
create_env_file() {
    print_status "Creating environment file..."
    
    cat > .env << 'EOF'
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
EOF

	print_success "Environment file created"
}

# Create README file
create_readme() {
    print_status "Creating README file..."
    
    cat > README.md << EOF
# $PROJECT_NAME

$PROJECT_DESCRIPTION

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Git (optional)

### Installation

1. Clone the repository:
\`\`\`bash
git clone <your-repo-url>
cd $PROJECT_NAME
\`\`\`

2. Install dependencies:
\`\`\`bash
go mod tidy
\`\`\`

3. Set up environment variables:
\`\`\`bash
cp .env.example .env
# Edit .env with your configuration
\`\`\`

4. Run the application:
\`\`\`bash
go run cmd/server/main.go
\`\`\`

## 📁 Project Structure

\`\`\`
$PROJECT_NAME/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── modules/         # Feature modules
│   ├── config/          # Configuration management
│   └── shared/          # Shared utilities
├── pkg/                 # Public packages
├── docs/                # Documentation
├── scripts/             # Build and deployment scripts
├── tests/               # Integration tests
└── deployments/         # Deployment configurations
\`\`\`

## 🎯 Features

- Modular architecture inspired by NestJS
- Dependency injection
- Built-in validation
- JWT authentication
- MongoDB integration
- Redis caching
- WebSocket support
- Rate limiting
- Comprehensive testing utilities

## 📚 Documentation

- [GoNest Framework Documentation](https://github.com/ulims/GoNest)
- [Architecture Guide](ARCHITECTURE.md)
- [API Reference](docs/API.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**$AUTHOR_NAME** - [$AUTHOR_EMAIL](mailto:$AUTHOR_EMAIL)

---

Built with [GoNest](https://github.com/ulims/GoNest) - The Go framework that brings NestJS elegance to Go! 🚀
EOF

	print_success "README file created"
}

# Create Makefile
create_makefile() {
    print_status "Creating Makefile..."
    
    cat > Makefile << 'EOF'
.PHONY: help build run test clean deps lint format docker-build docker-run

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  lint        - Run linter"
	@echo "  format      - Format code"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run Docker container"

# Build the application
build:
	@echo "Building application..."
	@mkdir -p bin
	go build -o bin/server cmd/server/main.go
	@echo "Build complete: bin/server"

# Run the application
run:
	@echo "Running application..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf build/
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies installed"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
format:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t PROJECT_NAME:latest .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 PROJECT_NAME:latest
EOF

	# Replace PROJECT_NAME placeholder
	if [[ "$OSTYPE" == "darwin"* ]]; then
		# macOS
		sed -i '' "s/PROJECT_NAME/$PROJECT_NAME/g" Makefile
	else
		# Linux/Windows
		sed -i "s/PROJECT_NAME/$PROJECT_NAME/g" Makefile
	fi

	print_success "Makefile created"
}

# Create Dockerfile
create_dockerfile() {
    print_status "Creating Dockerfile..."
    
    cat > Dockerfile << 'EOF'
# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk update && apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy environment file
COPY --from=builder /app/.env ./

# Change ownership to non-root user
RUN chown -R appuser:appgroup /root/

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
EOF

	print_success "Dockerfile created"
}

# Create docker-compose file
create_docker_compose() {
    print_status "Creating docker-compose file..."
    
    cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV=development
      - PORT=8080
      - HOST=0.0.0.0
    depends_on:
      - mongodb
      - redis
    networks:
      - app-network

  mongodb:
    image: mongo:6.0
    ports:
      - "27017:27017"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=admin
      - MONGO_INITDB_ROOT_PASSWORD=password
    volumes:
      - mongodb_data:/data/db
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - app-network

  mongo-express:
    image: mongo-express
    ports:
      - "8081:8081"
    environment:
      - ME_CONFIG_MONGODB_ADMINUSERNAME=admin
      - ME_CONFIG_MONGODB_ADMINPASSWORD=password
      - ME_CONFIG_MONGODB_URL=mongodb://admin:password@mongodb:27017/
    depends_on:
      - mongodb
    networks:
      - app-network

volumes:
  mongodb_data:
  redis_data:

networks:
  app-network:
    driver: bridge
EOF

	print_success "Docker Compose file created"
}

# Install dependencies
install_dependencies() {
    print_status "Installing GoNest dependencies..."
    
    # Add GoNest framework
    go get github.com/ulims/GoNest
    
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
    
    # Add environment loading package
    go get github.com/joho/godotenv
    
    # Download and tidy dependencies
    go mod tidy
    
    print_success "Dependencies installed"
}

# Create initial commit
create_initial_commit() {
    if [ "$GIT_AVAILABLE" = true ]; then
        print_status "Creating initial commit..."
        
        git add .
        git commit -m "Initial commit: GoNest project setup"
        
        print_success "Initial commit created"
    fi
}

# Print next steps
print_next_steps() {
    print_success "Project setup complete!"
    echo
    echo -e "${GREEN}🎉 Your GoNest project '$PROJECT_NAME' has been created successfully!${NC}"
    echo
    echo -e "${BLUE}Next steps:${NC}"
    echo "1. cd $PROJECT_NAME"
    echo "2. Review and customize the generated files"
    echo "3. Update the GoNest import path in main.go"
    echo "4. Run 'go mod tidy' to ensure dependencies are correct"
    echo "5. Start building your modules!"
    echo
    echo -e "${BLUE}Useful commands:${NC}"
    echo "  make build    - Build the application"
    echo "  make run      - Run the application"
    echo "  make test     - Run tests"
    echo "  make clean    - Clean build artifacts"
    echo
    echo -e "${BLUE}Documentation:${NC}"
    echo "  - README.md - Project overview and setup instructions"
    echo "  - ARCHITECTURE.md - GoNest architecture guide"
    echo "  - docs/ - Additional documentation"
    echo
    echo -e "${YELLOW}Happy coding with GoNest! 🚀${NC}"
}

# Main function
main() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}   GoNest Project Setup Script${NC}"
    echo -e "${BLUE}================================${NC}"
    echo
    
    # Check prerequisites
    check_prerequisites
    
    # Get project details
    get_project_details "$1"
    
    # Create project structure
    create_project_structure
    
    # Initialize Go module
    init_go_module
    
    # Initialize Git repository
    init_git_repository
    
    # Create main application file
    create_main_app
    
    # Create configuration file
    create_config_file
    
    # Create environment file
    create_env_file
    
    # Create README file
    create_readme
    
    # Create Makefile
    create_makefile
    
    # Create Dockerfile
    create_dockerfile
    
    # Create docker-compose file
    create_docker_compose
    
    # Install dependencies
    install_dependencies
    
    # Create initial commit
    create_initial_commit
    
    # Print next steps
    print_next_steps
}

# Run main function with command line argument
main "$@"

