@echo off
setlocal enabledelayedexpansion

echo 🚀 Installing GoNest CLI...

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go 1.21+ first.
    echo    Visit: https://golang.org/doc/install
    pause
    exit /b 1
)

echo ✅ Go detected

REM Create temporary directory
set "TEMP_DIR=%TEMP%\gonest-install-%RANDOM%"
mkdir "%TEMP_DIR%"
cd /d "%TEMP_DIR%"

echo 📥 Cloning GoNest repository...
git clone --quiet https://github.com/ulims/GoNest.git
cd GoNest

echo 🔨 Building and installing CLI tool...
go install ./cmd/gonest
if %errorlevel% neq 0 (
    echo ❌ Failed to install CLI tool. Please check the error above.
    pause
    exit /b 1
)

REM Clean up
cd /d "%TEMP%"
rmdir /s /q "%TEMP_DIR%"

echo ✅ GoNest CLI installed successfully!
echo.
echo 🎯 Usage:
echo    gonest --help                    # Show help
echo    gonest new my-app               # Create new project
echo    gonest new my-api --template=api --strict  # With template and strict mode
echo.
echo 🚀 Happy coding with GoNest!
pause
