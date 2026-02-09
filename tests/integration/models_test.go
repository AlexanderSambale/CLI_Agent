//go:build integration
// +build integration

package integration

import (
	"cli_agent/tests"
	"strings"
	"testing"
)

// This test requires a valid configuration file with API credentials
func TestIntegrationModels(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	root, cliAgent := tests.GetRootAndCLIAgent(t)

	configPath := tests.ConfigPathIfExisting(t, root, "config.yaml")
	stdout, stderr, exitCode := tests.RunCLICommand(t, cliAgent, "-c", configPath, "models", "--list")

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
		t.Errorf("Stderr: %s", stderr)
	}

	if !strings.Contains(string(stdout), "Owned By") {
		t.Errorf("Expected 'Owned By' in the output: %s", stdout)
	}

	modelID, _, _ := tests.GetFirstModel(stdout)
	stdout, stderr, exitCode = tests.RunCLICommand(t, cliAgent, "-c", configPath, "models", "--get", modelID)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
		t.Errorf("Stderr: %s", stderr)
	}

	if !strings.Contains(string(stdout), "ID:") {
		t.Errorf("Expected 'ID:' in the output: %s", stdout)
	}
}
