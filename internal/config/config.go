package config

import (
	defaults "cli_agent/internal/constants"
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
	GetModelConfig() ModelConfig
	SetModelConfig(ModelConfig)
	GetExecutionConfig() ExecutionConfig
	SetExecutionConfig(ExecutionConfig)
	GetAgentConfig() AgentConfig
	SetAgentConfig(AgentConfig)
}

// Config represents the application configuration
type Config struct {
	// General application settings
	Name     string         `mapstructure:"name"`
	Version  string         `mapstructure:"version"`
	Settings SettingsConfig `mapstructure:"settings"`

	// OpenAI configuration
	OpenAI OpenAIConfig `mapstructure:"openai"`

	// Model configuration
	Model ModelConfig `mapstructure:"model"`

	// Execution configuration
	Execution ExecutionConfig `mapstructure:"execution"`

	// Agent configuration
	Agent AgentConfig `mapstructure:"agent"`
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

// GetModelConfig returns the model configuration
func (c *Config) GetModelConfig() ModelConfig {
	return c.Model
}

func (c *Config) SetModelConfig(modelConfig ModelConfig) {
	c.Model = modelConfig
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
}

// HTTP client options
type HTTPClient struct {
	Timeout    int `mapstructure:"timeout"` // in seconds
	MaxRetries int `mapstructure:"max_retries"`
	RetryDelay int `mapstructure:"retry_delay"` // in milliseconds
}

// ModelConfig represents model-specific configuration
type ModelConfig struct {
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"` // (0.0,2.0]
	MaxTokens   int     `mapstructure:"max_tokens"`
	TopP        float64 `mapstructure:"top_p"`  // (0.0,1.0]
	System      string  `mapstructure:"system"` // System message for chat/agent context
}

// ExecutionConfig represents command execution configuration
type ExecutionConfig struct {
	Engine  string        `mapstructure:"engine"`  // Command prefix (e.g., "docker run --rm ubuntu bash -c")
	Timeout time.Duration `mapstructure:"timeout"` // timeout in seconds (default: 30)
}

// AgentConfig represents agent-specific configuration
type AgentConfig struct {
	MaxTurns int `mapstructure:"max_turns"` // Maximum number of agent turns (default: 10)
}

// GetAgentConfig returns the agent configuration
func (c *Config) GetAgentConfig() AgentConfig {
	return c.Agent
}

// SetAgentConfig sets the agent configuration
func (c *Config) SetAgentConfig(agentConfig AgentConfig) {
	c.Agent = agentConfig
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
		config.HTTPClient = &HTTPClient{Timeout: defaults.HttpTimeout, MaxRetries: defaults.MaxRetries, RetryDelay: defaults.RetryDelay}
	}

	c.SetOpenAIConfig(config)

	// Set default values for model config if not specified
	modelConfig := c.GetModelConfig()

	if modelConfig.Model == "" {
		modelConfig.Model = defaults.Model
	}
	if modelConfig.Temperature <= 0.0 || modelConfig.Temperature > 2.0 {
		modelConfig.Temperature = defaults.Temperature
	}
	if modelConfig.MaxTokens < 1 {
		modelConfig.MaxTokens = defaults.MaxTokens
	}
	if modelConfig.TopP <= 0.0 || modelConfig.TopP > 1.0 {
		modelConfig.TopP = defaults.TopP
	}

	c.SetModelConfig(modelConfig)

	// Set default values for execution config if not specified to max value 64-bit signed integer
	execConfig := c.GetExecutionConfig()
	if execConfig.Timeout == 0 {
		execConfig.Timeout = defaults.ExecTimeout
	}

	c.SetExecutionConfig(execConfig)

	// Set default values for agent config if not specified
	agentConfig := c.GetAgentConfig()
	if agentConfig.MaxTurns == 0 {
		agentConfig.MaxTurns = defaults.MaxTurns
	}

	c.SetAgentConfig(agentConfig)

	return nil
}
