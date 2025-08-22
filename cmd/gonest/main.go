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
			"examples",
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
	gonest "github.com/ulims/GoNest"
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

	// Generate main.go with template-specific content
	mainContent := generateMainGo(template, strict)
	writeFile(filepath.Join(projectName, "cmd/server/main.go"), mainContent)

	// Generate configuration files
	generateConfigFiles(projectName, template, strict)

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
}

func generateMainGo(template string, strict bool) string {
	baseContent := `package main

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

	// Create application
	app := gonest.NewApplication().
		Config(&gonest.Config{
			Port:        "8080",
			Host:        "localhost",
			Environment: "development",
		}).
		Logger(logger).
		Build()

	// Register your modules here
	// app.ModuleRegistry.Register(yourModule.GetModule())

	// Start the application
	if err := app.Start(); err != nil {
		logger.Fatal("Failed to start application:", err)
	}
}
`

	if strict {
		baseContent = strings.Replace(baseContent, "Environment: \"development\"", "Environment: \"development\",\n\t\t\tLogLevel:    \"info\",", 1)
	}

	return baseContent
}

func generateConfigFiles(projectName, template string, strict bool) {
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

// Component generation functions
func generateModule(name, moduleDir string) {
	modulePath := filepath.Join(moduleDir, strings.ToLower(name))
	os.MkdirAll(modulePath, 0755)

	// Generate module file
	moduleContent := fmt.Sprintf(`package %s

import (
	"github.com/sirupsen/logrus"
	gonest "github.com/ulims/GoNest"
)

type %sModule struct {
	*gonest.Module
}

func New%sModule(logger *logrus.Logger) *%sModule {
	// Create services
	%sService := New%sService(logger)
	
	// Create controllers
	%sController := New%sController(%sService)
	
	// Create and return module
	module := gonest.NewModule("%sModule").
		Controller(%sController).
		Service(%sService).
		Build()
	
	return &%sModule{
		Module: module,
	}
}
`, strings.ToLower(name), name, name, name, strings.ToLower(name), name, strings.ToLower(name), name, strings.ToLower(name), name, name, name, name, name)

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
