# GoNest Setup Scripts

This directory contains automated setup scripts to quickly initialize new GoNest projects with the recommended structure and configuration.

## 🚀 Available Scripts

### 1. **setup-project.sh** (Linux/macOS)
Bash script for Unix-like systems that creates a complete GoNest project structure.

### 2. **setup-project.bat** (Windows)
Batch script for Windows systems that creates a complete GoNest project structure.

## 📋 Prerequisites

Before running the setup scripts, ensure you have:

- **Go 1.21+** installed and in your PATH
- **Git** installed (optional, but recommended)
- **Basic Go knowledge** (packages, modules, structs)

### Verify Go Installation

```bash
go version
# Should output: go version go1.21.x windows/amd64 (or similar)
```

## 🔧 Usage

### Linux/macOS

```bash
# Make the script executable
chmod +x scripts/setup-project.sh

# Run the script
./scripts/setup-project.sh

# Or run with a project name
./scripts/setup-project.sh my-awesome-app
```

### Windows

```cmd
# Run the batch script
scripts\setup-project.bat
```

## 📁 What Gets Created

The setup scripts automatically create:

```
my-gonest-app/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── modules/                 # Feature modules directory
│   ├── config/
│   │   └── config.go           # Configuration management
│   └── shared/
│       ├── middleware/          # Shared middleware
│       ├── utils/              # Utility functions
│       └── constants/          # Application constants
├── pkg/                        # Public packages
├── docs/                       # Documentation
├── scripts/                    # Build and deployment scripts
├── tests/                      # Integration tests
├── bin/                        # Binary output
├── build/                      # Build artifacts
├── deployments/                # Deployment configurations
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── .env                        # Environment variables
├── .gitignore                  # Git ignore file
├── README.md                   # Project documentation
├── Makefile                    # Build automation
├── Dockerfile                  # Container configuration
└── docker-compose.yml          # Multi-service setup
```

## 🎯 Generated Files

### 1. **Main Application** (`cmd/server/main.go`)
- Complete GoNest application setup
- Environment variable loading
- Configuration management
- Lifecycle hooks registration

### 2. **Configuration** (`internal/config/config.go`)
- Environment-based configuration
- Server, database, JWT, and Redis settings
- Helper functions for type conversion

### 3. **Environment File** (`.env`)
- Pre-configured environment variables
- Development-ready defaults
- Secure JWT secret placeholder

### 4. **Project Documentation** (`README.md`)
- Project overview and setup instructions
- Feature list and architecture description
- Contributing guidelines

### 5. **Build Automation** (`Makefile`)
- Common development commands
- Build, test, and run targets
- Docker integration

### 6. **Containerization** (`Dockerfile`, `docker-compose.yml`)
- Multi-stage Docker build
- Development environment with MongoDB and Redis
- Non-root user security

## 🔄 Setup Process

The scripts perform these steps automatically:

1. **Prerequisites Check**: Verify Go installation and version
2. **Project Details**: Collect project name, description, and author info
3. **Directory Structure**: Create the recommended folder hierarchy
4. **Go Module**: Initialize Go module with proper naming
5. **Git Repository**: Initialize Git and create comprehensive `.gitignore`
6. **Source Files**: Generate main application and configuration files
7. **Dependencies**: Install GoNest and all required packages
8. **Documentation**: Create project README and setup guides
9. **Build Tools**: Add Makefile and Docker configuration
10. **Initial Commit**: Create first Git commit (if Git available)

## 🎨 Customization

After running the setup script, you can customize:

### 1. **Update Import Paths**
Edit `cmd/server/main.go` to use your actual GoNest import path:
```go
gonest "github.com/ulims/GoNest"
```

### 2. **Modify Configuration**
Update `internal/config/config.go` to match your specific needs:
- Database connection settings
- JWT configuration
- Redis settings

### 3. **Add Your Modules**
Create new modules in `internal/modules/` following the pattern:
```
modules/{feature}/
├── {feature}_module.go     # Module definition
├── {feature}_service.go    # Business logic
└── {feature}_controller.go # HTTP handlers
```

### 4. **Environment Variables**
Modify `.env` file with your specific configuration:
- Database credentials
- JWT secrets
- Service endpoints

## 🚨 Troubleshooting

### Common Issues

1. **Permission Denied** (Linux/macOS)
   ```bash
   chmod +x scripts/setup-project.sh
   ```

2. **Go Not Found**
   - Ensure Go is installed and in your PATH
   - Verify with `go version`

3. **Git Not Available**
   - Script will continue without Git initialization
   - You can manually initialize Git later

4. **Dependencies Fail**
   - Run `go mod tidy` manually
   - Check your Go version (1.21+ required)

### Getting Help

- Check the [GoNest Documentation](docs/DOCUMENTATION.md)
- Review the [Architecture Guide](ARCHITECTURE.md)
- Look at the [Examples](examples/) directory

## 🔮 Next Steps

After successful setup:

1. **Review Generated Code**: Understand the structure and patterns
2. **Customize Configuration**: Update settings for your environment
3. **Add Your Modules**: Start building your application features
4. **Run the Application**: Test with `make run` or `go run cmd/server/main.go`
5. **Build and Deploy**: Use `make build` and Docker for production

## 🤝 Contributing

To improve these setup scripts:

1. Fork the repository
2. Create a feature branch
3. Make your improvements
4. Test thoroughly on multiple platforms
5. Submit a pull request

## 📄 License

These scripts are part of the GoNest framework and are licensed under the MIT License.

---

**Happy coding with GoNest! 🚀**

