// Package main demonstrates the GoNest framework usage
package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/ulims/GoNest/gonest"
)

func main() {
	fmt.Println("🚀 Testing GoNest Framework Import")

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

	logger.Info("✅ GoNest framework imported successfully!")
	logger.Info("🌍 Framework is ready to use!")
	logger.Infof("📱 Application configured for %s:%s", app.Config.Host, app.Config.Port)

	fmt.Println("🎉 Manual installation test successful!")
	fmt.Println("✅ Framework can be imported and used")
	fmt.Println("✅ Manual installation method now works!")
}
