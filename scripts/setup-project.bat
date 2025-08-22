@echo off
setlocal enabledelayedexpansion

REM GoNest Project Setup Script for Windows
REM This script creates a new GoNest project with the recommended structure

echo ================================
echo    GoNest Project Setup Script
echo ================================
echo.

REM Check if Go is installed
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] Go is not installed. Please install Go 1.21+ first.
    pause
    exit /b 1
)

REM Check if Git is installed
where git >nul 2>nul
if %errorlevel% neq 0 (
    echo [WARNING] Git is not installed. Git repository will not be initialized.
    set GIT_AVAILABLE=false
) else (
    set GIT_AVAILABLE=true
)

REM Get project details
set /p PROJECT_NAME="Enter project name: "
if "%PROJECT_NAME%"=="" (
    echo [ERROR] Project name cannot be empty
    pause
    exit /b 1
)

REM Convert to lowercase and replace spaces with hyphens
set PROJECT_NAME=%PROJECT_NAME: =%
set PROJECT_NAME=%PROJECT_NAME: =%

set /p MODULE_NAME="Enter Go module name (default: %PROJECT_NAME%): "
if "%MODULE_NAME%"=="" set MODULE_NAME=%PROJECT_NAME%

set /p PROJECT_DESCRIPTION="Enter project description: "
if "%PROJECT_DESCRIPTION%"=="" set PROJECT_DESCRIPTION=A GoNest application

set /p AUTHOR_NAME="Enter author name: "
if "%AUTHOR_NAME%"=="" set AUTHOR_NAME=Your Name

set /p AUTHOR_EMAIL="Enter author email: "
if "%AUTHOR_EMAIL%"=="" set AUTHOR_EMAIL=your.email@example.com

echo [INFO] Creating project structure...

REM Create project directory
mkdir "%PROJECT_NAME%"
cd "%PROJECT_NAME%"

REM Create directory structure
mkdir cmd\server
mkdir internal\modules
mkdir internal\config
mkdir internal\shared\middleware
mkdir internal\shared\utils
mkdir internal\shared\constants
mkdir pkg
mkdir docs
mkdir scripts
mkdir tests
mkdir bin
mkdir build
mkdir deployments

echo [SUCCESS] Project structure created

echo [INFO] Initializing Go module...
go mod init %MODULE_NAME%
echo [SUCCESS] Go module initialized

REM Initialize Git repository if available
if "%GIT_AVAILABLE%"=="true" (
    echo [INFO] Initializing Git repository...
    git init
    
    REM Create .gitignore
    echo # GoNest Project > .gitignore
    echo # Binaries >> .gitignore
    echo bin/ >> .gitignore
    echo dist/ >> .gitignore
    echo *.exe >> .gitignore
    echo *.dll >> .gitignore
    echo *.so >> .gitignore
    echo *.dylib >> .gitignore
    echo. >> .gitignore
    echo # Test binary, built with 'go test -c' >> .gitignore
    echo *.test >> .gitignore
    echo. >> .gitignore
    echo # Output of the go coverage tool >> .gitignore
    echo *.out >> .gitignore
    echo. >> .gitignore
    echo # Dependency directories >> .gitignore
    echo vendor/ >> .gitignore
    echo. >> .gitignore
    echo # Go workspace file >> .gitignore
    echo go.work >> .gitignore
    echo. >> .gitignore
    echo # Environment files >> .gitignore
    echo .env >> .gitignore
    echo .env.local >> .gitignore
    echo .env.*.local >> .gitignore
    echo. >> .gitignore
    echo # IDE files >> .gitignore
    echo .vscode/ >> .gitignore
    echo .idea/ >> .gitignore
    echo *.swp >> .gitignore
    echo *.swo >> .gitignore
    echo. >> .gitignore
    echo # OS files >> .gitignore
    echo .DS_Store >> .gitignore
    echo Thumbs.db >> .gitignore
    echo. >> .gitignore
    echo # Logs >> .gitignore
    echo *.log >> .gitignore
    echo. >> .gitignore
    echo # Build artifacts >> .gitignore
    echo build/ >> .gitignore
    echo tmp/ >> .gitignore
    echo. >> .gitignore
    echo # Coverage reports >> .gitignore
    echo coverage/ >> .gitignore
    echo. >> .gitignore
    echo # Database files >> .gitignore
    echo *.db >> .gitignore
    echo *.sqlite >> .gitignore
    echo. >> .gitignore
    echo # Backup files >> .gitignore
    echo *.bak >> .gitignore
    echo *.backup >> .gitignore
    echo. >> .gitignore
    echo # Archive files >> .gitignore
    echo *.tar.gz >> .gitignore
    echo *.zip >> .gitignore
    
    echo [SUCCESS] Git repository initialized
) else (
    echo [WARNING] Skipping Git initialization (Git not available)
)

