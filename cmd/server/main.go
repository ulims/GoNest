package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create Echo instance
	e := echo.New()

	// Add a simple route
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "🚀 GoNest Application is running!")
	})

	// Add health check route
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "healthy",
			"message": "GoNest application is running successfully",
		})
	})

	// Start server
	addr := "localhost:8080"
	logger.Infof("🚀 Starting GoNest application on %s", addr)

	if err := e.Start(addr); err != nil {
		logger.Fatal("Failed to start application:", err)
	}
}
