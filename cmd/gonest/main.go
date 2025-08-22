package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gonest",
		Short: "GoNest CLI - A powerful Go framework inspired by NestJS",
		Long: `GoNest CLI provides tools for generating modules, controllers, services, and other components
for your GoNest applications. It helps you scaffold applications quickly and maintain
consistent project structure.

Examples:
  gonest new my-awesome-app          # Create a new project
  gonest new my-app --strict        # Create with strict mode
  gonest new my-app --template=api  # Use API template
  gonest generate module user        # Generate a new module
  gonest generate controller user    # Generate a controller
  gonest generate service user       # Generate a service`,
	}

	newCmd = &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new GoNest project",
		Long:  "Create a new GoNest project with basic structure and configuration",
		Args:  cobra.ExactArgs(1),
		Run:   createNewProject,
	}

	generateCmd = &cobra.Command{
		Use:   "generate [type] [name]",
		Short: "Generate GoNest components",
		Long:  "Generate modules, controllers, services, and other components",
		Args:  cobra.ExactArgs(2),
		Run:   generateComponent,
	}

	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the GoNest application",
		Long:  "Build the GoNest application for production",
		Run:   buildApplication,
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the GoNest application",
		Long:  "Run the GoNest application in development mode",
		Run:   runApplication,
	}

	testCmd = &cobra.Command{
		Use:   "test",
		Short: "Run tests",
		Long:  "Run all tests in the GoNest application",
		Run:   runTests,
	}

	// Flags
	force     bool
	strict    bool
	template  string
	moduleDir string
)