echo [INFO] Creating main application file...

REM Create main.go
echo package main > cmd\server\main.go
echo. >> cmd\server\main.go
echo import ( >> cmd\server\main.go
echo 	"context" >> cmd\server\main.go
echo 	"log" >> cmd\server\main.go
echo 	"os" >> cmd\server\main.go
echo. >> cmd\server\main.go
echo 	"github.com/joho/godotenv" >> cmd\server\main.go
echo 	"github.com/sirupsen/logrus" >> cmd\server\main.go
echo 	gonest "GoNest/gonest" >> cmd\server\main.go
echo 	"%MODULE_NAME%/internal/config" >> cmd\server\main.go
echo ) >> cmd\server\main.go
echo. >> cmd\server\main.go
echo func main() { >> cmd\server\main.go
echo 	// Load environment variables >> cmd\server\main.go
echo 	if err := godotenv.Load(); err != nil { >> cmd\server\main.go
echo 		log.Println("No .env file found, using system environment variables") >> cmd\server\main.go
echo 	} >> cmd\server\main.go
echo. >> cmd\server\main.go
echo 	// Load configuration >> cmd\server\main.go
echo 	cfg := config.LoadConfig() >> cmd\server\main.go
echo. >> cmd\server\main.go
echo 	// Initialize logger >> cmd\server\main.go
echo 	logger := logrus.New() >> cmd\server\main.go
echo 	logger.SetLevel(logrus.InfoLevel) >> cmd\server\main.go
echo 	logger.SetFormatter(&logrus.JSONFormatter{}) >> cmd\server\main.go
echo. >> cmd\server\main.go
echo 	// Create application with configuration >> cmd\server\main.go
echo 	app := gonest.NewApplication(). >> cmd\server\main.go
echo 		Config(&gonest.Config{ >> cmd\server\main.go
echo 			Port:        cfg.Server.Port, >> cmd\server\main.go
echo 			Host:        cfg.Server.Host, >> cmd\server\main.go
echo 			Environment: cfg.Server.Environment, >> cmd\server\main.go
echo 			LogLevel:    cfg.Server.LogLevel, >> cmd\server\main.go
echo 		}). >> cmd\server\main.go
echo 		Logger(logger). >> cmd\server\main.go
echo 		Build() >> cmd\server\main.go
echo. >> cmd\server\main.go
echo 	// Register application lifecycle hooks >> cmd\server\main.go
echo 	app.LifecycleManager.RegisterHook( >> cmd\server\main.go
echo 		gonest.EventApplicationStart, >> cmd\server\main.go
echo 		gonest.LifecycleHookFunc(func(ctx context.Context) error { >> cmd\server\main.go
echo 			logger.Info("🚀 Application starting up...") >> cmd\server\main.go
echo 			logger.Infof("🌍 Server running on %%s:%%s", cfg.Server.Host, cfg.Server.Port) >> cmd\server\main.go
echo 			return nil >> cmd\server\main.go
echo 		}), >> cmd\server\main.go
echo 		gonest.PriorityHigh, >> cmd\server\main.go
echo 	) >> cmd\server\main.go
echo. >> cmd\server\main.go
echo 	app.LifecycleManager.RegisterHook( >> cmd\server\main.go
echo 		gonest.EventApplicationStop, >> cmd\server\main.go
echo 		gonest.LifecycleHookFunc(func(ctx context.Context) error { >> cmd\server\main.go
echo 			logger.Info("🛑 Application shutting down...") >> cmd\server\main.go
echo 			return nil >> cmd\server\main.go
echo 		}), >> cmd\server\main.go
echo 		gonest.PriorityHigh, >> cmd\server\main.go
echo 	) >> cmd\server\main.go
echo. >> cmd\server\main.go
echo 	// Start the application >> cmd\server\main.go
echo 	logger.Info("Starting GoNest application...") >> cmd\server\main.go
echo 	if err := app.Start(); err != nil { >> cmd\server\main.go
echo 		logger.Fatal("Failed to start application:", err) >> cmd\server\main.go
echo 	} >> cmd\server\main.go
echo } >> cmd\server\main.go

