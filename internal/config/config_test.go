package config

import (
	tc "cli_agent/testdata/test_constants"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// Error message format constants
const (
	errExpectedName     = "Expected name '%s', got '%s'"
	errExpectedBaseURL  = "Expected base_url '%s', got '%s'"
	errExpectedAPIKey   = "Expected api_key '%s', got '%s'"
	errExpectedVersion  = "Expected version '%s', got '%s'"
	errFailedToLoad     = "Failed to load config: %v"
	errFailedToValidate = "Failed to validate config: %v"
)

func TestLoadValidYAML(t *testing.T) {
	cfg, err := Load("../../testdata/config/valid.yaml")
	if err != nil {
		t.Fatalf("Failed to load valid YAML config: %v", err)
	}

	if cfg.GetName() != tc.TestConfigName {
		t.Errorf(errExpectedName, tc.TestConfigName, cfg.GetName())
	}
	if cfg.GetVersion() != tc.TestVersion {
		t.Errorf(errExpectedVersion, tc.TestVersion, cfg.GetVersion())
	}

	openAIConfig := cfg.GetOpenAIConfig()
	if openAIConfig.BaseURL != tc.TestBaseURL {
		t.Errorf(errExpectedBaseURL, tc.TestBaseURL, openAIConfig.BaseURL)
	}
	if openAIConfig.APIKey != tc.TestAPIKey {
		t.Errorf(errExpectedAPIKey, tc.TestAPIKey, openAIConfig.APIKey)
	}
}

func TestLoadValidJSON(t *testing.T) {
	cfg, err := Load("../../testdata/config/valid.json")
	if err != nil {
		t.Fatalf("Failed to load valid JSON config: %v", err)
	}

	if cfg.GetName() != tc.TestConfigName {
		t.Errorf(errExpectedName, tc.TestConfigName, cfg.GetName())
	}

	openAIConfig := cfg.GetOpenAIConfig()
	if openAIConfig.BaseURL != tc.TestBaseURL {
		t.Errorf(errExpectedBaseURL, tc.TestBaseURL, openAIConfig.BaseURL)
	}
}

func TestLoadValidTOML(t *testing.T) {
	cfg, err := Load("../../testdata/config/valid.toml")
	if err != nil {
		t.Fatalf("Failed to load valid TOML config: %v", err)
	}

	if cfg.GetName() != tc.TestConfigName {
		t.Errorf(errExpectedName, tc.TestConfigName, cfg.GetName())
	}

	openAIConfig := cfg.GetOpenAIConfig()
	if openAIConfig.BaseURL != tc.TestBaseURL {
		t.Errorf(errExpectedBaseURL, tc.TestBaseURL, openAIConfig.BaseURL)
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := Load("../../testdata/config/nonexistent.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected file not found error, got: %v", err)
	}
}

func TestLoadUnsupportedFormat(t *testing.T) {
	// Create a temporary file with unsupported extension
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err := Load(tmpFile)
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
}

func TestValidateValidConfig(t *testing.T) {
	cfg := &Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: OpenAIConfig{
			BaseURL: tc.TestBaseURL,
			APIKey:  tc.TestAPIKey,
		},
	}

	err := ValidateAndSetDefaults(cfg)
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}
}

