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

	input, err := readInput(flagSet)
	if err != nil {
		return err
	}

	cfg, err := loadAndValidateConfig()
	if err != nil {
		return err
	}

	execConfig := cfg.GetExecutionConfig()
	exec := executor.NewExecutor(&execConfig)

	ctx := context.Background()
	result, err := exec.Execute(ctx, input)
	if err != nil {
		printExecutionResult(result)
		return fmt.Errorf("command execution failed: %w", err)
	}

	printExecutionResult(result)
	return nil
}

// readInput reads input from command-line argument or stdin
func readInput(flagSet *pflag.FlagSet) (string, error) {
	if flagSet.NArg() > 0 {
		return flagSet.Arg(0), nil
	}

	return readFromStdin()
}

// readFromStdin reads all input from stdin
func readFromStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading from stdin: %w", err)
	}
	return strings.Join(lines, "\n"), nil
}

// loadAndValidateConfig loads and validates the configuration file
func loadAndValidateConfig() (config.CLIConfig, error) {
	cfg, err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := config.ValidateAndSetDefaults(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// printExecutionResult prints the stdout, stderr, and execution info
func printExecutionResult(result *executor.Result) {
	if result == nil {
		return
	}

	if result.Stdout != "" {
		fmt.Println(result.Stdout)
	}

	if result.Stderr != "" {
		fmt.Fprintln(os.Stderr, result.Stderr)
	}

	fmt.Fprintf(os.Stderr, "Exit code: %d\n", result.ExitCode)
	fmt.Fprintf(os.Stderr, "Duration: %v\n", result.Duration)
}
