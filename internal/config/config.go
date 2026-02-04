package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Add configuration fields here as needed
	// Example:
	// Name     string `mapstructure:"name"`
	// Version  string `mapstructure:"version"`
	// Settings struct {
	//     Debug bool `mapstructure:"debug"`
	// } `mapstructure:"settings"`
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
	// Add validation logic here as needed
	// Example:
	// if c.Name == "" {
	//     return fmt.Errorf("name is required")
	// }
	return nil
}