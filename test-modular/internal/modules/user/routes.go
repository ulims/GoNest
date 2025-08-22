package user

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all user module routes
func (m *UserModule) RegisterRoutes(e *echo.Echo) {
	// Create route group for user module
	userGroup := e.Group("/users")
	
	// Register routes with the controller
	userGroup.POST("", m.Controller.(*UserController).CreateUser)
	userGroup.GET("/:id", m.Controller.(*UserController).GetUser)
	userGroup.GET("", m.Controller.(*UserController).ListUsers)
}
