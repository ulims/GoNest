// Package main demonstrates the GoNest framework usage
package main

import (
	"fmt"
	"log"

	"GoNest/gonest"

	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("🚀 GoNest Framework Demo")

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create application
	app := gonest.NewApplication().
		Config(&gonest.Config{
			Port:        "8080",
			Host:        "localhost",
			Environment: "development",
		}).
		Logger(logger).
		Build()

	logger.Info("✅ GoNest application created successfully!")
	logger.Info("🌍 Framework is ready to use!")
	logger.Infof("📱 Application configured for %s:%s", app.Config.Host, app.Config.Port)

	// Note: This is a demo - in a real application you would:
	// 1. Register modules
	// 2. Set up routes
	// 3. Start the server with app.Start()

	log.Println("Demo completed successfully!")
}
