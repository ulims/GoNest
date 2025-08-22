package main

import (
	"context"
	"github.com/sirupsen/logrus"
	gonest "github.com/ulims/GoNest"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create application
	app := gonest.NewApplication().
		Config(&gonest.Config{
			Port:        "8080",
			Host:        "localhost",
			Environment: "development",
		}).
		Logger(logger).
		Build()

	// Register your modules here
	// app.ModuleRegistry.Register(yourModule.GetModule())

	// Start the application
	if err := app.Start(); err != nil {
		logger.Fatal("Failed to start application:", err)
	}
}
