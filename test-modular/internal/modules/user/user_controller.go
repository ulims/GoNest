package user

import (
	"net/http"
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
		Username string `json:"username" validate:"required,min=3"`
		Email    string `json:"email" validate:"required,email"`
	}
	
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}
	
	// Create user
	user, err := c.userService.CreateUser(req.Username, req.Email)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	
	return ctx.JSON(http.StatusCreated, user)
}

// GetUser handles user retrieval by ID
func (c *UserController) GetUser(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
	}
	
	user, err := c.userService.GetUser(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}
	
	return ctx.JSON(http.StatusOK, user)
}

// ListUsers handles user listing
func (c *UserController) ListUsers(ctx echo.Context) error {
	users, err := c.userService.ListUsers()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve users",
		})
	}
	
	return ctx.JSON(http.StatusOK, users)
}