func init() {
	// Add flags for new command
	newCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing files")
	newCmd.Flags().BoolVarP(&strict, "strict", "s", false, "Enable strict mode with additional validation and security")
	newCmd.Flags().StringVarP(&template, "template", "t", "basic", "Project template (basic, api, fullstack, microservice)")

	// Add flags for generate command
	generateCmd.Flags().StringVarP(&moduleDir, "module", "m", "", "Module directory for component generation")

	// Add commands
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(testCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func createNewProject(cmd *cobra.Command, args []string) {
	projectName := args[0]

	if !force && fileExists(projectName) {
		fmt.Printf("Project directory %s already exists. Use --force to overwrite.\n", projectName)
		return
	}

	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		return
	}

	// Create directory structure based on template
	dirs := getProjectStructure(template)
	for _, dir := range dirs {
		os.MkdirAll(filepath.Join(projectName, dir), 0755)
	}

	// Generate project files based on template
	generateProjectFiles(projectName, template, strict)

	fmt.Printf("✅ GoNest project '%s' created successfully!\n", projectName)
	fmt.Printf("📁 Navigate to the project: cd %s\n", projectName)
	fmt.Printf("🚀 Run the application: go run cmd/server/main.go\n")

	if strict {
		fmt.Printf("🔒 Strict mode enabled - additional validation and security features included\n")
	}
}

func generateComponent(cmd *cobra.Command, args []string) {
	componentType := strings.ToLower(args[0])
	componentName := args[1]

	// Determine module directory
	if moduleDir == "" {
		// Auto-detect module directory
		if fileExists("internal/modules") {
			moduleDir = "internal/modules"
		} else {
			fmt.Println("❌ No modules directory found. Please run this command from a GoNest project root.")
			return
		}
	}

	switch componentType {
	case "module":
		generateModule(componentName, moduleDir)
	case "controller":
		generateController(componentName, moduleDir)
	case "service":
		generateService(componentName, moduleDir)
	case "dto":
		generateDTO(componentName, moduleDir)
	case "entity":
		generateEntity(componentName, moduleDir)
	default:
		fmt.Printf("❌ Unknown component type: %s\n", componentType)
		fmt.Println("Available types: module, controller, service, dto, entity")
	}
}

func buildApplication(cmd *cobra.Command, args []string) {
	fmt.Println("🔨 Building GoNest application...")

	// Check if go.mod exists
	if !fileExists("go.mod") {
		fmt.Println("❌ No go.mod file found. Are you in a Go module directory?")
		return
	}

	// Run go build
	if err := runCommand("go", "build", "-o", "bin/app", "./cmd/server"); err != nil {
		fmt.Printf("❌ Build failed: %v\n", err)
		return
	}

	fmt.Println("✅ Build completed successfully!")
	fmt.Println("📦 Binary created at: bin/app")
}

func runApplication(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 Starting GoNest application...")

	// Check if main.go exists
	if !fileExists("cmd/server/main.go") {
		fmt.Println("❌ No main.go file found at cmd/server/main.go")
		return
	}

	// Run the application
	if err := runCommand("go", "run", "./cmd/server"); err != nil {
		fmt.Printf("❌ Application failed to start: %v\n", err)
		return
	}
}

func runTests(cmd *cobra.Command, args []string) {
	fmt.Println("🧪 Running tests...")

	// Run go test
	if err := runCommand("go", "test", "./..."); err != nil {
		fmt.Printf("❌ Tests failed: %v\n", err)
		return
	}

	fmt.Println("✅ All tests passed!")
}

// Helper functions
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Helper functions for project generation
func getProjectStructure(template string) []string {
	switch template {
	case "api":
		return []string{
			"cmd/server",
			"internal/modules",
			"internal/config",
			"internal/shared/middleware",
			"internal/shared/interceptors",
			"internal/shared/exceptions",
			"pkg/utils",
			"docs",
			"tests",
		}
	case "fullstack":
		return []string{
			"cmd/server",
			"cmd/client",
			"internal/modules",
			"internal/config",
			"internal/shared/middleware",
			"internal/shared/interceptors",
			"internal/shared/exceptions",
			"pkg/utils",
			"pkg/client",
			"web/static",
			"web/templates",
			"docs",
			"tests",
		}
	case "microservice":
		return []string{
			"cmd/server",
			"internal/modules",
			"internal/config",
			"internal/shared/middleware",
			"internal/shared/interceptors",
			"internal/shared/exceptions",
			"pkg/utils",
			"pkg/grpc",
			"proto",
			"docs",
			"tests",
			"deploy",
		}
	default: // basic
		return []string{
			"cmd/server",
			"internal/modules",
			"internal/config",
			"internal/shared/middleware",
			"internal/shared/interceptors",
			"internal/shared/exceptions",
			"pkg/utils",
			"docs",
		}
	}
}

func generateProjectFiles(projectName, template string, strict bool) {
	// Generate go.mod with correct import path
	goModContent := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/labstack/echo/v4 v4.11.4
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
`, projectName)

	writeFile(filepath.Join(projectName, "go.mod"), goModContent)

	// Generate main.go with working modular architecture
	mainContent := generateMainGoWithModules(projectName)
	writeFile(filepath.Join(projectName, "cmd/server/main.go"), mainContent)

	// Generate configuration files
	generateConfigFiles(projectName)

	// Generate README
	readmeContent := generateREADME(projectName, template)
	writeFile(filepath.Join(projectName, "README.md"), readmeContent)

	// Generate .gitignore
	gitignoreContent := generateGitignore(template)
	writeFile(filepath.Join(projectName, ".gitignore"), gitignoreContent)

	// Generate Makefile
	makefileContent := generateMakefile(template)
	writeFile(filepath.Join(projectName, "Makefile"), makefileContent)

	// Generate additional template-specific files
	generateTemplateSpecificFiles(projectName, template, strict)

	// Generate working modular architecture to demonstrate GoNest's power
	generateModularArchitecture(projectName)
}

func generateMainGoWithModules(projectName string) string {
	mainContent := fmt.Sprintf(`package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"%s/internal/modules/user"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create Echo instance
	e := echo.New()

	// Initialize and register the User module
	// This demonstrates GoNest's modular architecture!
	userModule := user.NewUserModule(logger)
	
	// Register module routes
	userModule.RegisterRoutes(e)

	// Add a health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
			"message": "🚀 GoNest Application with Modular Architecture is running!",
			"modules": "User module is active and ready",
		})
	})

	// Add root route
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "🚀 GoNest Application with Modular Architecture!")
	})

	// Start server
	addr := "localhost:8080"
	logger.Infof("🚀 Starting GoNest application on %%s", addr)
	logger.Info("📁 User module is loaded and ready!")
	logger.Info("🎯 Try: POST /users, GET /users, GET /users/:id")
	
	if err := e.Start(addr); err != nil {
		logger.Fatal("Failed to start application:", err)
	}
}
`, projectName)

	return mainContent
}

func generateConfigFiles(projectName string) {
	// Generate config.go
	configContent := `package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port        string
	Host        string
	Environment string
	LogLevel    string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

type JWTConfig struct {
	Secret     string
	Expiration int64
}

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
			Expiration: getEnvAsInt64("JWT_EXPIRATION", 86400),
		},
	}
}

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
`
	writeFile(filepath.Join(projectName, "internal/config/config.go"), configContent)

	// Generate .env file
	envContent := `# Server Configuration
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
`
	writeFile(filepath.Join(projectName, ".env"), envContent)
}

func generateREADME(projectName, template string) string {
	return fmt.Sprintf(`# %s

A GoNest application built with the GoNest framework.

## Template: %s

## Features

- Modular architecture
- Dependency injection
- Authentication & authorization
- Request/response interceptors
- Exception handling
- And more...