echo [SUCCESS] Main application file created

echo [INFO] Creating configuration file...

REM Create config.go
echo package config > internal\config\config.go
echo. >> internal\config\config.go
echo import ( >> internal\config\config.go
echo 	"os" >> internal\config\config.go
echo 	"strconv" >> internal\config\config.go
echo ) >> internal\config\config.go
echo. >> internal\config\config.go
echo // Config holds application configuration >> internal\config\config.go
echo type Config struct { >> internal\config\config.go
echo 	Server   ServerConfig >> internal\config\config.go
echo 	Database DatabaseConfig >> internal\config\config.go
echo 	JWT      JWTConfig >> internal\config\config.go
echo 	Redis    RedisConfig >> internal\config\config.go
echo } >> internal\config\config.go
echo. >> internal\config\config.go
echo // ServerConfig holds server configuration >> internal\config\config.go
echo type ServerConfig struct { >> internal\config\config.go
echo 	Port        string >> internal\config\config.go
echo 	Host        string >> internal\config\config.go
echo 	Environment string >> internal\config\config.go
echo 	LogLevel    string >> internal\config\config.go
echo } >> internal\config\config.go
echo. >> internal\config\config.go
echo // DatabaseConfig holds database configuration >> internal\config\config.go
echo type DatabaseConfig struct { >> internal\config\config.go
echo 	Host     string >> internal\config\config.go
echo 	Port     int >> internal\config\config.go
echo 	Name     string >> internal\config\config.go
echo 	Username string >> internal\config\config.go
echo 	Password string >> internal\config\config.go
echo } >> internal\config\config.go
echo. >> internal\config\config.go
echo // JWTConfig holds JWT configuration >> internal\config\config.go
echo type JWTConfig struct { >> internal\config\config.go
echo 	Secret     string >> internal\config\config.go
echo 	Expiration int64 >> internal\config\config.go
echo } >> internal\config\config.go
echo. >> internal\config\config.go
echo // RedisConfig holds Redis configuration >> internal\config\config.go
echo type RedisConfig struct { >> internal\config\config.go
echo 	Host     string >> internal\config\config.go
echo 	Port     int >> internal\config\config.go
echo 	Password string >> internal\config\config.go
echo 	DB       int >> internal\config\config.go
echo } >> internal\config\config.go
echo. >> internal\config\config.go
echo // LoadConfig loads configuration from environment variables >> internal\config\config.go
echo func LoadConfig() *Config { >> internal\config\config.go
echo 	return &Config{ >> internal\config\config.go
echo 		Server: ServerConfig{ >> internal\config\config.go
echo 			Port:        getEnv("PORT", "8080"), >> internal\config\config.go
echo 			Host:        getEnv("HOST", "localhost"), >> internal\config\config.go
echo 			Environment: getEnv("ENV", "development"), >> internal\config\config.go
echo 			LogLevel:    getEnv("LOG_LEVEL", "info"), >> internal\config\config.go
echo 		}, >> internal\config\config.go
echo 		Database: DatabaseConfig{ >> internal\config\config.go
echo 			Host:     getEnv("DB_HOST", "localhost"), >> internal\config\config.go
echo 			Port:     getEnvAsInt("DB_PORT", 27017), >> internal\config\config.go
echo 			Name:     getEnv("DB_NAME", "gonest"), >> internal\config\config.go
echo 			Username: getEnv("DB_USERNAME", ""), >> internal\config\config.go
echo 			Password: getEnv("DB_PASSWORD", ""), >> internal\config\config.go
echo 		}, >> internal\config\config.go
echo 		JWT: JWTConfig{ >> internal\config\config.go
echo 			Secret:     getEnv("JWT_SECRET", "your-secret-key"), >> internal\config\config.go
echo 			Expiration: getEnvAsInt64("JWT_EXPIRATION", 86400), // 24 hours >> internal\config\config.go
echo 		}, >> internal\config\config.go
echo 		Redis: RedisConfig{ >> internal\config\config.go
echo 			Host:     getEnv("REDIS_HOST", "localhost"), >> internal\config\config.go
echo 			Port:     getEnvAsInt("REDIS_PORT", 6379), >> internal\config\config.go
echo 			Password: getEnv("REDIS_PASSWORD", ""), >> internal\config\config.go
echo 			DB:       getEnvAsInt("REDIS_DB", 0), >> internal\config\config.go
echo 		}, >> internal\config\config.go
echo 	} >> internal\config\config.go
echo } >> internal\config\config.go
echo. >> internal\config\config.go
echo // Helper functions >> internal\config\config.go
echo func getEnv(key, defaultValue string) string { >> internal\config\config.go
echo 	if value := os.Getenv(key); value != "" { >> internal\config\config.go
echo 		return value >> internal\config\config.go
echo 	} >> internal\config\config.go
echo 	return defaultValue >> internal\config\config.go
echo } >> internal\config\config.go
echo. >> internal\config\config.go
echo func getEnvAsInt(key string, defaultValue int) int { >> internal\config\config.go
echo 	if value := os.Getenv(key); value != "" { >> internal\config\config.go
echo 		if intValue, err := strconv.Atoi(value); err == nil { >> internal\config\config.go
echo 			return intValue >> internal\config\config.go
echo 		} >> internal\config\config.go
echo 	} >> internal\config\config.go
echo 	return defaultValue >> internal\config\config.go
echo } >> internal\config\config.go
echo. >> internal\config\config.go
echo func getEnvAsInt64(key string, defaultValue int64) int64 { >> internal\config\config.go
echo 	if value := os.Getenv(key); value != "" { >> internal\config\config.go
echo 		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil { >> internal\config\config.go
echo 			return intValue >> internal\config\config.go
echo 		} >> internal\config\config.go
echo 	} >> internal\config\config.go
echo 	return defaultValue >> internal\config\config.go
echo } >> internal\config\config.go

