package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// Test constants
const (
	testBaseURL     = "https://api.openai.com/v1"
	testAPIKey      = "sk-test-key"
	testConfigName  = "Test Config"
	testVersion     = "1.0.0"
	testName        = "Test"
	testModel       = "glm-4.7"
	testTemperature = 0.7
	testMaxTokens   = 128000
	testTopP        = 1.0
	testTimeout     = 120
	testMaxRetries  = 3
	testRetryDelay  = 1000
	testAltModel    = "gpt-3.5-turbo"
	testAltTimeout  = 60
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

	if cfg.GetName() != testConfigName {
		t.Errorf(errExpectedName, testConfigName, cfg.GetName())
	}
	if cfg.GetVersion() != testVersion {
		t.Errorf(errExpectedVersion, testVersion, cfg.GetVersion())
	}
	if cfg.GetOpenAIConfig().BaseURL != testBaseURL {
		t.Errorf(errExpectedBaseURL, testBaseURL, cfg.GetOpenAIConfig().BaseURL)
	}
	if cfg.GetOpenAIConfig().APIKey != testAPIKey {
		t.Errorf(errExpectedAPIKey, testAPIKey, cfg.GetOpenAIConfig().APIKey)
	}
}

func TestLoadValidJSON(t *testing.T) {
	cfg, err := Load("../../testdata/config/valid.json")
	if err != nil {
		t.Fatalf("Failed to load valid JSON config: %v", err)
	}

	if cfg.GetName() != testConfigName {
		t.Errorf(errExpectedName, testConfigName, cfg.GetName())
	}
	if cfg.GetOpenAIConfig().BaseURL != testBaseURL {
		t.Errorf(errExpectedBaseURL, testBaseURL, cfg.GetOpenAIConfig().BaseURL)
	}
}

func TestLoadValidTOML(t *testing.T) {
	cfg, err := Load("../../testdata/config/valid.toml")
	if err != nil {
		t.Fatalf("Failed to load valid TOML config: %v", err)
	}

	if cfg.GetName() != testConfigName {
		t.Errorf(errExpectedName, testConfigName, cfg.GetName())
	}
	if cfg.GetOpenAIConfig().BaseURL != testBaseURL {
		t.Errorf(errExpectedBaseURL, testBaseURL, cfg.GetOpenAIConfig().BaseURL)
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
		Name:    testName,
		Version: testVersion,
		OpenAI: OpenAIConfig{
			BaseURL: testBaseURL,
			APIKey:  testAPIKey,
		},
	}

	err := ValidateAndSetDefaults(cfg)
	if err != nil {
		t.Errorf("Expected valid config to pass validation, got error: %v", err)
	}
}

func TestValidateMissingBaseURL(t *testing.T) {
	cfg := &Config{
		Name:    testName,
		Version: testVersion,
		OpenAI: OpenAIConfig{
			APIKey: testAPIKey,
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
		Name:    testName,
		Version: testVersion,
		OpenAI: OpenAIConfig{
			BaseURL: testBaseURL,
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
		Name:    testName,
		Version: testVersion,
		OpenAI: OpenAIConfig{
			BaseURL: testBaseURL,
			APIKey:  testAPIKey,
		},
	}

	err := ValidateAndSetDefaults(cfg)
	if err != nil {
		t.Fatalf(errFailedToValidate, err)
	}

	// Check HTTP client defaults
	if cfg.GetOpenAIConfig().HTTPClient.Timeout != testTimeout {
		t.Errorf("Expected default timeout %d, got %d", testTimeout, cfg.GetOpenAIConfig().HTTPClient.Timeout)
	}
	if cfg.GetOpenAIConfig().HTTPClient.MaxRetries != testMaxRetries {
		t.Errorf("Expected default max_retries %d, got %d", testMaxRetries, cfg.GetOpenAIConfig().HTTPClient.MaxRetries)
	}
	if cfg.GetOpenAIConfig().HTTPClient.RetryDelay != testRetryDelay {
		t.Errorf("Expected default retry_delay %d, got %d", testRetryDelay, cfg.GetOpenAIConfig().HTTPClient.RetryDelay)
	}

	// Check defaults
	if cfg.GetOpenAIConfig().Defaults.Model != testModel {
		t.Errorf("Expected default model '%s', got '%s'", testModel, cfg.GetOpenAIConfig().Defaults.Model)
	}
	if cfg.GetOpenAIConfig().Defaults.Temperature != testTemperature {
		t.Errorf("Expected default temperature %f, got %f", testTemperature, cfg.GetOpenAIConfig().Defaults.Temperature)
	}
	if cfg.GetOpenAIConfig().Defaults.MaxTokens != testMaxTokens {
		t.Errorf("Expected default max_tokens %d, got %d", testMaxTokens, cfg.GetOpenAIConfig().Defaults.MaxTokens)
	}
	if cfg.GetOpenAIConfig().Defaults.TopP != testTopP {
		t.Errorf("Expected default top_p %f, got %f", testTopP, cfg.GetOpenAIConfig().Defaults.TopP)
	}
}

func TestValidatePartialConfig(t *testing.T) {
	cfg := &Config{
		Name:    testName,
		Version: testVersion,
		OpenAI: OpenAIConfig{
			BaseURL: testBaseURL,
			APIKey:  testAPIKey,
			HTTPClient: &HTTPClient{
				Timeout: testAltTimeout, // Override default
			},
			Defaults: &Defaults{
				Model: testAltModel, // Override default
			},
		},
	}

	err := ValidateAndSetDefaults(cfg)
	if err != nil {
		t.Fatalf("Failed to validate partial config: %v", err)
	}

	// Check that overridden values are preserved
	if cfg.GetOpenAIConfig().HTTPClient.Timeout != testAltTimeout {
		t.Errorf("Expected timeout %d, got %d", testAltTimeout, cfg.GetOpenAIConfig().HTTPClient.Timeout)
	}
	if cfg.GetOpenAIConfig().Defaults.Model != testAltModel {
		t.Errorf("Expected model '%s', got '%s'", testAltModel, cfg.GetOpenAIConfig().Defaults.Model)
	}

	// Check that defaults are still set for non-overridden values
	if cfg.GetOpenAIConfig().HTTPClient.MaxRetries != 0 {
		t.Errorf("Expected max_retries %d, got %d", 0, cfg.GetOpenAIConfig().HTTPClient.MaxRetries)
	}
	if cfg.GetOpenAIConfig().Defaults.Temperature != 0 {
		t.Errorf("Expected temperature %f, got %f", 0.0, cfg.GetOpenAIConfig().Defaults.Temperature)
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
