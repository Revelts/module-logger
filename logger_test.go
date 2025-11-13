package logger

import (
	"testing"
)

func TestGetInstance(t *testing.T) {
	// Test that GetInstance returns the same instance
	instance1 := GetInstance()
	instance2 := GetInstance()

	if instance1 != instance2 {
		t.Error("GetInstance should return the same instance (singleton)")
	}
}

func TestInit(t *testing.T) {
	logger := GetInstance()

	// Test initialization with empty DSN (should work, just won't send to Sentry)
	err := logger.Init("", "test")
	if err != nil {
		t.Errorf("Init should not fail with empty DSN: %v", err)
	}

	// Test that re-initialization fails
	err = logger.Init("", "test")
	if err == nil {
		t.Error("Re-initialization should fail")
	}
}

func TestLogLevels(t *testing.T) {
	logger := GetInstance()

	// Reset for testing
	logger.mu.Lock()
	logger.initialized = false
	logger.mu.Unlock()

	// Initialize without Sentry for testing
	logger.Init("", "test")

	// Test all log levels (these should not panic)
	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	// Test with fields
	fields := map[string]interface{}{
		"user_id": 123,
		"action":  "test",
	}
	logger.Info("Info with fields", fields)
	logger.Error("Error with fields", fields)
}

func TestErrorWithErr(t *testing.T) {
	logger := GetInstance()

	// Reset for testing
	logger.mu.Lock()
	logger.initialized = false
	logger.mu.Unlock()

	logger.Init("", "test")

	// Test with error
	testErr := &CustomError{Message: "test error"}
	logger.ErrorWithErr(testErr, "Something went wrong")

	// Test with nil error
	logger.ErrorWithErr(nil, "No error provided")
}

type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

