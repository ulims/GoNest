package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gonest "github.com/ulims/GoNest/gonest"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// User represents a user entity
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-" validate:"required,min=6"`
	Role      string    `json:"role" validate:"required,oneof=user admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserService handles user business logic
type UserService struct {
	users map[string]*User
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
	}
}

func (us *UserService) CreateUser(ctx context.Context, user *User) error {
	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.ID = fmt.Sprintf("user_%d", len(us.users)+1)

	// Store user
	us.users[user.ID] = user

	return nil
}

func (us *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	// Get from memory
	user, exists := us.users[id]
	if !exists {
		return nil, gonest.NotFoundException(fmt.Sprintf("User not found: %s", id))
	}

	return user, nil
}

// UserController handles HTTP requests for users
type UserController struct {
	userService *UserService
}

func NewUserController() *UserController {
	return &UserController{}
}

// CreateUser endpoint with validation
func (c *UserController) CreateUser(ctx echo.Context) error {
	var user User
	if err := ctx.Bind(&user); err != nil {
		return gonest.BadRequestException("Invalid request body")
	}

	// Validate user data
	validator := gonest.NewDTOValidator()
	if err := gonest.ValidateStruct(&user, validator); err != nil {
		return gonest.BadRequestException(fmt.Sprintf("Validation failed: %v", err))
	}

	if err := c.userService.CreateUser(ctx.Request().Context(), &user); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, user)
}

// GetUser endpoint
func (c *UserController) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return gonest.BadRequestException("User ID is required")
	}

	user, err := c.userService.GetUser(ctx.Request().Context(), id)
	if err != nil {
		return gonest.NotFoundException(fmt.Sprintf("User not found: %s", id))
	}

	return ctx.JSON(http.StatusOK, user)
}

// GetProfile returns current user profile
func (c *UserController) GetProfile(ctx echo.Context) error {
	user, err := gonest.GetCurrentUser(ctx)
	if err != nil {
		return gonest.UnauthorizedException("User not authenticated")
	}

	return ctx.JSON(http.StatusOK, user)
}

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Initialize JWT service
	jwtConfig := gonest.DefaultJWTConfig()
	authService := gonest.NewAuthService(jwtConfig, logger)

	// Initialize passport service
	passportService := gonest.NewPassportService(logger)

	// Create auth controller
	authController := gonest.NewAuthController(authService, passportService, logger)

	// Create module
	userModule := gonest.NewModule("UserModule").
		Controller(NewUserController()).
		Service(NewUserService()).
		Provider(authService).
		Build()

	// Create application
	app := gonest.NewApplication().
		Config(&gonest.Config{
			Port: "8080",
		}).
		Logger(logger).
		Build()

	// Register module
	app.ModuleRegistry.Register(userModule)

	// Register lifecycle hook for route setup
	app.LifecycleManager.RegisterHook(gonest.EventApplicationStart, gonest.LifecycleHookFunc(func(ctx context.Context) error {
		// Auth routes
		authGroup := app.Group("/auth")
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/refresh", authController.RefreshToken)

		// Protected routes
		apiGroup := app.Group("/api")
		apiGroup.Use(authService.JWTMiddleware())

		// User routes
		userGroup := apiGroup.Group("/users")

		// Inject dependencies into controller
		userController := &UserController{}
		app.ServiceRegistry.Inject(userController)

		userGroup.POST("", userController.CreateUser)
		userGroup.GET("/:id", userController.GetUser)
		userGroup.GET("/profile", userController.GetProfile)

		// Health check
		app.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status":    "ok",
				"timestamp": time.Now(),
				"features": map[string]bool{
					"authentication":       true,
					"validation":           true,
					"dependency_injection": true,
				},
			})
		})

		logger.Info("Advanced GoNest application configured with core features")
		return nil
	}), gonest.PriorityNormal)

	// Start the application
	if err := app.Start(); err != nil {
		logger.Fatal("Failed to start application:", err)
	}
}
