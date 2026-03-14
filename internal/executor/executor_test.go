package executor

import (
	"context"
	"strings"
	"testing"
	"time"

	"cli_agent/internal/config"
)

const (
	EXECUTE_ERROR         = "Execute() error = %v"
	EXECUTE_EXIT_CODE     = "Execute() exitCode = %v, want 0"
	EXECUTE_ERROR_GOT_NIL = "Execute() expected error, got nil"
)

func TestNewExecutor(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.ExecutionConfig
		expected string
	}{
		{
			name: "empty engine",
			config: &config.ExecutionConfig{
				Engine:  "",
				Timeout: 30 * time.Second,
			},
			expected: "",
		},
		{
			name: "docker engine",
			config: &config.ExecutionConfig{
				Engine:  "docker run --rm ubuntu bash -c",
				Timeout: 30 * time.Second,
			},
			expected: "docker run --rm ubuntu bash -c",
		},
		{
			name: "podman engine",
			config: &config.ExecutionConfig{
				Engine:  "podman run --rm alpine sh -c",
				Timeout: 30 * time.Second,
			},
			expected: "podman run --rm alpine sh -c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutor(tt.config)
			if exec.GetEngine() != tt.expected {
				t.Errorf("NewExecutor() engine = %v, want %v", exec.GetEngine(), tt.expected)
			}
		})
	}
}

func TestExecuteSuccess(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "echo 'hello world'")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "hello world") {
		t.Errorf("Execute() stdout = %v, want to contain 'hello world'", result.Stdout)
	}

	if result.Stderr != "" {
		t.Errorf("Execute() stderr = %v, want empty", result.Stderr)
	}

	if result.Duration == 0 {
		t.Errorf("Execute() duration = %v, want > 0", result.Duration)
	}
}

func TestExecuteNonZeroExitCode(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "exit 42")
	if err == nil {
		t.Fatal(EXECUTE_ERROR_GOT_NIL)
	}

	if result.ExitCode != 42 {
		t.Errorf("Execute() exitCode = %v, want 42", result.ExitCode)
	}
}

func TestExecuteTimeout(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 100 * time.Millisecond,
	})

	result, err := exec.Execute(context.Background(), "sleep 10")
	if err == nil {
		t.Fatal("Execute() expected timeout error, got nil")
	}

	if result.ExitCode != -1 {
		t.Errorf("Execute() exitCode = %v, want -1 for timeout", result.ExitCode)
	}

	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("Execute() error = %v, want to contain 'timed out'", err)
	}
}

func TestExecuteStderrCapture(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "echo 'error message' >&2")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stderr, "error message") {
		t.Errorf("Execute() stderr = %v, want to contain 'error message'", result.Stderr)
	}
}

func TestExecuteBothStdoutAndStderr(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "echo 'stdout' && echo 'stderr' >&2")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "stdout") {
		t.Errorf("Execute() stdout = %v, want to contain 'stdout'", result.Stdout)
	}

	if !strings.Contains(result.Stderr, "stderr") {
		t.Errorf("Execute() stderr = %v, want to contain 'stderr'", result.Stderr)
	}
}

func TestExecuteWithEngine(t *testing.T) {
	// Test with a simple echo command as engine
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "echo",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "test command")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	// The engine should prepend "echo" to the command
	if !strings.Contains(result.Stdout, "test command") {
		t.Errorf("Execute() stdout = %v, want to contain 'test command'", result.Stdout)
	}
}

func TestExecuteEmptyCommand(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	// Empty command should still execute (bash -c "")
	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}
}

func TestExecuteCommandWithSpecialCharacters(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "echo 'hello $USER && world'")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("Execute() stdout = %v, want to contain 'hello'", result.Stdout)
	}
}

func TestExecuteCommandWithQuotes(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), `echo "hello world"`)
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "hello world") {
		t.Errorf("Execute() stdout = %v, want to contain 'hello world'", result.Stdout)
	}
}

