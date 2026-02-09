package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func ConfigPathIfExisting(t *testing.T, root string, filename string) string {
	configPath := filepath.Join(root, filename)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file not found: %s", configPath)
	}
	return configPath
}

// RunCLICommand executes the CLI binary and returns output and exit code
func RunCLICommand(t *testing.T, command string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(command, args...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("Command failed with unexpected error: %v", err)
		}
	}

	return stdout, stderr, exitCode
}

func GetRootAndCLIAgent(t *testing.T) (string, string) {
	root := os.Getenv("PROJECT_ROOT")
	t.Log(root)
	cwd, _ := os.Getwd()
	t.Log(cwd)

	cliAgent := filepath.Join(root, "cli-agent")
	return root, cliAgent
}

func GetFirstModel(output string) (string, string, string) {

	lines := strings.Split(output, "\n")

	index := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "ID") {
			index = i
			break
		}
	}

	lineWithFirstModel := lines[index+1]

	re := regexp.MustCompile(`^([^\s]+)\s+([^\s]+)\s+(\d+)$`)
	matches := re.FindStringSubmatch(lineWithFirstModel)

	id := matches[1]
	ownedBy := matches[2]
	created := matches[3]

	return id, ownedBy, created
}
