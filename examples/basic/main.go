package main

import (
	"net/http"
	"time"

	gonest "GoNest"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// User DTOs
type CreateUserDTO struct {
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=18,max=120"`
}

type UserResponseDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

// User model
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

// User Repository (simulated)
type UserRepository struct {
	users map[string]*User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*User),
	}
}

func (r *UserRepository) Create(user *User) error {
	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) FindByID(id string) (*User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, echo.NewHTTPError(http.StatusNotFound, "User not found")
}

func (r *UserRepository) FindAll() ([]*User, error) {
	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

// User Service
type UserService struct {
	userRepo *UserRepository `inject:"UserRepository"`
	logger   *logrus.Logger  `inject:"Logger"`
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(dto *CreateUserDTO) (*User, error) {
	s.logger.Infof("Creating user: %s", dto.Name)

	user := &User{
		ID:        generateID(),
		Name:      dto.Name,
		Email:     dto.Email,
		Age:       dto.Age,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(id string) (*User, error) {
	s.logger.Infof("Fetching user with ID: %s", id)
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetAllUsers() ([]*User, error) {
	s.logger.Info("Fetching all users")
	return s.userRepo.FindAll()
}

// User Controller
type UserController struct {
	userService *UserService `inject:"UserService"`
}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) CreateUser(ctx echo.Context) error {
	var dto CreateUserDTO
	if err := ctx.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate DTO
	validator := gonest.NewDTOValidator()
	if err := validator.Validate(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	user, err := c.userService.CreateUser(&dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID is required")
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetAllUsers(ctx echo.Context) error {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}

	return ctx.JSON(http.StatusOK, users)
}

// Health Controller
type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
		"service":   "gonest-example",
	})
}

// Utility function
func generateID() string {
	return time.Now().Format("20060102150405")
}

func main() {
	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Create application
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
	app.RegisterService("UserRepository", NewUserRepository())
	app.RegisterService("UserService", NewUserService())
	app.RegisterService("Logger", logger)

	// Create user module
	userModule := gonest.NewModule("user").
		Service(NewUserService()).
		Provider(NewUserRepository()).
		Build()

	// Create user controller
	userController := gonest.NewController().
		Path("/users").
		Post("/", NewUserController().CreateUser).
		Get("/", NewUserController().GetAllUsers).
		Get("/:id", NewUserController().GetUser).
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
	app.Module(userModule)

	// Add some basic middleware
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Infof("%s %s", c.Request().Method, c.Request().URL.Path)
			return next(c)
		}
	})

	logger.Info("Starting GoNest example application...")
	logger.Info("Server will be available at http://localhost:8080")
	logger.Info("Available endpoints:")
	logger.Info("  POST   /users/     - Create a new user")
	logger.Info("  GET    /users/     - Get all users")
	logger.Info("  GET    /users/:id  - Get user by ID")
	logger.Info("  GET    /health/    - Health check")

	// Start application
	if err := app.Start(); err != nil {
		logger.Fatalf("Failed to start application: %v", err)
	}
}
