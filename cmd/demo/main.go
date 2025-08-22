package main

import (
	"net/http"
	"time"

	gonest "github.com/ulims/GoNest/gonest"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Simple example to demonstrate GoNest framework

func main() {
	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Create application
	app := gonest.NewApplication().
		Config(&gonest.Config{
			Port:         "8080",
			Host:         "localhost",
			Environment:  "development",
			LogLevel:     "info",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		}).
		Logger(logger).
		Build()

	// Add a simple route
	app.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to GoNest Framework!",
			"author":  "Agbama Ulimhuka Akem",
			"email":   "ulimhukaakem@gmail.com",
			"time":    time.Now(),

		})
	})

	// Add health check
	app.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now(),
			"service":   "gonest-demo",
		})
	})

	logger.Info("Starting GoNest Demo Application...")
	logger.Info("Server will be available at http://localhost:8080")
	logger.Info("Available endpoints:")
	logger.Info("  GET    /       - Welcome message")
	logger.Info("  GET    /health - Health check")

	// Start application
	if err := app.Start(); err != nil {
		logger.Fatalf("Failed to start application: %v", err)
	}
}
