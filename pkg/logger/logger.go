package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger is a global variable that holds the logger instance
var Log *logrus.Logger

// InitLogger initializes the logger with different log levels and log output
func InitLogger(debugMode bool) {
	// Create a new logger instance
	Log = logrus.New()

	// Set the default log level
	if debugMode {
		Log.SetLevel(logrus.DebugLevel) // For debugging purposes
	} else {
		Log.SetLevel(logrus.InfoLevel) // Default to info in production
	}

	// Set the log output format
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,         // Include timestamp in log entries
		TimestampFormat: time.RFC3339, // Set the timestamp format
	})

	// Log to both file and console
	Log.SetOutput(os.Stdout)

	// Optionally, create a log file for logging in production (if needed)
	if !debugMode {
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			Log.SetOutput(logFile)
		} else {
			Log.Warn("Could not open log file, using default stderr")
		}
	}
}

// Info logs an info message
func Info(message string, args ...interface{}) {
	Log.Infof(message, args...)
}

// Error logs an error message
func Error(message string, args ...interface{}) {
	Log.Errorf(message, args...)
}

// Debug logs a debug message
func Debug(message string, args ...interface{}) {
	Log.Debugf(message, args...)
}

// Warn logs a warning message
func Warn(message string, args ...interface{}) {
	Log.Warnf(message, args...)
}

// Fatal logs a fatal message and exits the application
func Fatal(message string, args ...interface{}) {
	Log.Fatalf(message, args...)
	os.Exit(1)
}
