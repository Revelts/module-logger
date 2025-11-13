package main

import (
	"fmt"
	"module-logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize logger with Sentry DSN and environment
	logger := logger.GetInstance()
	
	// Get Sentry DSN from environment variable (or use empty string to disable Sentry)
	sentryDSN := os.Getenv("SENTRY_DSN")
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	
	err := logger.Init(sentryDSN, environment)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	
	// Always flush Sentry events before application exit
	defer logger.Flush()
	
	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		logger.Info("Shutting down gracefully")
		logger.Flush()
		os.Exit(0)
	}()
	
	// Example usage
	logger.Info("Application started")
	
	// Log different levels
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	
	// Log with custom fields
	fields := map[string]interface{}{
		"user_id": 12345,
		"action":  "login",
		"ip":      "192.168.1.1",
		"timestamp": time.Now().Unix(),
	}
	logger.Info("User logged in", fields)
	
	// Log an error (will be sent to Sentry)
	logger.Error("This is an error - will be sent to Sentry", map[string]interface{}{
		"component": "example",
		"severity":  "high",
	})
	
	// Log error with error object
	testErr := fmt.Errorf("database connection failed")
	logger.ErrorWithErr(testErr, "Failed to connect to database", map[string]interface{}{
		"database": "postgres",
		"host":     "localhost",
		"port":     5432,
	})
	
	// Simulate some work
	time.Sleep(2 * time.Second)
	logger.Info("Application running...")
	
	// Keep running until interrupted
	select {}
}

