//go:build integration
// +build integration

package integration

import (
	"cli_agent/tests"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAgent tests the agent command
func TestAgent(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	root, cliAgent := tests.GetRootAndCLIAgent(t)

	path := "/v1/chat/completions"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`{
				"id": "chatcmpl-test",
				"object": "chat.completion",
				"created": 1710000000,
				"model": "gpt-5.3",
				"choices": [
					{
						"index": 0,
						"message": {
							"role": "assistant",
							"content": "<do>ls -l README.md</do>"
						},
						"finish_reason": "stop"
					}
				],
				"usage": {
					"prompt_tokens": 0,
					"completion_tokens": 0,
					"total_tokens": 0
				}
			}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	baseURL := server.URL + "/v1"

	// This test requires a valid configuration file with API credentials
	t.Run("Run agent with file or piped input", func(t *testing.T) {
		configPath := tests.ConfigPathIfExisting(t, root, "testdata/config/valid.yaml")
		testFilePath := filepath.Join(root, "testdata/files/test.txt")

		// Check if test file exists
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			t.Fatalf("Test file not found: %s", testFilePath)
		}

		// Run agent with file input by reading the file content and passing it as stdin
		// This test simulates: ./cli-agent -c config.yaml agent < test.txt
		content, err := os.ReadFile(testFilePath)
		if err != nil {
			t.Fatalf("Failed to read test file: %v", err)
		}
		stdout, stderr, exitCode := tests.RunCLICommandWithStdin(t, cliAgent, string(content), "-c", configPath, "--base-url", baseURL, "agent")

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
			t.Errorf("Stderr: %s", stderr)
		}

		// Check that the output contains expected execution markers
		output := string(stdout)

		// Check that a command was executed (should contain ls or similar)
		if !strings.Contains(output, "ls") && !strings.Contains(output, "README.md") {
			t.Errorf("Expected 'ls' and 'README' execution output in: %s", output)
		}
	})
}
