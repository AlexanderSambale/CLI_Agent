package cmd

import (
	"fmt"
	"os"

	"cli_agent/internal/parser"

	"github.com/spf13/pflag"
)

// ExecuteParse runs the parse command
func ExecuteParse(args []string) error {
	flagSet := pflag.NewFlagSet("parse", pflag.ExitOnError)
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// Get input from command-line argument or stdin
	input, err := readInput(flagSet, os.Stdin)
	if err != nil {
		return err
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
