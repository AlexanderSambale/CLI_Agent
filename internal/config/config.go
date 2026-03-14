package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// CLIConfig defines the interface for CLI configuration
type CLIConfig interface {
	GetName() string
	GetVersion() string
	GetDebug() bool
	GetVerbose() bool
	GetOpenAIConfig() OpenAIConfig
	SetOpenAIConfig(OpenAIConfig)
	GetExecutionConfig() ExecutionConfig
	SetExecutionConfig(ExecutionConfig)
}

// Config represents the application configuration
type Config struct {
	// General application settings
	Name     string         `mapstructure:"name"`
	Version  string         `mapstructure:"version"`
	Settings SettingsConfig `mapstructure:"settings"`

	// OpenAI configuration
	OpenAI OpenAIConfig `mapstructure:"openai"`

	// Execution configuration
	Execution ExecutionConfig `mapstructure:"execution"`
}

type SettingsConfig struct {
	Debug   bool `mapstructure:"debug"`
	Verbose bool `mapstructure:"verbose"`
}

// GetName returns the application name
func (c *Config) GetName() string {
	return c.Name
}

// GetVersion returns the application version
func (c *Config) GetVersion() string {
	return c.Version
}

// GetDebug returns the debug setting
func (c *Config) GetDebug() bool {
	return c.Settings.Debug
}

// GetVerbose returns the verbose setting
func (c *Config) GetVerbose() bool {
	return c.Settings.Verbose
}

// GetOpenAIConfig returns the OpenAI configuration
func (c *Config) GetOpenAIConfig() OpenAIConfig {
	return c.OpenAI
}

func (c *Config) SetOpenAIConfig(openAIConfig OpenAIConfig) {
	c.OpenAI = openAIConfig
}

// GetExecutionConfig returns the execution configuration
func (c *Config) GetExecutionConfig() ExecutionConfig {
	return c.Execution
}

func (c *Config) SetExecutionConfig(executionConfig ExecutionConfig) {
	c.Execution = executionConfig
}

// OpenAIConfig represents OpenAI-specific configuration
type OpenAIConfig struct {
	// Required fields
	BaseURL    string      `mapstructure:"base_url"`
	APIKey     string      `mapstructure:"api_key"`
	HTTPClient *HTTPClient `mapstructure:"http_client"`
	Defaults   *Defaults   `mapstructure:"defaults"`
}

// HTTP client options
type HTTPClient struct {
	Timeout    int `mapstructure:"timeout"` // in seconds
	MaxRetries int `mapstructure:"max_retries"`
	RetryDelay int `mapstructure:"retry_delay"` // in milliseconds
}

// Default request parameters
type Defaults struct {
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"` // 0.0 – 2.0
	MaxTokens   int     `mapstructure:"max_tokens"`
	TopP        float64 `mapstructure:"top_p"` // 0.0 – 1.0
}

// ExecutionConfig represents command execution configuration
type ExecutionConfig struct {
	Engine  string        `mapstructure:"engine"`  // Command prefix (e.g., "docker run --rm ubuntu bash -c")
	Timeout time.Duration `mapstructure:"timeout"` // timeout in seconds (default: 30)
}

// Load reads and parses the configuration file from the given path
func Load(configPath string) (CLIConfig, error) {
	// Validate that the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s: %w", configPath, err)
	}

	// Get the file extension to determine the format
	ext := filepath.Ext(configPath)

	// Create a new viper instance
	v := viper.New()

	// Set the config file path
	v.SetConfigFile(configPath)

	// Set the config type based on file extension
	switch ext {
	case ".yaml", ".yml":
		v.SetConfigType("yaml")
	case ".json":
		v.SetConfigType("json")
	case ".toml":
		v.SetConfigType("toml")
	default:
		return nil, fmt.Errorf("unsupported configuration file format: %s (supported: .yaml, .yml, .json, .toml)", ext)
	}

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Unmarshal the config into our Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
// and sets Default values if not specified
func ValidateAndSetDefaults(c CLIConfig) error {
	config := c.GetOpenAIConfig()

	// Validate OpenAI configuration
	if config.BaseURL == "" {
		return fmt.Errorf("openai.base_url is required")
	}
	if config.APIKey == "" {
		return fmt.Errorf("openai.api_key is required")
	}

	// Set default values for HTTP client if not specified
	if config.HTTPClient == nil {
		config.HTTPClient = &HTTPClient{Timeout: 120, MaxRetries: 3, RetryDelay: 1000}
	}

	// Set default values for request parameters if not specified
	if config.Defaults == nil {
		config.Defaults = &Defaults{Model: "glm-4.7", Temperature: 0.7, MaxTokens: 128000, TopP: 1.0}
	}

	c.SetOpenAIConfig(config)

	// Set default values for execution config if not specified to max value 64-bit signed integer
	execConfig := c.GetExecutionConfig()
	if execConfig.Timeout == 0 {
		execConfig.Timeout = 1<<63 - 1
	}

	c.SetExecutionConfig(execConfig)

	return nil
}
