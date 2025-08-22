package user

import (
	"net/http"

	gonest "github.com/ulims/GoNest/gonest"

	"github.com/labstack/echo/v4"
)

// UserController handles HTTP requests for user operations
type UserController struct {
	userService *UserService
}

// NewUserController creates a new user controller
func NewUserController(userService *UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser handles user creation
func (c *UserController) CreateUser(ctx echo.Context) error {
	var req struct {
		Username  string `json:"username" validate:"required,min=3,max=50"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8"`
		FirstName string `json:"first_name" validate:"required,min=2,max=50"`
		LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	}

	if err := ctx.Bind(&req); err != nil {
		return gonest.BadRequestException("Invalid request body")
	}

	// Validate request
	if err := gonest.ValidateStruct(req, nil); err != nil {
		return gonest.BadRequestException(err.Error())
	}

	// Create user
	user, err := c.userService.CreateUser(
		req.Username,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		return gonest.BadRequestException(err.Error())
	}

	return ctx.JSON(http.StatusCreated, user)
}

// GetUser handles user retrieval by ID
func (c *UserController) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return gonest.BadRequestException("User ID is required")
	}

	user, err := c.userService.GetUser(id)
	if err != nil {
		return gonest.NotFoundException("User not found")
	}

	return ctx.JSON(http.StatusOK, user)
}

// GetUserByUsername handles user retrieval by username
func (c *UserController) GetUserByUsername(ctx echo.Context) error {
	username := ctx.Param("username")
	if username == "" {
		return gonest.BadRequestException("Username is required")
	}

	user, err := c.userService.GetUserByUsername(username)
	if err != nil {
		return gonest.NotFoundException("User not found")
	}

	return ctx.JSON(http.StatusOK, user)
}

// UpdateUser handles user updates
func (c *UserController) UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return gonest.BadRequestException("User ID is required")
	}

	var req struct {
		FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
		LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
		Email     string `json:"email" validate:"omitempty,email"`
		Status    string `json:"status" validate:"omitempty,oneof=active inactive suspended"`
	}

	if err := ctx.Bind(&req); err != nil {
		return gonest.BadRequestException("Invalid request body")
	}

	// Validate request
	if err := gonest.ValidateStruct(req, nil); err != nil {
		return gonest.BadRequestException(err.Error())
	}

	// Update user
	user, err := c.userService.UpdateUser(id, req.FirstName, req.LastName, req.Email, req.Status)
	if err != nil {
		return gonest.BadRequestException(err.Error())
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
		return gonest.NotFoundException("User not found")
	}

	return ctx.NoContent(http.StatusNoContent)
}

// ListUsers handles user listing
func (c *UserController) ListUsers(ctx echo.Context) error {
	users, err := c.userService.ListUsers()
	if err != nil {
		return gonest.BadRequestException("Failed to retrieve users")
	}

	return ctx.JSON(http.StatusOK, users)
}

// GetUserProfile handles user profile retrieval
func (c *UserController) GetUserProfile(ctx echo.Context) error {
	// In a real app, get user ID from JWT token
	userID := ctx.Get("user_id").(string)
	if userID == "" {
		return gonest.UnauthorizedException("User not authenticated")
	}

	user, err := c.userService.GetUser(userID)
	if err != nil {
		return gonest.NotFoundException("User not found")
	}

	return ctx.JSON(http.StatusOK, user)
}

// UpdateUserProfile handles user profile updates
func (c *UserController) UpdateUserProfile(ctx echo.Context) error {
	// In a real app, get user ID from JWT token
	userID := ctx.Get("user_id").(string)
	if userID == "" {
		return gonest.UnauthorizedException("User not authenticated")
	}

	var req struct {
		FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
		LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
		Email     string `json:"email" validate:"omitempty,email"`
	}

	if err := ctx.Bind(&req); err != nil {
		return gonest.BadRequestException("Invalid request body")
	}

	// Validate request
	if err := gonest.ValidateStruct(req, nil); err != nil {
		return gonest.BadRequestException(err.Error())
	}

	// Update user
	user, err := c.userService.UpdateUser(userID, req.FirstName, req.LastName, req.Email, "")
	if err != nil {
		return gonest.BadRequestException(err.Error())
	}

	return ctx.JSON(http.StatusOK, user)
}
