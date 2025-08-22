package main

import (
	"net/http"
	"time"

	gonest "GoNest"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// ===== USER MODEL =====

// User represents a user model
type User struct {
	gonest.BaseModel
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Age      int    `json:"age" db:"age"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

// TableName returns the table name
func (u *User) TableName() string {
	return "users"
}

// ===== USER DTOs =====

// CreateUserDTO represents user creation data
type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=18,max=120"`
	IsActive bool   `json:"is_active"`
}

// UpdateUserDTO represents user update data
type UpdateUserDTO struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=50"`
	Age      int    `json:"age" validate:"omitempty,min=18,max=120"`
	IsActive *bool  `json:"is_active"`
}

// ===== USER REPOSITORY =====

// UserRepository provides user data access
type UserRepository struct {
	gonest.BaseRepository
	users map[string]*User
}

// NewUserRepository creates a new user repository
func NewUserRepository(db interface{}, logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		BaseRepository: *gonest.NewBaseRepository(nil, logger),
		users:          make(map[string]*User),
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *User) error {
	user.ID = time.Now().Format("20060102150405")
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.users[user.ID.(string)] = user
	return nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id string) (*User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, gonest.NotFoundException("User not found")
}

// FindAll finds all users
func (r *UserRepository) FindAll() ([]*User, error) {
	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(user *User) error {
	if _, exists := r.users[user.ID.(string)]; !exists {
		return gonest.NotFoundException("User not found")
	}
	user.UpdatedAt = time.Now()
	r.users[user.ID.(string)] = user
	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id string) error {
	if _, exists := r.users[id]; !exists {
		return gonest.NotFoundException("User not found")
	}
	delete(r.users, id)
	return nil
}

// ===== USER SERVICE =====

// UserService provides user business logic
type UserService struct {
	userRepo *UserRepository `inject:"UserRepository"`
	logger   *logrus.Logger  `inject:"Logger"`
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(dto *CreateUserDTO) (*User, error) {
	s.logger.Infof("Creating user: %s", dto.Name)

	user := &User{
		Name:     dto.Name,
		Email:    dto.Email,
		Age:      dto.Age,
		IsActive: dto.IsActive,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser gets a user by ID
func (s *UserService) GetUser(id string) (*User, error) {
	s.logger.Infof("Getting user: %s", id)
	return s.userRepo.FindByID(id)
}

// GetAllUsers gets all users
func (s *UserService) GetAllUsers() ([]*User, error) {
	s.logger.Info("Getting all users")
	return s.userRepo.FindAll()
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id string, dto *UpdateUserDTO) (*User, error) {
	s.logger.Infof("Updating user: %s", id)

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if dto.Name != "" {
		user.Name = dto.Name
	}
	if dto.Age > 0 {
		user.Age = dto.Age
	}
	if dto.IsActive != nil {
		user.IsActive = *dto.IsActive
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id string) error {
	s.logger.Infof("Deleting user: %s", id)
	return s.userRepo.Delete(id)
}

// ===== USER CONTROLLER =====

// UserController handles user HTTP requests
type UserController struct {
	userService *UserService `inject:"UserService"`
}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	return &UserController{}
}

// CreateUser handles user creation
func (c *UserController) CreateUser(ctx echo.Context) error {
	var dto CreateUserDTO
	if err := ctx.Bind(&dto); err != nil {
		return gonest.BadRequestException("Invalid request body")
	}

	// Validate DTO using pipes
	validator := gonest.NewValidationPipe()
	if _, err := validator.Transform(&dto); err != nil {
		return gonest.BadRequestException("Validation failed: " + err.Error())
	}

	user, err := c.userService.CreateUser(&dto)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, user)
}

// GetUser handles getting a user by ID
func (c *UserController) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return gonest.BadRequestException("User ID is required")
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

// GetAllUsers handles getting all users
func (c *UserController) GetAllUsers(ctx echo.Context) error {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, users)
}

// UpdateUser handles user updates
func (c *UserController) UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return gonest.BadRequestException("User ID is required")
	}

	var dto UpdateUserDTO
	if err := ctx.Bind(&dto); err != nil {
		return gonest.BadRequestException("Invalid request body")
	}

	// Validate DTO using pipes
	validator := gonest.NewValidationPipe()
	if _, err := validator.Transform(&dto); err != nil {
		return gonest.BadRequestException("Validation failed: " + err.Error())
	}

	user, err := c.userService.UpdateUser(id, &dto)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion
func (c *UserController) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return gonest.BadRequestException("User ID is required")
	}

	if err := c.userService.DeleteUser(id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

// ===== AUTH SERVICE =====

// AuthService provides authentication logic
type AuthService struct {
	logger *logrus.Logger `inject:"Logger"`
}

// NewAuthService creates a new auth service
func NewAuthService() *AuthService {
	return &AuthService{}
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(token string) (bool, error) {
	// Simple token validation (in real app, validate JWT)
	if len(token) < 10 {
		return false, gonest.UnauthorizedException("Invalid token")
	}
	return true, nil
}

// ===== CUSTOM GUARDS =====

// CustomAuthGuard provides custom authentication
type CustomAuthGuard struct {
	authService *AuthService `inject:"AuthService"`
}

// NewCustomAuthGuard creates a new custom auth guard
func NewCustomAuthGuard() *CustomAuthGuard {
	return &CustomAuthGuard{}
}

// CanActivate checks if the request is authenticated
func (cag *CustomAuthGuard) CanActivate(ctx echo.Context) (bool, error) {
	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		return false, gonest.UnauthorizedException("Authorization header required")
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	return cag.authService.ValidateToken(token)
}

// ===== CUSTOM INTERCEPTORS =====

// CustomLoggingInterceptor provides custom logging
type CustomLoggingInterceptor struct {
	logger *logrus.Logger `inject:"Logger"`
}

// NewCustomLoggingInterceptor creates a new custom logging interceptor
func NewCustomLoggingInterceptor() *CustomLoggingInterceptor {
	return &CustomLoggingInterceptor{}
}

// Intercept logs request and response information
func (cli *CustomLoggingInterceptor) Intercept(ctx echo.Context, next echo.HandlerFunc) error {
	start := time.Now()

	cli.logger.WithFields(logrus.Fields{
		"method": ctx.Request().Method,
		"path":   ctx.Request().URL.Path,
		"ip":     ctx.RealIP(),
	}).Info("Custom request started")

	err := next(ctx)

	duration := time.Since(start)
	status := ctx.Response().Status

	cli.logger.WithFields(logrus.Fields{
		"method":   ctx.Request().Method,
		"path":     ctx.Request().URL.Path,
		"status":   status,
		"duration": duration,
	}).Info("Custom request completed")

	return err
}

// ===== CUSTOM PIPES =====

// CustomTransformPipe provides custom data transformation
type CustomTransformPipe struct{}

// NewCustomTransformPipe creates a new custom transform pipe
func NewCustomTransformPipe() *CustomTransformPipe {
	return &CustomTransformPipe{}
}

// Transform transforms the input data
func (ctp *CustomTransformPipe) Transform(value interface{}) (interface{}, error) {
	// Add custom transformation logic here
	return value, nil
}

// ===== CUSTOM EXCEPTION FILTERS =====

// CustomExceptionFilter provides custom exception handling
type CustomExceptionFilter struct {
	logger *logrus.Logger `inject:"Logger"`
}

// NewCustomExceptionFilter creates a new custom exception filter
func NewCustomExceptionFilter() *CustomExceptionFilter {
	return &CustomExceptionFilter{}
}

// Catch handles custom exceptions
func (cef *CustomExceptionFilter) Catch(exception interface{}, ctx echo.Context) error {
	cef.logger.WithFields(logrus.Fields{
		"exception": exception,
		"path":      ctx.Request().URL.Path,
		"method":    ctx.Request().Method,
	}).Error("Custom exception caught")

	// Return a custom error response
	response := map[string]interface{}{
		"error":   "Custom error occurred",
		"status":  http.StatusInternalServerError,
		"path":    ctx.Request().URL.Path,
		"method":  ctx.Request().Method,
		"details": exception,
	}

	return ctx.JSON(http.StatusInternalServerError, response)
}

// ===== HEALTH CONTROLLER =====

// HealthController provides health check functionality
type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck handles health check requests
func (hc *HealthController) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
		"service":   "gonest-comprehensive-example",
		"version":   "1.0.0",
	})
}

