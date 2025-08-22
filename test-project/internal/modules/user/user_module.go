package user

import (
	"github.com/sirupsen/logrus"
)

// UserModule demonstrates modular architecture
type UserModule struct {
	userService    *UserService
	userController *UserController
	logger         *logrus.Logger
}

// NewUserModule creates a new user module with all its components
func NewUserModule(logger *logrus.Logger) *UserModule {
	// Create services
	userService := NewUserService(logger)
	
	// Create controllers
	userController := NewUserController(userService)
	
	// Create and return module - this is where the magic happens!
	return &UserModule{
		userService:    userService,
		userController: userController,
		logger:         logger,
	}
}
