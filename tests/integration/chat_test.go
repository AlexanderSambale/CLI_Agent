//go:build integration
// +build integration

package integration

import (
	"cli_agent/tests"
	"strings"
	"testing"
)

// TestChat tests the chat command
func TestChat(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	root, cliAgent := tests.GetRootAndCLIAgent(t)

	// This test requires a valid configuration file with API credentials
	t.Run("Run chat with valid config", func(t *testing.T) {
		configPath := tests.ConfigPathIfExisting(t, root, "config.yaml")
		stdout, stderr, exitCode := tests.RunCLICommand(t, cliAgent, "-c", configPath, "chat", "What is the capital of France?")

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
			t.Errorf("Stderr: %s", stderr)
		}

		if !strings.Contains(string(stdout), "Paris") {
			t.Errorf("Expected Paris in the output: %s", stdout)
		}
	})
	t.Run("Run chat with invalid config", func(t *testing.T) {
		configPath := tests.ConfigPathIfExisting(t, root, "testdata/config/invalid.yaml")
		_, stderr, exitCode := tests.RunCLICommand(t, cliAgent, "-c", configPath, "chat", "What is the capital of France?")

		// Check exit code
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
			t.Errorf("Stderr: %s", stderr)
		}
	})
}
