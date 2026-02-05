package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// General application settings
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	Settings struct {
		Debug   bool `mapstructure:"debug"`
		Verbose bool `mapstructure:"verbose"`
	} `mapstructure:"settings"`

	// OpenAI configuration
	OpenAI OpenAIConfig `mapstructure:"openai"`
}

// OpenAIConfig represents OpenAI-specific configuration
type OpenAIConfig struct {
	// Required fields
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`

	// HTTP client options
	HTTPClient struct {
		Timeout    int `mapstructure:"timeout"`    // in seconds
		MaxRetries int `mapstructure:"max_retries"`
		RetryDelay int `mapstructure:"retry_delay"` // in milliseconds
	} `mapstructure:"http_client"`

	// Default request parameters
	Defaults struct {
		Model       string  `mapstructure:"model"`
		Temperature float64 `mapstructure:"temperature"`
		MaxTokens   int     `mapstructure:"max_tokens"`
		TopP        float64 `mapstructure:"top_p"`
	} `mapstructure:"defaults"`
}

// Load reads and parses the configuration file from the given path
func Load(configPath string) (*Config, error) {
	// Validate that the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
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
func (c *Config) Validate() error {
	// Validate OpenAI configuration
	if c.OpenAI.BaseURL == "" {
		return fmt.Errorf("openai.base_url is required")
	}
	if c.OpenAI.APIKey == "" {
		return fmt.Errorf("openai.api_key is required")
	}

	// Set default values for HTTP client if not specified
	if c.OpenAI.HTTPClient.Timeout == 0 {
		c.OpenAI.HTTPClient.Timeout = 120 // Default 2 minutes
	}
	if c.OpenAI.HTTPClient.MaxRetries == 0 {
		c.OpenAI.HTTPClient.MaxRetries = 3
	}
	if c.OpenAI.HTTPClient.RetryDelay == 0 {
		c.OpenAI.HTTPClient.RetryDelay = 1000 // Default 1 second
	}

	// Set default values for request parameters if not specified
	if c.OpenAI.Defaults.Model == "" {
		c.OpenAI.Defaults.Model = "gpt-4o"
	}
	if c.OpenAI.Defaults.Temperature == 0 {
		c.OpenAI.Defaults.Temperature = 0.7
	}
	if c.OpenAI.Defaults.MaxTokens == 0 {
		c.OpenAI.Defaults.MaxTokens = 2048
	}
	if c.OpenAI.Defaults.TopP == 0 {
		c.OpenAI.Defaults.TopP = 1.0
	}

	return nil
}