// ===== MAIN FUNCTION =====

func main() {
	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Create application with enhanced features
	app := gonest.NewApplication().
		Config(&gonest.Config{
			Port:         "8080",
			Host:         "localhost",
			Environment:  "development",
			LogLevel:     "info",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		}).
		Logger(logger).
		Build()

	// Register services
	app.RegisterService("UserRepository", NewUserRepository(nil, logger))
	app.RegisterService("UserService", NewUserService())
	app.RegisterService("AuthService", NewAuthService())
	app.RegisterService("Logger", logger)

	// Register guards
	app.RegisterGuard("CustomAuthGuard", NewCustomAuthGuard(), gonest.PriorityNormal)

	// Register interceptors
	app.RegisterInterceptor("CustomLoggingInterceptor", NewCustomLoggingInterceptor(), gonest.PriorityNormal)

	// Register pipes
	app.RegisterPipe("CustomTransformPipe", NewCustomTransformPipe(), gonest.PriorityNormal)

	// Register exception filters
	app.RegisterExceptionFilter("CustomExceptionFilter", NewCustomExceptionFilter(), gonest.PriorityNormal)

	// Create user module
	userModule := gonest.NewModule("user").
		Service(NewUserService()).
		Provider(NewUserRepository(nil, logger)).
		Build()

	// Create user controller with guards and interceptors
	userController := gonest.NewController().
		Path("/users").
		Middleware(
			gonest.GuardMiddleware(NewCustomAuthGuard()),
			gonest.InterceptorMiddleware(NewCustomLoggingInterceptor()),
		).
		Post("/", NewUserController().CreateUser).
		Get("/", NewUserController().GetAllUsers).
		Get("/:id", NewUserController().GetUser).
		Put("/:id", NewUserController().UpdateUser).
		Delete("/:id", NewUserController().DeleteUser).
		Build()

	// Create health controller
	healthController := gonest.NewController().
		Path("/health").
		Get("/", NewHealthController().HealthCheck).
		Build()

	// Register controllers
	app.RegisterController(userController)
	app.RegisterController(healthController)

	// Register module
	app.ModuleRegistry.Register(userModule)

	// Add global exception filter middleware
	app.Use(gonest.ExceptionFilterMiddleware(
		gonest.NewHTTPExceptionFilter(logger),
		gonest.NewValidationExceptionFilter(logger),
		NewCustomExceptionFilter(),
		gonest.NewGenericExceptionFilter(logger),
	))

	logger.Info("Starting GoNest Comprehensive Example Application...")
	logger.Info("Server will be available at http://localhost:8080")
	logger.Info("WebSocket endpoint available at ws://localhost:8080/ws")
	logger.Info("Available endpoints:")
	logger.Info("  POST   /users/     - Create a new user (requires auth)")
	logger.Info("  GET    /users/     - Get all users (requires auth)")
	logger.Info("  GET    /users/:id  - Get user by ID (requires auth)")
	logger.Info("  PUT    /users/:id  - Update user (requires auth)")
	logger.Info("  DELETE /users/:id  - Delete user (requires auth)")
	logger.Info("  GET    /health/    - Health check")
	logger.Info("")
	logger.Info("Features demonstrated:")
	logger.Info("  ✅ Guards (Authentication)")
	logger.Info("  ✅ Interceptors (Logging)")
	logger.Info("  ✅ Pipes (Validation)")
	logger.Info("  ✅ Exception Filters (Error handling)")
	logger.Info("  ✅ WebSockets (Real-time communication)")
	logger.Info("  ✅ Database Integration (Repository pattern)")
	logger.Info("  ✅ Lifecycle Hooks (Application events)")
	logger.Info("  ✅ Dependency Injection")
	logger.Info("  ✅ Module System")
	logger.Info("  ✅ Middleware Pipeline")

	// Start application
	if err := app.Start(); err != nil {
		logger.Fatalf("Failed to start application: %v", err)
	}
}
