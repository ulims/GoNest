package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"test-project/internal/modules/user"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create Echo instance
	e := echo.New()

	// Initialize and register the User module
	// This demonstrates GoNest's modular architecture!
	userModule := user.NewUserModule(logger)
	
	// Register module routes
	userModule.RegisterRoutes(e)

	// Add a health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
			"message": "🚀 GoNest Application with Modular Architecture is running!",
			"modules": "User module is active and ready",
		})
	})

	// Add root route
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "🚀 GoNest Application with Modular Architecture!")
	})

	// Start server
	addr := "localhost:8080"
	logger.Infof("🚀 Starting GoNest application on %s", addr)
	logger.Info("📁 User module is loaded and ready!")
	logger.Info("🎯 Try: POST /users, GET /users, GET /users/:id")
	
	if err := e.Start(addr); err != nil {
		logger.Fatal("Failed to start application:", err)
	}
}
