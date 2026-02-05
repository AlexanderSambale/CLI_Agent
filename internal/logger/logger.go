package logger

import (
	"fmt"
	"os"
	"time"
)

// Logger provides verbose logging functionality
type Logger struct {
	verbose bool
	debug   bool
	output  *os.File
}

// NewLogger creates a new logger
func NewLogger(verbose, debug bool) *Logger {
	return &Logger{
		verbose: verbose,
		debug:   debug,
		output:  os.Stderr,
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

// Verbose logs a verbose message (only if verbose is enabled)
func (l *Logger) Verbose(format string, args ...interface{}) {
	if l.verbose {
		l.log("VERBOSE", format, args...)
	}
}

// Debug logs a debug message (only if debug is enabled)
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.debug {
		l.log("DEBUG", format, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

// log is the internal logging method
func (l *Logger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "[%s] [%s] %s\n", timestamp, level, message)
}

// Infof is an alias for Info
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(format, args...)
}

// Verbosef is an alias for Verbose
func (l *Logger) Verbosef(format string, args ...interface{}) {
	l.Verbose(format, args...)
}

// Debugf is an alias for Debug
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(format, args...)
}

// Errorf is an alias for Error
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(format, args...)
}