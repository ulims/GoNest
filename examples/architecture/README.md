# GoNest Architecture Example

This example demonstrates the **NestJS-style modular architecture** implemented in Go using the GoNest framework. The architecture follows the principle of **separation of concerns** and **modular design** that makes applications scalable, maintainable, and testable.

## 🏗️ Directory Structure

```
examples/architecture/
├── main.go                 # Application entry point
├── main_module.go          # Root module that imports all feature modules
├── go.mod                  # Go module definition
├── modules/                # Feature modules directory
│   └── user/              # User feature module
│       ├── user_module.go    # Module definition and registration
│       ├── user_service.go   # Business logic layer
│       └── user_controller.go # HTTP request handling
└── README.md              # This file
```

## 🔑 Key Architectural Principles

### 1. **Modular Design**
- Each feature is encapsulated in its own module
- Modules can be imported and reused across the application
- Clear separation between different business domains

### 2. **Layered Architecture**
- **Controller Layer**: Handles HTTP requests and responses
- **Service Layer**: Contains business logic and data operations
- **Model Layer**: Defines data structures and validation rules

### 3. **Dependency Injection**
- Services are injected into controllers
- Modules manage their own dependencies
- Clean, testable code structure

## 📦 Module Structure

### User Module (`modules/user/`)
The user module demonstrates a complete feature implementation:

- **`user_module.go`**: Defines the module, registers services and controllers
- **`user_service.go`**: Implements user business logic (CRUD operations)
- **`user_controller.go`**: Handles HTTP endpoints for user operations

## 🚀 Module Registration Order

The application follows a specific module registration order:

1. **Main Module** (`main_module.go`) - Root module that imports feature modules
2. **Feature Modules** - Individual business domain modules (e.g., `user`)

## 🎯 Key Features Demonstrated

### Service Layer (`user_service.go`)
- **Business Logic**: User creation, retrieval, update, deletion
- **Data Management**: In-memory storage with thread safety
- **Validation**: Input validation and error handling
- **Model Definition**: User entity with validation tags

### Controller Layer (`user_controller.go`)
- **HTTP Endpoints**: RESTful API endpoints
- **Request Validation**: Input validation using DTOs
- **Error Handling**: Proper HTTP status codes and error responses
- **Response Formatting**: JSON responses with appropriate status codes

### Module Layer (`user_module.go`)
- **Dependency Registration**: Services and controllers registration
- **Module Configuration**: Module-specific settings and dependencies
- **Integration**: Connects services and controllers

## 🔧 Running the Example

```bash
# Navigate to the architecture example directory
cd examples/architecture

# Build the application
go build .

# Run the application
./architecture-example.exe
```

## 📚 API Endpoints

Once running, the application provides these endpoints:

- `POST /users` - Create a new user
- `GET /users/:id` - Get user by ID
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user
- `GET /users` - List all users

## 🧪 Testing the Architecture

The example demonstrates:

1. **Module Registration**: How modules are created and registered
2. **Dependency Injection**: How services are injected into controllers
3. **Request Flow**: Complete request → controller → service → response flow
4. **Error Handling**: Proper error responses and validation
5. **Data Validation**: Input validation using struct tags

## 🎯 Learning Objectives

This architecture example teaches:

- **GoNest Framework Usage**: How to use the framework's module system
- **Modular Design**: How to structure applications with clear separation of concerns
- **Best Practices**: Following Go and NestJS architectural patterns
- **Scalability**: How to organize code for growth and maintenance

## 🔄 Extending the Example

To add new features:

1. **Create a new module** in `modules/` directory
2. **Follow the same structure**: `{feature}_module.go`, `{feature}_service.go`, `{feature}_controller.go`
3. **Register the module** in `main_module.go`
4. **Import dependencies** as needed

This architecture provides a solid foundation for building enterprise-grade applications with GoNest!