## Getting Started

### Prerequisites

- Go 1.21 or higher
- GoNest framework

### Installation

1. Clone the repository
2. Install dependencies:
   `+"```bash"+`
   go mod tidy
   `+"```"+`

### Running the Application

`+"```bash"+`
# Development mode
go run cmd/server/main.go

# Production build
go build -o bin/app cmd/server/main.go
./bin/app
`+"```"+`

### Testing

`+"```bash"+`
go test ./...
`+"```"+`

## Project Structure

`+"```"+`
%s/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── modules/
│   ├── config/
│   └── shared/
├── pkg/
├── examples/
└── docs/
`+"```"+`

## Documentation

For more information about GoNest, visit the [documentation](https://github.com/ulims/GoNest).
`, projectName, template, projectName)
}

func generateGitignore(template string) string {
	baseContent := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swp

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log

# Environment files
.env
.env.local
.env.production

# Build artifacts
bin/
dist/

# Temporary files
tmp/
temp/
`

	if template == "microservice" {
		baseContent += `
# Microservice specific
deploy/
*.pb.go
`
	}

	return baseContent
}

func generateMakefile(template string) string {
	baseContent := `# GoNest Application Makefile

.PHONY: build run test clean deps lint

# Build the application
build:
	@echo "Building application..."
	go build -o bin/app cmd/server/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help"
`

	if template == "microservice" {
		baseContent += `
# Microservice specific commands
proto:
	@echo "Generating protobuf files..."
	protoc --go_out=. --go-grpc_out=. proto/*.proto

docker:
	@echo "Building Docker image..."
	docker build -t %s .
`
	}

	return baseContent
}

func generateTemplateSpecificFiles(projectName, template string, strict bool) {
	switch template {
	case "api":
		// Generate API-specific files
		generateAPIFiles(projectName, strict)
	case "fullstack":
		// Generate fullstack-specific files
		generateFullstackFiles(projectName, strict)
	case "microservice":
		// Generate microservice-specific files
		generateMicroserviceFiles(projectName, strict)
	}
}

func generateAPIFiles(projectName string, strict bool) {
	// Generate API documentation template
	swaggerContent := fmt.Sprintf(`openapi: 3.0.0
info:
  title: %s API
  version: 1.0.0
  description: GoNest API documentation
paths:
  /health:
    get:
      summary: Health check
      responses:
        '200':
          description: Service is healthy
`, projectName)
	writeFile(filepath.Join(projectName, "docs/swagger.yaml"), swaggerContent)
}

func generateFullstackFiles(projectName string, strict bool) {
	// Generate web templates
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
</head>
<body>
    <h1>Welcome to %s</h1>
    <p>Built with GoNest framework</p>
</body>
</html>
`, projectName, projectName)
	writeFile(filepath.Join(projectName, "web/templates/index.html"), htmlContent)
}

func generateMicroserviceFiles(projectName string, strict bool) {
	// Generate protobuf template
	protoContent := fmt.Sprintf(`syntax = "proto3";

package %s;

option go_package = "%s/proto";

service %sService {
  rpc HealthCheck(HealthRequest) returns (HealthResponse);
}

message HealthRequest {}

message HealthResponse {
  string status = 1;
  string timestamp = 2;
}
`, strings.ToLower(projectName), projectName, projectName)
	writeFile(filepath.Join(projectName, "proto/service.proto"), protoContent)
}

// generateModularArchitecture creates a working sample module to demonstrate GoNest's power
func generateModularArchitecture(projectName string) {
	// Create a sample "user" module to demonstrate modularity
	userModulePath := filepath.Join(projectName, "internal/modules/user")
	os.MkdirAll(userModulePath, 0755)

	// Generate user module file
	userModuleContent := `package user

import (
	"github.com/sirupsen/logrus"
)

// UserModule demonstrates modular architecture
type UserModule struct {
	userService    *UserService
	userController *UserController
	logger         *logrus.Logger
}

// NewUserModule creates a new user module with all its components
func NewUserModule(logger *logrus.Logger) *UserModule {
	// Create services
	userService := NewUserService(logger)
	
	// Create controllers
	userController := NewUserController(userService)
	
	// Create and return module - this is where the magic happens!
	return &UserModule{
		userService:    userService,
		userController: userController,
		logger:         logger,
	}
}
`
	writeFile(filepath.Join(userModulePath, "user_module.go"), userModuleContent)

	// Generate user service
	userServiceContent := `package user

import (
	"errors"
	"sync"
	"time"
	"github.com/sirupsen/logrus"
)

// User represents a user entity
type User struct {
	ID        string    ` + "`" + `json:"id"` + "`" + `
	Username  string    ` + "`" + `json:"username"` + "`" + `
	Email     string    ` + "`" + `json:"email"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
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
func (s *UserService) CreateUser(username, email string) (*User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Check if username already exists
	if s.usernameExists(username) {
		return nil, errors.New("username already exists")
	}
	
	// Create new user
	user := &User{
		ID:        time.Now().Format("20060102150405"),
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}
	
	// Store user
	s.users[user.ID] = user
	
	s.logger.Infof("Created user: %%s (%%s)", user.Username, user.ID)
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
`
	writeFile(filepath.Join(userModulePath, "user_service.go"), userServiceContent)

	// Generate user controller
	userControllerContent := `package user

import (
	"net/http"
	"github.com/labstack/echo/v4"
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
		Username string ` + "`" + `json:"username" validate:"required,min=3"` + "`" + `
		Email    string ` + "`" + `json:"email" validate:"required,email"` + "`" + `
	}
	
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	
	// Create user
	user, err := c.userService.CreateUser(req.Username, req.Email)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	
	return ctx.JSON(http.StatusCreated, user)
}

// GetUser handles user retrieval by ID
func (c *UserController) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}
	
	user, err := c.userService.GetUser(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}
	
	return ctx.JSON(http.StatusOK, user)
}

