package user

import (
	"github.com/sirupsen/logrus"
	gonest "github.com/ulims/GoNest/gonest"
)

// UserModule demonstrates GoNest's modular architecture
type UserModule struct {
	*gonest.Module
}

// NewUserModule creates a new user module with all its components
func NewUserModule(logger *logrus.Logger) *UserModule {
	// Create services
	userService := NewUserService(logger)
	
	// Create controllers
	userController := NewUserController(userService)
	
	// Create and return module - this is where the magic happens!
	module := gonest.NewModule("UserModule").
		Controller(userController).
		Service(userService).
		Build()
	
	return &UserModule{
		Module: module,
	}
}
