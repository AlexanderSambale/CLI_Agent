package executor

import (
	"bytes"
	"context"
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

	// Build the command with proper argument splitting
	var cmd *exec.Cmd
	if e.engine == "" {
		// No engine prefix, just run bash -c with the command
		cmd = exec.CommandContext(cmdCtx, "bash", "-c", command)
	} else {
		// Prepend engine prefix and split into arguments
		parts := strings.Fields(e.engine)
		if len(parts) == 0 {
			// Fallback to bash if engine is empty after splitting
			cmd = exec.CommandContext(cmdCtx, "bash", "-c", command)
		} else {
			// First part is the command, rest are arguments
			args := append(parts[1:], command)
			cmd = exec.CommandContext(cmdCtx, parts[0], args...)
		}
	}

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
		} else {
			result.ExitCode = -1
		}
		return result, err
	}

	result.ExitCode = 0
	return result, nil
}