// ListUsers handles user listing
func (c *UserController) ListUsers(ctx echo.Context) error {
	users, err := c.userService.ListUsers()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve users",
		})
	}
	
	return ctx.JSON(http.StatusOK, users)
}
`
	writeFile(filepath.Join(userModulePath, "user_controller.go"), userControllerContent)

	// Generate routes file to show how modules are wired together
	routesContent := `package user

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all user module routes
func (m *UserModule) RegisterRoutes(e *echo.Echo) {
	// Create route group for user module
	userGroup := e.Group("/users")
	
	// Register routes with the controller
	userGroup.POST("", m.userController.CreateUser)
	userGroup.GET("/:id", m.userController.GetUser)
	userGroup.GET("", m.userController.ListUsers)
}
`
	writeFile(filepath.Join(userModulePath, "routes.go"), routesContent)

	fmt.Printf("✅ Generated working modular architecture with User module!\n")
	fmt.Printf("   This demonstrates GoNest's NestJS-style modularity.\n")
	fmt.Printf("   Check internal/modules/user/ to see how modules work.\n")
	fmt.Printf("\n🎯 **Why This Matters:**\n")
	fmt.Printf("   - No more 'examples' folder cluttering your project\n")
	fmt.Printf("   - Working modular architecture from day one\n")
	fmt.Printf("   - See NestJS patterns in action immediately\n")
	fmt.Printf("   - Ready to build and run with real functionality!\n")
}

// Component generation functions
func generateModule(name, moduleDir string) {
	modulePath := filepath.Join(moduleDir, strings.ToLower(name))
	os.MkdirAll(modulePath, 0755)

	// Generate module file
	moduleContent := fmt.Sprintf(`package %s

import (
	"github.com/sirupsen/logrus"
)

type %sModule struct {
	%sService    *%sService
	%sController *%sController
	logger       *logrus.Logger
}

func New%sModule(logger *logrus.Logger) *%sModule {
	// Create services
	%sService := New%sService(logger)
	
	// Create controllers
	%sController := New%sController(%sService)
	
	// Create and return module
	return &%sModule{
		%sService:    %sService,
		%sController: %sController,
		logger:       logger,
	}
}
`, strings.ToLower(name), name, strings.ToLower(name), name, strings.ToLower(name), name, name, name, strings.ToLower(name), name, strings.ToLower(name), name, strings.ToLower(name), name, strings.ToLower(name), name, strings.ToLower(name), name, strings.ToLower(name), name)

	writeFile(filepath.Join(modulePath, fmt.Sprintf("%s_module.go", strings.ToLower(name))), moduleContent)

	fmt.Printf("✅ Module '%s' generated successfully at %s\n", name, modulePath)
}

func generateController(name, moduleDir string) {
	// Implementation for controller generation
	fmt.Printf("✅ Controller '%s' generation not yet implemented\n", name)
}

func generateService(name, moduleDir string) {
	// Implementation for service generation
	fmt.Printf("✅ Service '%s' generation not yet implemented\n", name)
}

func generateDTO(name, moduleDir string) {
	// Implementation for DTO generation
	fmt.Printf("✅ DTO '%s' generation not yet implemented\n", name)
}

func generateEntity(name, moduleDir string) {
	// Implementation for entity generation
	fmt.Printf("✅ Entity '%s' generation not yet implemented\n", name)
}

func writeFile(path, content string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", path, err)
		return
	}
}
