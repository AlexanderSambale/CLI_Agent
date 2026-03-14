package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"cli_agent/internal/config"
)

// Executor defines the interface for command execution
type Executor interface {
	// Execute runs a command and returns the result
	Execute(ctx context.Context, command string) (*Result, error)

	// GetEngine returns the engine prefix
	GetEngine() string
}

// Result represents the result of a command execution
type Result struct {
	ExitCode int           // Process exit code
	Stdout   string        // Standard output
	Stderr   string        // Standard error
	Duration time.Duration // Execution duration
}

// executor implements the Executor interface
type executor struct {
	engine  string
	timeout time.Duration
}

// NewExecutor creates a new executor with the given configuration
func NewExecutor(config *config.ExecutionConfig) Executor {
	return &executor{
		engine:  config.Engine,
		timeout: config.Timeout,
	}
}

// GetEngine returns the engine prefix
func (e *executor) GetEngine() string {
	return e.engine
}

// Execute runs a command and returns the result
func (e *executor) Execute(ctx context.Context, command string) (*Result, error) {
	// Create context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// Build the full command
	var fullCommand string
	if e.engine == "" {
		// No engine prefix, just run bash -c
		fullCommand = fmt.Sprintf("bash -c %s", quoteCommand(command))
	} else {
		// Prepend engine prefix
		fullCommand = fmt.Sprintf("%s %s", e.engine, quoteCommand(command))
	}

	// Create command with bash
	cmd := exec.CommandContext(cmdCtx, fullCommand)

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment variables
	cmd.Env = os.Environ()

	// Execute and measure time
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	// Build result
	result := &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}

	// Handle exit code
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else if ctx.Err() == context.DeadlineExceeded {
			result.ExitCode = -1
			return result, fmt.Errorf("command execution timed out after %v", e.timeout)
		} else {
			result.ExitCode = -1
		}
		return result, err
	}

	result.ExitCode = 0
	return result, nil
}

// quoteCommand properly quotes a command for shell execution
func quoteCommand(command string) string {
	// If the command is already quoted, return it as-is
	if (strings.HasPrefix(command, "\"") && strings.HasSuffix(command, "\"")) ||
		(strings.HasPrefix(command, "'") && strings.HasSuffix(command, "'")) {
		return command
	}

	// Otherwise, wrap in double quotes and escape existing double quotes
	escaped := strings.ReplaceAll(command, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", escaped)
}
