package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"cli_agent/internal/config"
	"cli_agent/internal/executor"

	"github.com/spf13/pflag"
)

// ExecuteExecute runs the execute command
func ExecuteExecute(args []string) error {
	flagSet := pflag.NewFlagSet("execute", pflag.ExitOnError)
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	var input string

	// Check if input is provided as a command-line argument
	if flagSet.NArg() > 0 {
		input = flagSet.Arg(0)
	} else {
		// Read from stdin if no argument provided
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		input = strings.Join(lines, "\n")
	}

	// Load configuration to get execution settings
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate the configuration
	if err := config.ValidateAndSetDefaults(cfg); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Get execution config
	execConfig := cfg.GetExecutionConfig()

	// Create executor
	exec := executor.NewExecutor(&execConfig)

	// Execute the command
	ctx := context.Background()
	result, err := exec.Execute(ctx, input)
	if err != nil {
		// Print result even if there was an error (e.g., non-zero exit code)
		if result != nil {
			if result.Stdout != "" {
				fmt.Println(result.Stdout)
			}
			if result.Stderr != "" {
				fmt.Fprintln(os.Stderr, result.Stderr)
			}
		}
		return fmt.Errorf("command execution failed: %w", err)
	}

	// Print stdout
	if result.Stdout != "" {
		fmt.Println(result.Stdout)
	}

	// Print stderr if present
	if result.Stderr != "" {
		fmt.Fprintln(os.Stderr, result.Stderr)
	}

	// Print execution info
	fmt.Fprintf(os.Stderr, "Exit code: %d\n", result.ExitCode)
	fmt.Fprintf(os.Stderr, "Duration: %v\n", result.Duration)

	return nil
}