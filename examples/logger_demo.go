package main

import (
	"errors"
	"os"
	"time"

	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

func main() {
	// Set environment variables for demo
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FORMAT", "json")

	// Initialize logger
	log := logger.InitLogger()

	// Basic logging examples
	log.Info("Logger demo started")
	log.Debug("This is a debug message with details", "user_id", "123", "action", "demo")
	log.Warn("This is a warning message", "warning_type", "demo")

	// Component-specific logging
	dbLogger := log.WithComponent("database")
	dbLogger.Info("Database connection established", "host", "localhost", "port", 9042)

	// Request logging
	requestLogger := log.WithRequest("GET", "/api/v1/posts", "req-123")
	requestLogger.Info("Processing API request")

	// Error logging
	err := errors.New("sample error for demonstration")
	errorLogger := log.WithError(err)
	errorLogger.Error("An error occurred during processing")

	// Specialized logging methods
	log.LogRequest("GET", "/api/v1/posts", "192.168.1.100", 200, 45)
	log.LogRequest("POST", "/api/v1/posts", "192.168.1.101", 400, 120)

	// Database operation logging
	log.LogDatabaseOperation("SELECT", "posts", 25, nil)
	log.LogDatabaseOperation("INSERT", "users", 1500, errors.New("connection timeout"))

	// Service operation logging
	log.LogServiceOperation("post", "create", true, 50)
	log.LogServiceOperation("user", "update", false, 200)

	// Authentication logging
	log.LogAuth("login", "user123", true, "")
	log.LogAuth("login", "user456", false, "invalid credentials")

	// Startup and shutdown logging
	startupConfig := map[string]interface{}{
		"version":     "1.0.0",
		"environment": "demo",
		"port":        8080,
		"debug":       true,
	}
	log.LogStartup("demo-service", startupConfig)

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	log.LogShutdown("demo-service", "demo completed successfully")

	// Demonstrate different log formats
	log.Info("Now switching to text format...")

	// Switch to text format
	os.Setenv("LOG_FORMAT", "text")
	textLogger := logger.InitLogger()

	textLogger.Info("This message is in text format")
	textLogger.WithComponent("demo").Warn("Text format warning message")
	textLogger.LogRequest("GET", "/api/v1/demo", "127.0.0.1", 200, 10)

	textLogger.Info("Logger demo completed successfully")
}
