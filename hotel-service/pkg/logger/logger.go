package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger interface defines the methods for logging at different levels.
type Logger interface {
	Info(msg string, keysAndValues ...interface{})  // Log an informational message.
	Error(msg string, keysAndValues ...interface{}) // Log an error message.
	Fatal(msg string, keysAndValues ...interface{}) // Log a fatal message and exit the application.
}

// simpleLogger struct implements the Logger interface.
type simpleLogger struct {
	infoLogger  *log.Logger // Logger for info messages.
	errorLogger *log.Logger // Logger for error messages.
}

// New function creates a new instance of simpleLogger.
func New() Logger {
	return &simpleLogger{
		infoLogger:  log.New(os.Stdout, "", 0), // Initialize info logger to write to standard output.
		errorLogger: log.New(os.Stderr, "", 0), // Initialize error logger to write to standard error.
	}
}

// log method handles the actual logging logic for different levels.
func (l *simpleLogger) log(level, msg string, keysAndValues ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)             // Get the current timestamp.
	logMsg := fmt.Sprintf("%s %s %s", timestamp, level, msg) // Format the log message.
	if len(keysAndValues) > 0 {
		logMsg += fmt.Sprintf(" %v", keysAndValues) // Append any additional key-value pairs.
	}
	if level == "INFO" {
		l.infoLogger.Println(logMsg) // Log info messages.
	} else {
		l.errorLogger.Println(logMsg) // Log error and fatal messages.
	}
}

// Info method logs an informational message.
func (l *simpleLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log("INFO", msg, keysAndValues...)
}

// Error method logs an error message.
func (l *simpleLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log("ERROR", msg, keysAndValues...)
}

// Fatal method logs a fatal message and exits the application.
func (l *simpleLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log("FATAL", msg, keysAndValues...)
	os.Exit(1) // Exit the application with a non-zero status.
}