func TestExecuteCommandWithNewlines(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "echo 'line1'\necho 'line2'")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "line1") {
		t.Errorf("Execute() stdout = %v, want to contain 'line1'", result.Stdout)
	}

	if !strings.Contains(result.Stdout, "line2") {
		t.Errorf("Execute() stdout = %v, want to contain 'line2'", result.Stdout)
	}
}

func TestExecuteContextCancellation(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := exec.Execute(ctx, "sleep 10")
	if err == nil {
		t.Fatal(EXECUTE_ERROR_GOT_NIL)
	}

	if result.ExitCode != -1 {
		t.Errorf("Execute() exitCode = %v, want -1 for context cancellation", result.ExitCode)
	}
}

func TestExecuteCommandNotFound(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "nonexistentcommand12345")
	if err == nil {
		t.Fatal(EXECUTE_ERROR_GOT_NIL)
	}

	// Command not found typically returns exit code 127
	if result.ExitCode != 127 {
		t.Errorf("Execute() exitCode = %v, want 127 for command not found", result.ExitCode)
	}

	if !strings.Contains(result.Stderr, "not found") && !strings.Contains(result.Stderr, "command") {
		t.Errorf("Execute() stderr = %v, want to contain 'not found' or 'command'", result.Stderr)
	}
}

func TestQuoteCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "simple command",
			command:  "ls -la",
			expected: "\"ls -la\"",
		},
		{
			name:     "command with double quotes",
			command:  `echo "hello"`,
			expected: `echo "hello"`,
		},
		{
			name:     "command with single quotes",
			command:  `echo 'hello'`,
			expected: `echo 'hello'`,
		},
		{
			name:     "command with escaped quotes",
			command:  `echo "hello \"world\""`,
			expected: `echo "hello \"world\""`,
		},
		{
			name:     "command with special characters",
			command:  `echo $HOME`,
			expected: `"echo \$HOME"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := quoteCommand(tt.command)
			if result != tt.expected {
				t.Errorf("quoteCommand() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExecuteMultiLineCommand(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "echo 'first'\necho 'second'\necho 'third'")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "first") {
		t.Errorf("Execute() stdout = %v, want to contain 'first'", result.Stdout)
	}

	if !strings.Contains(result.Stdout, "second") {
		t.Errorf("Execute() stdout = %v, want to contain 'second'", result.Stdout)
	}

	if !strings.Contains(result.Stdout, "third") {
		t.Errorf("Execute() stdout = %v, want to contain 'third'", result.Stdout)
	}
}

func TestExecutePipeCommand(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "echo 'hello world' | grep hello")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("Execute() stdout = %v, want to contain 'hello'", result.Stdout)
	}
}

func TestExecuteRedirection(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "echo 'test' > /tmp/test_executor.txt && cat /tmp/test_executor.txt && rm /tmp/test_executor.txt")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "test") {
		t.Errorf("Execute() stdout = %v, want to contain 'test'", result.Stdout)
	}
}

func TestExecuteComplexEngine(t *testing.T) {
	// Test with a more complex engine prefix
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "echo 'engine:'",
		Timeout: 30 * time.Second,
	})

	result, err := exec.Execute(context.Background(), "test")
	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	// The engine should prepend "echo 'engine:'" to the command
	if !strings.Contains(result.Stdout, "engine:") {
		t.Errorf("Execute() stdout = %v, want to contain 'engine:'", result.Stdout)
	}
}

func TestExecuteDurationMeasurement(t *testing.T) {
	exec := NewExecutor(&config.ExecutionConfig{
		Engine:  "",
		Timeout: 30 * time.Second,
	})

	start := time.Now()
	result, err := exec.Execute(context.Background(), "sleep 0.1")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf(EXECUTE_ERROR, err)
	}

	if result.ExitCode != 0 {
		t.Errorf(EXECUTE_EXIT_CODE, result.ExitCode)
	}

	// Duration should be approximately 100ms
	if result.Duration < 100*time.Millisecond {
		t.Errorf("Execute() duration = %v, want >= 100ms", result.Duration)
	}

	// The measured duration should be close to the actual elapsed time
	if result.Duration > elapsed+50*time.Millisecond {
		t.Errorf("Execute() duration = %v, want close to %v", result.Duration, elapsed)
	}
}
