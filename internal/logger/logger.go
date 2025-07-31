package logger

import (
	"b3-ingest/internal/infra/settings"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// LogLevel defines the log severity levels.
type LogLevel int

const (
	DEBUG   LogLevel = iota // Most detailed level, for debugging.
	INFO                    // General information about the application flow.
	WARNING                 // Events that may indicate a problem, but do not stop execution.
	ERROR                   // Errors that affect functionality, but the application can continue.
	FATAL                   // Critical errors that cause the application to terminate.
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger is the custom structure for logging.
type Logger struct {
	mu       sync.Mutex  // Mutex to ensure concurrency safety when writing logs.
	std      *log.Logger // Standard Go logger for formatted output.
	minLevel LogLevel    // Minimum log level to be displayed.
}

// NewLogger creates and returns a new Logger instance.
// writer: Where logs will be written (e.g., os.Stdout, a file).
// prefix: Prefix for each log line (e.g., "[APP]").
// flag: Flags for the standard logger (e.g., log.Ldate|log.Ltime|log.Lshortfile).
// minLevel: The minimum log level to display.
func NewLogger(writer io.Writer, prefix string, flag int, minLevel LogLevel) *Logger {
	return &Logger{
		std:      log.New(writer, prefix, flag),
		minLevel: minLevel,
	}
}

// SetMinLevel sets the minimum log level to be displayed.
func (l *Logger) SetMinLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// logf is the internal method to format and write logs, checking the level.
func (l *Logger) logf(level LogLevel, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.minLevel {
		return // Do not log if the level is below the configured minimum.
	}

	// Add the log level to the message prefix.
	message := fmt.Sprintf("[%s] %s", level.String(), fmt.Sprintf(format, v...))
	l.std.Output(3, message) // 3 to skip internal logger calls and show the correct file/line.
}

// Debug logs a message at DEBUG level.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.logf(DEBUG, format, v...)
}

// Info logs a message at INFO level.
func (l *Logger) Info(format string, v ...interface{}) {
	l.logf(INFO, format, v...)
}

// Warning logs a message at WARNING level.
func (l *Logger) Warning(format string, v ...interface{}) {
	l.logf(WARNING, format, v...)
}

// Error logs a message at ERROR level.
func (l *Logger) Error(format string, v ...interface{}) {
	l.logf(ERROR, format, v...)
}

// Fatal logs a message at FATAL level and terminates the application.
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logf(FATAL, format, v...)
	os.Exit(1) // Exit the program with exit code 1.
}

// Global usage example (optional, but common for convenience)
var defaultLogger *Logger

// InitDefaultLogger initializes the default logger after envs are loaded.
func InitDefaultLogger() {
	defaultLogger = NewLogger(os.Stdout, "["+settings.GetEnvs().APPName+"] ", log.Ldate|log.Ltime, INFO)
}

// GetDefaultLogger returns the instance of the default logger.
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// SetDefaultLogger allows replacing the default logger instance.
func SetDefaultLogger(l *Logger) {
	defaultLogger = l
}
