package logger

import (
	"fmt"
	"os"
	"time"
)

// Logger provides verbose logging functionality
type CLILogger interface {
	GetVerbose() bool
	GetDebug() bool
	GetOutput() *os.File
	Info(format string, args ...interface{})
	Verbosef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Error(format string, args ...interface{})
	log(level, format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Logger provides verbose logging functionality
type Logger struct {
	Verbose bool     `mapstructure:"verbose"`
	Debug   bool     `mapstructure:"debug"`
	Output  *os.File `mapstructure:"output"`
}

func (c *Logger) GetVerbose() bool {
	return c.Verbose
}

func (c *Logger) GetDebug() bool {
	return c.Debug
}

func (c *Logger) GetOutput() *os.File {
	return c.Output
}

// NewLogger creates a new logger
func NewLogger(verbose, debug bool) CLILogger {
	return &Logger{
		Verbose: verbose,
		Debug:   debug,
		Output:  os.Stderr,
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

// Verbose logs a verbose message (only if verbose is enabled)
func (l *Logger) Verbosef(format string, args ...interface{}) {
	if l.Verbose {
		l.log("VERBOSE", format, args...)
	}
}

// Debug logs a debug message (only if debug is enabled)
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Debug {
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
	fmt.Fprintf(l.Output, "[%s] [%s] %s\n", timestamp, level, message)
}

// Infof is an alias for Info
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(format, args...)
}

// Errorf is an alias for Error
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(format, args...)
}
