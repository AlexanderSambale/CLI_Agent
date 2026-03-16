package openai

import (
	"cli_agent/internal/config"
	"cli_agent/internal/logger"
	tc "cli_agent/testdata/test_constants"
	"reflect"
	"testing"
)

const (
	errExpectedClientNonNil = "Expected client to be non-nil"
)

// Helper function to create a test config
func createTestConfig() *config.Config {
	return &config.Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: config.OpenAIConfig{
			BaseURL: tc.TestBaseURL,
			APIKey:  tc.TestAPIKey,
			HTTPClient: &config.HTTPClient{
				Timeout:    120,
				MaxRetries: 3,
				RetryDelay: 1000,
			},
		},
		Model: config.ModelConfig{
			Model:       "gpt-4",
			Temperature: 0.7,
			MaxTokens:   2048,
			TopP:        1.0,
			System:      "",
		},
	}
}

func TestNewClientValidConfig(t *testing.T) {
	cfg := createTestConfig()

	log := logger.NewLogger(false, false)

	client, err := NewClient(cfg, log)
	if err != nil {
		t.Fatalf(tc.ErrFailedToCreateClient, err)
	}

	if client == nil {
		t.Error(errExpectedClientNonNil)
	}

	if !reflect.DeepEqual(client.GetCLIConfig(), cfg) {
		t.Errorf("Expected client config %+v to match input config %+v", client.GetCLIConfig(), cfg)
	}

	// Compare Logger attributes
	if !reflect.DeepEqual(client.GetLogger(), log) {
		t.Errorf("Expected client logger %+v to match input logger %+v", client.GetLogger(), log)
	}
}

func TestNewClientMissingBaseURL(t *testing.T) {
	cfg := &config.Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: config.OpenAIConfig{
			APIKey: "sk-test-key",
		},
	}

	log := logger.NewLogger(false, false)

	_, err := NewClient(cfg, log)
	if err == nil {
		t.Error("Expected error for missing base_url, got nil")
	}

	if err.Error() != BaseURLIsRequired {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestNewClientMissingAPIKey(t *testing.T) {
	cfg := &config.Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: config.OpenAIConfig{
			BaseURL: tc.TestBaseURL,
		},
	}

	log := logger.NewLogger(false, false)

	_, err := NewClient(cfg, log)
	if err == nil {
		t.Error("Expected error for missing api_key, got nil")
	}

	if err.Error() != ApiKeyIsRequired {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestNewClientWithDefaults(t *testing.T) {
	cfg := &config.Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: config.OpenAIConfig{
			BaseURL: tc.TestBaseURL,
			APIKey:  tc.TestAPIKey,
			// HTTPClient and Defaults are empty, should work
		},
	}

	config.ValidateAndSetDefaults(cfg)

	log := logger.NewLogger(false, false)

	client, err := NewClient(cfg, log)
	if err != nil {
		t.Fatalf("Failed to create client with defaults: %v", err)
	}

	if client == nil {
		t.Error(errExpectedClientNonNil)
	}
}
