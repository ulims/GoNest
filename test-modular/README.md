# test-modular

A GoNest application built with the GoNest framework.

## Template: basic

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
   ```bash
   go mod tidy
   ```

### Running the Application

```bash
# Development mode
go run cmd/server/main.go

# Production build
go build -o bin/app cmd/server/main.go
./bin/app
```

### Testing

```bash
go test ./...
```

## Project Structure

```
test-modular/
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
```

## Documentation

For more information about GoNest, visit the [documentation](https://github.com/ulims/GoNest).
