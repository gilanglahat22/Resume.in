package utils

import (
	"log"
	"os"
)

var (
	// InfoLogger for informational logs
	InfoLogger *log.Logger
	// ErrorLogger for error logs
	ErrorLogger *log.Logger
	// DebugLogger for debug logs
	DebugLogger *log.Logger
)

// InitLoggers initializes the logger instances
func InitLoggers() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info logs an informational message
func Info(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	DebugLogger.Printf(format, v...)
} 