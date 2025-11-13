# Module Logger

A Go logging library with console output and Sentry integration for error tracking. Uses singleton pattern for easy use across multiple repositories.

## Features

- ✅ Console logging for all log levels (Debug, Info, Warn, Error)
- ✅ Automatic Sentry integration for error-level logs
- ✅ Singleton pattern for easy usage across projects
- ✅ Structured logging with custom fields
- ✅ Thread-safe implementation

## Installation

```bash
go get module-logger
```

## Quick Start

### 1. Initialize the Logger

Initialize the logger once in your application (typically in `main.go`):

```go
package main

import (
    "module-logger"
    "os"
)

func main() {
    // Initialize logger with Sentry DSN and environment
    logger := logger.GetInstance()
    err := logger.Init(
        os.Getenv("SENTRY_DSN"),           // Your Sentry DSN
        os.Getenv("ENVIRONMENT"),           // e.g., "production", "development"
    )
    if err != nil {
        panic(err)
    }
    
    // Always flush before application exit
    defer logger.Flush()
    
    // Your application code...
}
```

### 2. Use the Logger

```go
package main

import "module-logger"

func main() {
    logger := logger.GetInstance()
    
    // Log different levels
    logger.Debug("Debug information")
    logger.Info("Application started")
    logger.Warn("This is a warning")
    logger.Error("This is an error - will be sent to Sentry")
    
    // Log with custom fields
    fields := map[string]interface{}{
        "user_id": 12345,
        "action":  "login",
        "ip":      "192.168.1.1",
    }
    logger.Info("User logged in", fields)
    
    // Log errors with error objects
    err := someFunction()
    if err != nil {
        logger.ErrorWithErr(err, "Failed to process request", fields)
    }
}
```

## API Reference

### Initialization

#### `GetInstance() *Logger`
Returns the singleton logger instance.

#### `Init(sentryDSN string, environment string) error`
Initializes the logger with Sentry configuration.
- `sentryDSN`: Your Sentry DSN (leave empty string to disable Sentry)
- `environment`: Environment name (e.g., "production", "development", "staging")

### Logging Methods

#### `Debug(message string, fields ...map[string]interface{})`
Logs a debug message to console.

#### `Info(message string, fields ...map[string]interface{})`
Logs an info message to console.

#### `Warn(message string, fields ...map[string]interface{})`
Logs a warning message to console.

#### `Error(message string, fields ...map[string]interface{})`
Logs an error message to console **and sends it to Sentry**.

#### `ErrorWithErr(err error, message string, fields ...map[string]interface{})`
Logs an error with an error object. Includes error details in Sentry.

#### `Flush()`
Flushes any pending Sentry events. **Always call this before application exit.**

## Example: Complete Application

```go
package main

import (
    "fmt"
    "module-logger"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // Initialize logger
    logger := logger.GetInstance()
    err := logger.Init(
        os.Getenv("SENTRY_DSN"),
        os.Getenv("ENVIRONMENT"),
    )
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize logger: %v", err))
    }
    defer logger.Flush()
    
    logger.Info("Application starting")
    
    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        logger.Info("Shutting down gracefully")
        logger.Flush()
        os.Exit(0)
    }()
    
    // Your application logic
    if err := runApplication(); err != nil {
        logger.ErrorWithErr(err, "Application failed", map[string]interface{}{
            "component": "main",
        })
        os.Exit(1)
    }
}

func runApplication() error {
    logger := logger.GetInstance()
    
    logger.Info("Processing request")
    
    // Simulate an error
    if someCondition {
        return fmt.Errorf("something went wrong")
    }
    
    return nil
}
```

## Environment Variables

```bash
# Required for Sentry integration
export SENTRY_DSN="https://your-sentry-dsn@sentry.io/project-id"

# Environment name
export ENVIRONMENT="production"  # or "development", "staging"
```

## Why Singleton Pattern?

The singleton pattern is chosen for this library because:

1. **Ease of Use**: Simple API - just call `logger.GetInstance().Info(...)` from anywhere
2. **Single Configuration**: Configure once at application startup
3. **Standard Practice**: Most Go logging libraries (logrus, zap) use similar patterns
4. **Cross-Repository**: Easy to use across multiple repositories without passing logger instances

## Thread Safety

The logger is thread-safe and can be used concurrently from multiple goroutines.

## License

MIT

# module-logger
