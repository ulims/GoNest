package user

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all user module routes
func (m *UserModule) RegisterRoutes(e *echo.Echo) {
	// Create route group for user module
	userGroup := e.Group("/users")
	
	// Register routes with the controller
	userGroup.POST("", m.userController.CreateUser)
	userGroup.GET("/:id", m.userController.GetUser)
	userGroup.GET("", m.userController.ListUsers)
}
