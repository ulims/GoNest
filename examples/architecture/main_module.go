package main

import (
	"architecture-example/modules/user"

	gonest "github.com/ulims/GoNest/gonest"

	"github.com/sirupsen/logrus"
)

// MainModule represents the main application module
type MainModule struct {
	*gonest.Module
}

// NewMainModule creates the main application module
func NewMainModule(logger *logrus.Logger) *MainModule {
	// Create user module
	userModule := user.NewUserModule(logger)

	// Create and return main module
	module := gonest.NewModule("MainModule").
		Import(userModule.Module).
		Build()

	return &MainModule{
		Module: module,
	}
}
