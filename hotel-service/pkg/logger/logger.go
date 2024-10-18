package logger

import (
	"log"
	"os"
)

// Logger interface defines the methods for logging at different levels.
type Logger interface {
	Info(msg string, keysAndValues ...interface{})  // Logs informational messages
	Error(msg string, keysAndValues ...interface{}) // Logs error messages
	Fatal(msg string, keysAndValues ...interface{}) // Logs fatal messages and exits
}

// simpleLogger struct implements the Logger interface.
type simpleLogger struct {
	infoLogger  *log.Logger // Logger for info messages
	errorLogger *log.Logger // Logger for error messages
}

// New creates a new instance of simpleLogger with configured loggers.
func New() Logger {
	return &simpleLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),  // Info logger writes to stdout
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile), // Error logger writes to stderr
	}
}

// Info logs an informational message.
func (l *simpleLogger) Info(msg string, keysAndValues ...interface{}) {
	l.infoLogger.Println(msg, keysAndValues) // Log the info message
}

// Error logs an error message.
func (l *simpleLogger) Error(msg string, keysAndValues ...interface{}) {
	l.errorLogger.Println(msg, keysAndValues) // Log the error message
}

// Fatal logs a fatal message and exits the program.
func (l *simpleLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.errorLogger.Fatal(msg, keysAndValues) // Log the fatal message and exit
}
