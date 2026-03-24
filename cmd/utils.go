package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// readInput reads input from command-line argument or stdin
func readInput(flagSet *pflag.FlagSet, stdin *os.File) (string, error) {
	if flagSet.NArg() > 0 {
		return flagSet.Arg(0), nil
	}

	return readFromStdin(stdin)
}

// readFromStdin reads all input from stdin
// Returns empty string if stdin is a terminal (not a pipe or redirected file)
func readFromStdin(stdin *os.File) (string, error) {
	// Check if stdin is a pipe or redirected file
	stat, _ := stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// stdin is a terminal, not a pipe or file
		return "", nil
	}

	// Read from stdin
	reader := bufio.NewReader(stdin)
	var builder strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Add the last line if it doesn't end with newline
				if line != "" {
					builder.WriteString(line)
				}
				break
			}
			return "", fmt.Errorf("error reading from stdin: %w", err)
		}
		builder.WriteString(line)
	}

	return strings.TrimSpace(builder.String()), nil
}