func TestValidateMissingBaseURL(t *testing.T) {
	cfg := &Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: OpenAIConfig{
			APIKey: tc.TestAPIKey,
		},
	}

	err := ValidateAndSetDefaults(cfg)
	if err == nil {
		t.Error("Expected error for missing base_url, got nil")
	}
	if err.Error() != "openai.base_url is required" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestValidateMissingAPIKey(t *testing.T) {
	cfg := &Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: OpenAIConfig{
			BaseURL: tc.TestBaseURL,
		},
	}

	err := ValidateAndSetDefaults(cfg)
	if err == nil {
		t.Error("Expected error for missing api_key, got nil")
	}
	if err.Error() != "openai.api_key is required" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestValidateDefaultValues(t *testing.T) {
	cfg := &Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: OpenAIConfig{
			BaseURL: tc.TestBaseURL,
			APIKey:  tc.TestAPIKey,
		},
	}

	err := ValidateAndSetDefaults(cfg)
	if err != nil {
		t.Fatalf(errFailedToValidate, err)
	}

	openAIConfig := cfg.GetOpenAIConfig()
	// Check HTTP client defaults
	if openAIConfig.HTTPClient.Timeout != tc.TestTimeout {
		t.Errorf("Expected default timeout %d, got %d", tc.TestTimeout, openAIConfig.HTTPClient.Timeout)
	}
	if openAIConfig.HTTPClient.MaxRetries != tc.TestMaxRetries {
		t.Errorf("Expected default max_retries %d, got %d", tc.TestMaxRetries, openAIConfig.HTTPClient.MaxRetries)
	}
	if openAIConfig.HTTPClient.RetryDelay != tc.TestRetryDelay {
		t.Errorf("Expected default retry_delay %d, got %d", tc.TestRetryDelay, openAIConfig.HTTPClient.RetryDelay)
	}

	// Check defaults
	if openAIConfig.Defaults.Model != tc.TestModel {
		t.Errorf("Expected default model '%s', got '%s'", tc.TestModel, openAIConfig.Defaults.Model)
	}
	if openAIConfig.Defaults.Temperature != tc.TestTemperature {
		t.Errorf("Expected default temperature %f, got %f", tc.TestTemperature, openAIConfig.Defaults.Temperature)
	}
	if openAIConfig.Defaults.MaxTokens != tc.TestMaxTokens {
		t.Errorf("Expected default max_tokens %d, got %d", tc.TestMaxTokens, openAIConfig.Defaults.MaxTokens)
	}
	if openAIConfig.Defaults.TopP != tc.TestTopP {
		t.Errorf("Expected default top_p %f, got %f", tc.TestTopP, openAIConfig.Defaults.TopP)
	}
}

func TestValidatePartialConfig(t *testing.T) {
	cfg := &Config{
		Name:    tc.TestName,
		Version: tc.TestVersion,
		OpenAI: OpenAIConfig{
			BaseURL: tc.TestBaseURL,
			APIKey:  tc.TestAPIKey,
			HTTPClient: &HTTPClient{
				Timeout: tc.TestAltTimeout, // Override default
			},
			Defaults: &Defaults{
				Model: tc.TestAltModel, // Override default
			},
		},
	}

	err := ValidateAndSetDefaults(cfg)
	if err != nil {
		t.Fatalf("Failed to validate partial config: %v", err)
	}

	openAIConfig := cfg.GetOpenAIConfig()
	// Check that overridden values are preserved
	if openAIConfig.HTTPClient.Timeout != tc.TestAltTimeout {
		t.Errorf("Expected timeout %d, got %d", tc.TestAltTimeout, openAIConfig.HTTPClient.Timeout)
	}
	if openAIConfig.Defaults.Model != tc.TestAltModel {
		t.Errorf("Expected model '%s', got '%s'", tc.TestAltModel, openAIConfig.Defaults.Model)
	}

	// Check that defaults are still set for non-overridden values
	if openAIConfig.HTTPClient.MaxRetries != 0 {
		t.Errorf("Expected max_retries %d, got %d", 0, openAIConfig.HTTPClient.MaxRetries)
	}
	if openAIConfig.Defaults.Temperature != 0 {
		t.Errorf("Expected temperature %f, got %f", 0.0, openAIConfig.Defaults.Temperature)
	}
}

func TestLoadAndValidateValidConfig(t *testing.T) {
	cfg, err := Load("../../testdata/config/valid.yaml")
	if err != nil {
		t.Fatalf(errFailedToLoad, err)
	}

	err = ValidateAndSetDefaults(cfg)
	if err != nil {
		t.Errorf("Failed to validate loaded config: %v", err)
	}
}

func TestLoadAndValidateInvalidConfig(t *testing.T) {
	cfg, err := Load("../../testdata/config/invalid.yaml")
	if err != nil {
		t.Fatalf(errFailedToLoad, err)
	}

	err = ValidateAndSetDefaults(cfg)
	if err == nil {
		t.Error("Expected validation error for invalid config, got nil")
	}
}