echo [SUCCESS] Configuration file created

echo [INFO] Creating environment file...

REM Create .env file
echo # Server Configuration > .env
echo PORT=8080 >> .env
echo HOST=localhost >> .env
echo ENV=development >> .env
echo LOG_LEVEL=info >> .env
echo. >> .env
echo # Database Configuration >> .env
echo DB_HOST=localhost >> .env
echo DB_PORT=27017 >> .env
echo DB_NAME=gonest >> .env
echo DB_USERNAME= >> .env
echo DB_PASSWORD= >> .env
echo. >> .env
echo # JWT Configuration >> .env
echo JWT_SECRET=your-super-secret-jwt-key-change-in-production >> .env
echo JWT_EXPIRATION=86400 >> .env
echo. >> .env
echo # Redis Configuration >> .env
echo REDIS_HOST=localhost >> .env
echo REDIS_PORT=6379 >> .env
echo REDIS_PASSWORD= >> .env
echo REDIS_DB=0 >> .env

echo [SUCCESS] Environment file created

echo [INFO] Creating README file...

REM Create README.md
echo # %PROJECT_NAME% > README.md
echo. >> README.md
echo %PROJECT_DESCRIPTION% >> README.md
echo. >> README.md
echo ## 🚀 Quick Start >> README.md
echo. >> README.md
echo ### Prerequisites >> README.md
echo. >> README.md
echo - Go 1.21+ >> README.md
echo - Git (optional) >> README.md
echo. >> README.md
echo ### Installation >> README.md
echo. >> README.md
echo 1. Clone the repository: >> README.md
echo ```bash >> README.md
echo git clone ^<your-repo-url^> >> README.md
echo cd %PROJECT_NAME% >> README.md
echo ``` >> README.md
echo. >> README.md
echo 2. Install dependencies: >> README.md
echo ```bash >> README.md
echo go mod tidy >> README.md
echo ``` >> README.md
echo. >> README.md
echo 3. Set up environment variables: >> README.md
echo ```bash >> README.md
echo cp .env.example .env >> README.md
echo # Edit .env with your configuration >> README.md
echo ``` >> README.md
echo. >> README.md
echo 4. Run the application: >> README.md
echo ```bash >> README.md
echo go run cmd\server\main.go >> README.md
echo ``` >> README.md
echo. >> README.md
echo ## 📁 Project Structure >> README.md
echo. >> README.md
echo ``` >> README.md
echo %PROJECT_NAME%/ >> README.md
echo ├── cmd/ >> README.md
echo │   └── server/          # Application entry point >> README.md
echo ├── internal/ >> README.md
echo │   ├── modules/         # Feature modules >> README.md
echo │   ├── config/          # Configuration management >> README.md
echo │   └── shared/          # Shared utilities >> README.md
echo ├── pkg/                 # Public packages >> README.md
echo ├── docs/                # Documentation >> README.md
echo ├── scripts/             # Build and deployment scripts >> README.md
echo ├── tests/               # Integration tests >> README.md
echo └── deployments/         # Deployment configurations >> README.md
echo ``` >> README.md
echo. >> README.md
echo ## 🎯 Features >> README.md
echo. >> README.md
echo - Modular architecture inspired by NestJS >> README.md
echo - Dependency injection >> README.md
echo - Built-in validation >> README.md
echo - JWT authentication >> README.md
echo - MongoDB integration >> README.md
echo - Redis caching >> README.md
echo - WebSocket support >> README.md
echo - Rate limiting >> README.md
echo - Comprehensive testing utilities >> README.md
echo. >> README.md
echo ## 📚 Documentation >> README.md
echo. >> README.md
echo - [GoNest Framework Documentation](https://github.com/ulims/GoNest) >> README.md
echo - [Architecture Guide](ARCHITECTURE.md) >> README.md
echo - [API Reference](docs/API.md) >> README.md
echo. >> README.md
echo ## 🤝 Contributing >> README.md
echo. >> README.md
echo 1. Fork the repository >> README.md
echo 2. Create a feature branch >> README.md
echo 3. Make your changes >> README.md
echo 4. Add tests >> README.md
echo 5. Submit a pull request >> README.md
echo. >> README.md
echo ## 📄 License >> README.md
echo. >> README.md
echo This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. >> README.md
echo. >> README.md
echo ## 👨‍💻 Author >> README.md
echo. >> README.md
echo **%AUTHOR_NAME%** - [%AUTHOR_EMAIL%](mailto:%AUTHOR_EMAIL%) >> README.md
echo. >> README.md
echo --- >> README.md
echo. >> README.md
echo Built with [GoNest](https://github.com/ulims/GoNest) - The Go framework that brings NestJS elegance to Go! 🚀 >> README.md

echo [SUCCESS] README file created

echo [INFO] Installing GoNest dependencies...

REM Install dependencies
go get github.com/ulims/GoNest
go get github.com/labstack/echo/v4
go get github.com/sirupsen/logrus
go get github.com/go-playground/validator/v10
go get github.com/golang-jwt/jwt/v5
go get github.com/gorilla/websocket
go get github.com/redis/go-redis/v9
go get go.mongodb.org/mongo-driver/mongo
go get github.com/joho/godotenv

REM Download and tidy dependencies
go mod tidy

echo [SUCCESS] Dependencies installed

REM Create initial commit if Git is available
if "%GIT_AVAILABLE%"=="true" (
    echo [INFO] Creating initial commit...
    git add .
    git commit -m "Initial commit: GoNest project setup"
    echo [SUCCESS] Initial commit created
)

echo.
echo [SUCCESS] Project setup complete!
echo.
echo 🎉 Your GoNest project '%PROJECT_NAME%' has been created successfully!
echo.
echo Next steps:
echo 1. cd %PROJECT_NAME%
echo 2. Review and customize the generated files
echo 3. Update the GoNest import path in main.go
echo 4. Run 'go mod tidy' to ensure dependencies are correct
echo 5. Start building your modules!
echo.
echo Documentation:
echo   - README.md - Project overview and setup instructions
echo   - ARCHITECTURE.md - GoNest architecture guide
echo   - docs/ - Additional documentation
echo.
echo Happy coding with GoNest! 🚀
echo.
pause

