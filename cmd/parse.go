package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cli_agent/internal/parser"

	"github.com/spf13/pflag"
)

// ExecuteParse runs the parse command
func ExecuteParse(args []string) error {
	flagSet := pflag.NewFlagSet("parse", pflag.ExitOnError)
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

	// Extract bash command from input
	command, err := parser.ExtractBashCommand(input)
	if err != nil {
		switch {
		case err == parser.ErrNoBashAction:
			return fmt.Errorf("no bash action found in input (expected <do>...</do> tags)")
		case err == parser.ErrMultipleBashActions:
			return fmt.Errorf("multiple bash actions found in input (only one <do>...</do> block allowed)")
		case err == parser.ErrEmptyBashAction:
			return fmt.Errorf("empty bash action found in <do>...</do> tags")
		default:
			return fmt.Errorf("error parsing input: %w", err)
		}
	}

	// Output the extracted command
	fmt.Println(command)

	return nil
}