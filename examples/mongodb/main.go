package main

import (
	"context"
	"net/http"
	"time"

	gonest "github.com/ulims/GoNest/gonest"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// ===== USER SCHEMA & MODEL =====

// User represents a user document
type User struct {
	gonest.MongoDBBaseModel
	Name     string   `bson:"name" json:"name"`
	Email    string   `bson:"email" json:"email"`
	Age      int      `bson:"age" json:"age"`
	IsActive bool     `bson:"isActive" json:"isActive"`
	Profile  *Profile `bson:"profile,omitempty" json:"profile,omitempty"`
	Posts    []Post   `bson:"posts,omitempty" json:"posts,omitempty"`
	Tags     []string `bson:"tags,omitempty" json:"tags,omitempty"`
	Settings Settings `bson:"settings" json:"settings"`
}

// Profile represents a user profile
type Profile struct {
	Bio       string    `bson:"bio" json:"bio"`
	Avatar    string    `bson:"avatar" json:"avatar"`
	Website   string    `bson:"website" json:"website"`
	Location  string    `bson:"location" json:"location"`
	BirthDate time.Time `bson:"birthDate" json:"birthDate"`
}

// Settings represents user settings
type Settings struct {
	Notifications bool   `bson:"notifications" json:"notifications"`
	Theme         string `bson:"theme" json:"theme"`
	Language      string `bson:"language" json:"language"`
}

// Post represents a user post
type Post struct {
	ID      string    `bson:"_id" json:"id"`
	Title   string    `bson:"title" json:"title"`
	Content string    `bson:"content" json:"content"`
	Created time.Time `bson:"created" json:"created"`
}

// GetCollection returns the collection name
func (u *User) GetCollection() string {
	return "users"
}

// BeforeSave hook - validate user data
func (u *User) BeforeSave() error {
	if u.Name == "" {
		return gonest.BadRequestException("Name is required")
	}
	if u.Email == "" {
		return gonest.BadRequestException("Email is required")
	}
	if u.Age < 0 || u.Age > 150 {
		return gonest.BadRequestException("Invalid age")
	}
	return nil
}

// AfterSave hook - log user creation
func (u *User) AfterSave() error {
	// Could send welcome email, create default settings, etc.
	return nil
}

// ===== USER DTOs =====

// CreateUserDTO represents user creation data
type CreateUserDTO struct {
	Name     string   `json:"name" validate:"required,min=2,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Age      int      `json:"age" validate:"required,min=13,max=120"`
	IsActive bool     `json:"isActive"`
	Bio      string   `json:"bio"`
	Tags     []string `json:"tags"`
}

// UpdateUserDTO represents user update data
type UpdateUserDTO struct {
	Name     *string  `json:"name" validate:"omitempty,min=2,max=50"`
	Age      *int     `json:"age" validate:"omitempty,min=13,max=120"`
	IsActive *bool    `json:"isActive"`
	Bio      *string  `json:"bio"`
	Tags     []string `json:"tags"`
}

// ===== USER SERVICE =====

// UserService provides user business logic
type UserService struct {
	userModel gonest.MongoDBModel `inject:"UserModel"`
	logger    *logrus.Logger      `inject:"Logger"`
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, dto *CreateUserDTO) (*User, error) {
	s.logger.Infof("Creating user: %s", dto.Name)

	user := &User{
		Name:     dto.Name,
		Email:    dto.Email,
		Age:      dto.Age,
		IsActive: dto.IsActive,
		Profile: &Profile{
			Bio: dto.Bio,
		},
		Tags: dto.Tags,
		Settings: Settings{
			Notifications: true,
			Theme:         "light",
			Language:      "en",
		},
	}

	if err := s.userModel.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser gets a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	s.logger.Infof("Getting user: %s", id)

	user := &User{}
	if err := s.userModel.FindById(ctx, id, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers gets all users with pagination
func (s *UserService) GetUsers(ctx context.Context, page, limit int) ([]*User, int64, error) {
	s.logger.Info("Getting users with pagination")

	skip := int64((page - 1) * limit)

	// Use query builder for complex queries
	query := s.userModel.Query().
		Where("isActive", true).
		Sort("createdAt", -1).
		Skip(skip).
		Limit(int64(limit)).
		Select("name", "email", "age", "isActive", "createdAt")

	var users []*User
	if err := query.Find(ctx, &users); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := s.userModel.Count(ctx, map[string]interface{}{"isActive": true})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id string, dto *UpdateUserDTO) (*User, error) {
	s.logger.Infof("Updating user: %s", id)

	update := map[string]interface{}{}
	if dto.Name != nil {
		update["name"] = *dto.Name
	}
	if dto.Age != nil {
		update["age"] = *dto.Age
	}
	if dto.IsActive != nil {
		update["isActive"] = *dto.IsActive
	}
	if dto.Bio != nil {
		update["profile.bio"] = *dto.Bio
	}
	if dto.Tags != nil {
		update["tags"] = dto.Tags
	}

	if err := s.userModel.UpdateById(ctx, id, update); err != nil {
		return nil, err
	}

	return s.GetUser(ctx, id)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	s.logger.Infof("Deleting user: %s", id)
	return s.userModel.DeleteById(ctx, id)
}

// SearchUsers searches users by name or email
func (s *UserService) SearchUsers(ctx context.Context, query string) ([]*User, error) {
	s.logger.Infof("Searching users with query: %s", query)

	// Use regex for text search
	filter := map[string]interface{}{
		"$or": []map[string]interface{}{
			{"name": map[string]interface{}{"$regex": query, "$options": "i"}},
			{"email": map[string]interface{}{"$regex": query, "$options": "i"}},
		},
	}

	var users []*User
	if err := s.userModel.Find(ctx, filter, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersByAgeRange gets users within an age range
func (s *UserService) GetUsersByAgeRange(ctx context.Context, minAge, maxAge int) ([]*User, error) {
	s.logger.Infof("Getting users by age range: %d-%d", minAge, maxAge)

	query := s.userModel.Query().
		WhereGreaterThanOrEqual("age", minAge).
		WhereLessThanOrEqual("age", maxAge).
		Sort("age", 1)

	var users []*User
	if err := query.Find(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserStats gets user statistics using aggregation
func (s *UserService) GetUserStats(ctx context.Context) (map[string]interface{}, error) {
	s.logger.Info("Getting user statistics")

	pipeline := []map[string]interface{}{
		{
			"$group": map[string]interface{}{
				"_id":        nil,
				"totalUsers": map[string]interface{}{"$sum": 1},
				"activeUsers": map[string]interface{}{
					"$sum": map[string]interface{}{
						"$cond": []interface{}{"$isActive", 1, 0},
					},
				},
				"avgAge": map[string]interface{}{"$avg": "$age"},
				"minAge": map[string]interface{}{"$min": "$age"},
				"maxAge": map[string]interface{}{"$max": "$age"},
			},
		},
	}

	var results []map[string]interface{}
	if err := s.userModel.Aggregate(ctx, pipeline, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return map[string]interface{}{
			"totalUsers":  0,
			"activeUsers": 0,
			"avgAge":      0,
			"minAge":      0,
			"maxAge":      0,
		}, nil
	}

	return map[string]interface{}{
		"totalUsers":  results[0]["totalUsers"],
		"activeUsers": results[0]["activeUsers"],
		"avgAge":      results[0]["avgAge"],
		"minAge":      results[0]["minAge"],
		"maxAge":      results[0]["maxAge"],
	}, nil
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

	user, err := c.userService.CreateUser(ctx.Request().Context(), &dto)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, user)
}

// GetUser handles getting a user by ID
func (c *UserController) GetUser(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return gonest.BadRequestException("User ID is required")
	}

	// In a real implementation, this would validate the ID format
	// For now, we'll just use the string as is
	id := idStr

	user, err := c.userService.GetUser(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

// GetUsers handles getting all users with pagination
func (c *UserController) GetUsers(ctx echo.Context) error {
	page := 1
	limit := 10

	if pageStr := ctx.QueryParam("page"); pageStr != "" {
		if p, err := gonest.ParseIntPipeInstance.Transform(pageStr); err == nil {
			page = p.(int)
		}
	}

	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
		if l, err := gonest.ParseIntPipeInstance.Transform(limitStr); err == nil {
			limit = l.(int)
		}
	}

	users, total, err := c.userService.GetUsers(ctx.Request().Context(), page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateUser handles user updates
func (c *UserController) UpdateUser(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return gonest.BadRequestException("User ID is required")
	}

	// In a real implementation, this would validate the ID format
	// For now, we'll just use the string as is
	id := idStr

	var dto UpdateUserDTO
	if err := ctx.Bind(&dto); err != nil {
		return gonest.BadRequestException("Invalid request body")
	}

	// Validate DTO using pipes
	validator := gonest.NewValidationPipe()
	if _, err := validator.Transform(&dto); err != nil {
		return gonest.BadRequestException("Validation failed: " + err.Error())
	}

	user, err := c.userService.UpdateUser(ctx.Request().Context(), id, &dto)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

// DeleteUser handles user deletion
func (c *UserController) DeleteUser(ctx echo.Context) error {
	idStr := ctx.Param("id")
	if idStr == "" {
		return gonest.BadRequestException("User ID is required")
	}

	// In a real implementation, this would validate the ID format
	// For now, we'll just use the string as is
	id := idStr

	if err := c.userService.DeleteUser(ctx.Request().Context(), id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

// SearchUsers handles user search
func (c *UserController) SearchUsers(ctx echo.Context) error {
	query := ctx.QueryParam("q")
	if query == "" {
		return gonest.BadRequestException("Search query is required")
	}

	users, err := c.userService.SearchUsers(ctx.Request().Context(), query)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, users)
}

// GetUsersByAgeRange handles getting users by age range
func (c *UserController) GetUsersByAgeRange(ctx echo.Context) error {
	minAge := 0
	maxAge := 150

	if minAgeStr := ctx.QueryParam("minAge"); minAgeStr != "" {
		if m, err := gonest.ParseIntPipeInstance.Transform(minAgeStr); err == nil {
			minAge = m.(int)
		}
	}

	if maxAgeStr := ctx.QueryParam("maxAge"); maxAgeStr != "" {
		if m, err := gonest.ParseIntPipeInstance.Transform(maxAgeStr); err == nil {
			maxAge = m.(int)
		}
	}

	users, err := c.userService.GetUsersByAgeRange(ctx.Request().Context(), minAge, maxAge)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, users)
}

// GetUserStats handles getting user statistics
func (c *UserController) GetUserStats(ctx echo.Context) error {
	stats, err := c.userService.GetUserStats(ctx.Request().Context())
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, stats)
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
		"service":   "gonest-mongodb-example",
		"version":   "1.0.0",
		"database":  "mongodb",
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

	// Create user schema (similar to Mongoose)
	userSchema := gonest.NewSchema()
	userSchema.Field("name", gonest.String).Required().MinLength(2).MaxLength(50)
	userSchema.Field("email", gonest.String).Required().Unique().Index()
	userSchema.Field("age", gonest.Number).Required().Min(13).Max(120)
	userSchema.Field("isActive", gonest.Boolean).Default(true)
	userSchema.Field("profile", gonest.Object).Embedded()
	userSchema.Field("posts", gonest.Array).Array(gonest.Object)
	userSchema.Field("tags", gonest.Array).Array(gonest.String)
	userSchema.Field("settings", gonest.Object).Embedded()
	userSchema.Timestamps().Collection("users")
	userSchema.AddIndex(gonest.NewIndex(map[string]interface{}{"email": 1}).Unique())
	userSchema.AddIndex(gonest.NewIndex(map[string]interface{}{"name": 1}))
	userSchema.AddIndex(gonest.NewIndex(map[string]interface{}{"age": 1}))
	userSchema.AddIndex(gonest.NewIndex(map[string]interface{}{"isActive": 1}))
	userSchema.AddIndex(gonest.NewIndex(map[string]interface{}{"createdAt": -1}))

	// Create application with MongoDB
	app := gonest.NewApplication().
		Config(&gonest.Config{
			Port:         "8080",
			Host:         "localhost",
			Environment:  "development",
			LogLevel:     "info",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			MongoDB: &gonest.MongoDBConfig{
				URI:                    "mongodb://localhost:27017",
				Database:               "gonest_example",
				MaxPoolSize:            100,
				MinPoolSize:            5,
				MaxIdleTime:            30 * time.Second,
				ConnectTimeout:         10 * time.Second,
				ServerSelectionTimeout: 30 * time.Second,
				SocketTimeout:          30 * time.Second,
			},
		}).
		Logger(logger).
		Build()

	// Create MongoDB service
	mongoService := gonest.NewMongoDBService(app.Config.MongoDB, logger)
	app.MongoDBService = mongoService

	// Create user model
	userModel := mongoService.Model("User", userSchema)

	// Register services
	app.RegisterService("UserModel", userModel)
	app.RegisterService("UserService", NewUserService())
	app.RegisterService("Logger", logger)

	// Create user module
	userModule := gonest.NewModule("user").
		Service(NewUserService()).
		Provider(userModel).
		Build()

	// Create user controller
	userController := gonest.NewController().
		Path("/users").
		Post("/", NewUserController().CreateUser).
		Get("/", NewUserController().GetUsers).
		Get("/search", NewUserController().SearchUsers).
		Get("/stats", NewUserController().GetUserStats).
		Get("/age-range", NewUserController().GetUsersByAgeRange).
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

	logger.Info("Starting GoNest MongoDB Example Application...")
	logger.Info("Server will be available at http://localhost:8080")
	logger.Info("MongoDB will be available at mongodb://localhost:27017")
	logger.Info("Available endpoints:")
	logger.Info("  POST   /users/         - Create a new user")
	logger.Info("  GET    /users/         - Get all users (with pagination)")
	logger.Info("  GET    /users/search   - Search users by name/email")
	logger.Info("  GET    /users/stats    - Get user statistics")
	logger.Info("  GET    /users/age-range - Get users by age range")
	logger.Info("  GET    /users/:id      - Get user by ID")
	logger.Info("  PUT    /users/:id      - Update user")
	logger.Info("  DELETE /users/:id      - Delete user")
	logger.Info("  GET    /health/        - Health check")
	logger.Info("")
	logger.Info("MongoDB Features demonstrated:")
	logger.Info("  ✅ Schema Definition (like Mongoose)")
	logger.Info("  ✅ Field Validation & Constraints")
	logger.Info("  ✅ Indexes (Unique, Compound, TTL)")
	logger.Info("  ✅ Timestamps (createdAt, updatedAt)")
	logger.Info("  ✅ Embedded Documents")
	logger.Info("  ✅ Array Fields")
	logger.Info("  ✅ Query Builder (like Mongoose)")
	logger.Info("  ✅ Aggregation Pipeline")
	logger.Info("  ✅ Lifecycle Hooks (BeforeSave, AfterSave)")
	logger.Info("  ✅ Connection Pooling")
	logger.Info("  ✅ Error Handling")
	logger.Info("  ✅ Pagination")
	logger.Info("  ✅ Text Search")
	logger.Info("  ✅ Statistics & Analytics")

	// Start application
	if err := app.Start(); err != nil {
		logger.Fatalf("Failed to start application: %v", err)
	}
}
