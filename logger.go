package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

// Logger is the singleton logger instance
type Logger struct {
	initialized bool
	sentryDSN   string
	environment string
	mu          sync.RWMutex
}

var (
	instance *Logger
	once     sync.Once
)

// GetInstance returns the singleton logger instance
func GetInstance() *Logger {
	once.Do(func() {
		instance = &Logger{
			initialized: false,
		}
	})
	return instance
}

// Init initializes the logger with Sentry configuration
// sentryDSN: Your Sentry DSN (Data Source Name)
// environment: Environment name (e.g., "production", "development", "staging")
func (l *Logger) Init(sentryDSN string, environment string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.initialized {
		return fmt.Errorf("logger already initialized")
	}

	l.sentryDSN = sentryDSN
	l.environment = environment

	// Initialize Sentry
	if sentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			Environment:      environment,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			return fmt.Errorf("failed to initialize Sentry: %w", err)
		}
	}

	l.initialized = true
	return nil
}

// log prints to console and sends to Sentry if error level
func (l *Logger) log(level LogLevel, message string, fields ...map[string]interface{}) {
	if !l.initialized {
		// Fallback to basic logging if not initialized
		log.Printf("[%s] %s", level, message)
		return
	}

	// Prepare log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)

	// Add fields if provided
	if len(fields) > 0 && fields[0] != nil {
		for key, value := range fields[0] {
			logEntry += fmt.Sprintf(" %s=%v", key, value)
		}
	}

	// Always print to console
	fmt.Fprintln(os.Stdout, logEntry)

	// Send to Sentry only for error level
	if level == LevelError {
		l.sendToSentry(message, fields...)
	}
}

// sendToSentry sends error to Sentry
func (l *Logger) sendToSentry(message string, fields ...map[string]interface{}) {
	if l.sentryDSN == "" {
		return
	}

	// Create Sentry event
	event := sentry.NewEvent()
	event.Message = message
	event.Level = sentry.LevelError
	event.Environment = l.environment

	// Add extra fields if provided
	if len(fields) > 0 && fields[0] != nil {
		event.Extra = fields[0]
	}

	// Capture the event
	sentry.CaptureEvent(event)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(LevelDebug, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log(LevelInfo, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(LevelWarn, message, fields...)
}

// Error logs an error message and sends it to Sentry
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	l.log(LevelError, message, fields...)
}

// ErrorWithErr logs an error with an error object and sends it to Sentry
func (l *Logger) ErrorWithErr(err error, message string, fields ...map[string]interface{}) {
	if err == nil {
		l.Error(message, fields...)
		return
	}

	// Merge error into fields
	errorFields := make(map[string]interface{})
	if len(fields) > 0 && fields[0] != nil {
		for k, v := range fields[0] {
			errorFields[k] = v
		}
	}
	errorFields["error"] = err.Error()
	errorFields["error_type"] = fmt.Sprintf("%T", err)

	fullMessage := fmt.Sprintf("%s: %v", message, err)
	l.log(LevelError, fullMessage, errorFields)
}

// Flush flushes any pending Sentry events (should be called before application exit)
func (l *Logger) Flush() {
	if l.sentryDSN != "" {
		sentry.Flush(time.Second * 2)
	}
}
