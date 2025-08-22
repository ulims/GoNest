package user

import (
	gonest "github.com/ulims/GoNest/gonest"

	"github.com/sirupsen/logrus"
)

// UserModule represents the user module
type UserModule struct {
	*gonest.Module
}

// NewUserModule creates a new user module
func NewUserModule(logger *logrus.Logger) *UserModule {
	// Create services
	userService := NewUserService(logger)

	// Create controllers
	userController := NewUserController(userService)

	// Create and return module
	module := gonest.NewModule("UserModule").
		Controller(userController).
		Service(userService).
		Build()

	return &UserModule{
		Module: module,
	}
}
