package main

import (
	"context"

	gonest "GoNest"

	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create application with configuration
	app := gonest.NewApplication().
		Config(&gonest.Config{
			Port:        "8080",
			Host:        "localhost",
			Environment: "development",
			LogLevel:    "info",
		}).
		Logger(logger).
		Build()

	// Register main module
	mainModule := NewMainModule(logger)
	app.ModuleRegistry.Register(mainModule.Module)

	// Register application lifecycle hooks
	app.LifecycleManager.RegisterHook(
		gonest.EventApplicationStart,
		gonest.LifecycleHookFunc(func(ctx context.Context) error {
			logger.Info("🚀 Application starting up...")
			logger.Info("📁 Architecture Example: Demonstrating NestJS-style modular structure")
			logger.Info("🏗️  Each module contains: models/, dto/, services/, controllers/")
			return nil
		}),
		gonest.PriorityHigh,
	)

	app.LifecycleManager.RegisterHook(
		gonest.EventApplicationStop,
		gonest.LifecycleHookFunc(func(ctx context.Context) error {
			logger.Info("🛑 Application shutting down...")
			return nil
		}),
		gonest.PriorityHigh,
	)

	// Start the application
	logger.Info("Starting GoNest Architecture Example...")
	logger.Info("This example demonstrates NestJS-style modular architecture")
	logger.Info("Check the modules/ folder structure for the complete example")

	if err := app.Start(); err != nil {
		logger.Fatal("Failed to start application:", err)
	}
}
