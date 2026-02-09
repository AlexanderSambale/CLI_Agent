package parser

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrNoBashAction        = errors.New("No bash action found")
	ErrMultipleBashActions = errors.New("Multiple bash actions found in input")
	ErrEmptyBashAction     = errors.New("Empty bash action")
)

// ExtractBashCommand extracts a bash command from a markdown code block.
// It searches for the first occurrence of a code block marked with the "bash" language identifier.
// Returns the extracted command with whitespace trimmed, or an error if no valid bash block is found.
//
// Parameters:
//   - input: string containing markdown with ```bash...``` code blocks
//
// Returns:
//   - string: the extracted bash command
//   - error: ErrNoBashAction if no bash block found, ErrMultipleBashCommands if multiple blocks found
func ExtractBashCommand(input string) (string, error) {
	// Pattern to match <do>...</do> code blocks
	// (easier to parse than nested ```bash ...``` which could be in markdown files if you want to create a markdown file from bash)
	// Matches: <do> until </do> newlines are explicitly allowed!
	pattern := "(?s)<do>(.*?)</do>"
	re := regexp.MustCompile(pattern)

	// Find all matches to check for multiple bash blocks
	matches := re.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return "", ErrNoBashAction
	}

	// Check for multiple bash code blocks
	if len(matches) > 1 {
		return "", ErrMultipleBashActions
	}

	// Extract the command from the first (and only) match
	untrimmedCommand := matches[0][1]
	command := strings.TrimSpace(untrimmedCommand)

	if command == "" {
		return "", ErrEmptyBashAction
	}

	return command, nil
}
