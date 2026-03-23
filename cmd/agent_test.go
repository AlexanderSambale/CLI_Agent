package cmd

import (
	"cli_agent/internal/config"
	"cli_agent/internal/logger"
	mock_openai "cli_agent/internal/mocks"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

// TestLoadAgentConfig_Success tests successful loading of agent configuration
func TestLoadAgentConfig_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
			System:      "You are a helpful coding assistant.",
		},
		Execution: config.ExecutionConfig{
			Engine:  "",
			Timeout: 30 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 10,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with command-line argument
	args := []string{"test prompt"}

	cfg, err := loadAgentConfig(mockClient, args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify configuration
	if cfg.model != "glm-4.7" {
		t.Errorf("Expected model 'glm-4.7', got '%s'", cfg.model)
	}
	if cfg.temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", cfg.temperature)
	}
	if cfg.maxTokens != 128000 {
		t.Errorf("Expected maxTokens 128000, got %d", cfg.maxTokens)
	}
	if cfg.topP != 1.0 {
		t.Errorf("Expected topP 1.0, got %f", cfg.topP)
	}
	if cfg.systemMessage != "You are a helpful coding assistant." {
		t.Errorf("Expected system message 'You are a helpful coding assistant.', got '%s'", cfg.systemMessage)
	}
	if cfg.maxTurnsLimit != 10 {
		t.Errorf("Expected maxTurnsLimit 10, got %d", cfg.maxTurnsLimit)
	}

	// Verify messages count (system + user)
	if len(cfg.messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(cfg.messages))
	}

	// Verify executor is created
	if cfg.executor == nil {
		t.Error("Expected executor to be created")
	}

	// Verify logger is set
	if cfg.logger == nil {
		t.Error("Expected logger to be set")
	}
}

// TestLoadAgentConfig_WithFlags tests loading agent configuration with command-line flags
func TestLoadAgentConfig_WithFlags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with defaults
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
			System:      "Default system message",
		},
		Execution: config.ExecutionConfig{
			Engine:  "",
			Timeout: 30 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 10,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with command-line flags overriding defaults
	args := []string{
		"--model", "gpt-4",
		"--temperature", "0.5",
		"--max-tokens", "4096",
		"--top-p", "0.9",
		"--system", "Custom system message",
		"--max-turns", "5",
		"--engine", "docker run --rm alpine bash -c",
		"--timeout", "60",
		"test prompt",
	}

	cfg, err := loadAgentConfig(mockClient, args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify flags override config defaults
	if cfg.model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", cfg.model)
	}
	if cfg.temperature != 0.5 {
		t.Errorf("Expected temperature 0.5, got %f", cfg.temperature)
	}
	if cfg.maxTokens != 4096 {
		t.Errorf("Expected maxTokens 4096, got %d", cfg.maxTokens)
	}
	if cfg.topP != 0.9 {
		t.Errorf("Expected topP 0.9, got %f", cfg.topP)
	}
	if cfg.systemMessage != "Custom system message" {
		t.Errorf("Expected system message 'Custom system message', got '%s'", cfg.systemMessage)
	}
	if cfg.maxTurnsLimit != 5 {
		t.Errorf("Expected maxTurnsLimit 5, got %d", cfg.maxTurnsLimit)
	}

	// Verify executor is created with overridden engine and timeout
	if cfg.executor == nil {
		t.Error("Expected executor to be created")
	}
}

// TestLoadAgentConfig_NoInput tests error when no input is provided
func TestLoadAgentConfig_NoInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
		},
		Execution: config.ExecutionConfig{
			Engine:  "",
			Timeout: 30 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 10,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with no input (no arguments and stdin is a terminal)
	args := []string{}

	_, err := loadAgentConfig(mockClient, args)
	if err == nil {
		t.Error("Expected error for no input, got nil")
	}
	if err.Error() != "no input provided" {
		t.Errorf("Expected error 'no input provided', got '%v'", err)
	}
}

// TestLoadAgentConfig_NoSystemMessage tests configuration without system message
func TestLoadAgentConfig_NoSystemMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config without system message
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
			System:      "", // No system message
		},
		Execution: config.ExecutionConfig{
			Engine:  "",
			Timeout: 30 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 10,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with command-line argument
	args := []string{"test prompt"}

	cfg, err := loadAgentConfig(mockClient, args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify only user message is present (no system message)
	if len(cfg.messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(cfg.messages))
	}

	// Verify system message is empty
	if cfg.systemMessage != "" {
		t.Errorf("Expected empty system message, got '%s'", cfg.systemMessage)
	}
}

// TestLoadAgentConfig_WithSystemMessageFlag tests configuration with system message flag
func TestLoadAgentConfig_WithSystemMessageFlag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with default system message
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
			System:      "Default system message",
		},
		Execution: config.ExecutionConfig{
			Engine:  "",
			Timeout: 30 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 10,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with system message flag
	args := []string{
		"--system", "Custom system message from flag",
		"test prompt",
	}

	cfg, err := loadAgentConfig(mockClient, args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify system message from flag is used
	if cfg.systemMessage != "Custom system message from flag" {
		t.Errorf("Expected system message 'Custom system message from flag', got '%s'", cfg.systemMessage)
	}

	// Verify both system and user messages are present
	if len(cfg.messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(cfg.messages))
	}
}

// TestLoadAgentConfig_ExecutionConfig tests execution configuration
func TestLoadAgentConfig_ExecutionConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with execution settings
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
		},
		Execution: config.ExecutionConfig{
			Engine:  "docker run --rm alpine bash -c",
			Timeout: 60 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 10,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test without overriding execution config
	args := []string{"test prompt"}

	cfg, err := loadAgentConfig(mockClient, args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify executor is created
	if cfg.executor == nil {
		t.Error("Expected executor to be created")
	}
}

// TestLoadAgentConfig_ExecutionConfigOverride tests execution configuration override with flags
func TestLoadAgentConfig_ExecutionConfigOverride(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with default execution settings
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
		},
		Execution: config.ExecutionConfig{
			Engine:  "",
			Timeout: 30 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 10,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with execution config flags
	args := []string{
		"--engine", "podman run --rm alpine bash -c",
		"--timeout", "120",
		"test prompt",
	}

	cfg, err := loadAgentConfig(mockClient, args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify executor is created
	if cfg.executor == nil {
		t.Error("Expected executor to be created")
	}
}

// TestLoadAgentConfig_AgentConfig tests agent configuration
func TestLoadAgentConfig_AgentConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with agent settings
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
		},
		Execution: config.ExecutionConfig{
			Engine:  "",
			Timeout: 30 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 20,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test without overriding agent config
	args := []string{"test prompt"}

	cfg, err := loadAgentConfig(mockClient, args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify max turns from config is used
	if cfg.maxTurnsLimit != 20 {
		t.Errorf("Expected maxTurnsLimit 20, got %d", cfg.maxTurnsLimit)
	}
}

// TestLoadAgentConfig_AgentConfigOverride tests agent configuration override with flags
func TestLoadAgentConfig_AgentConfigOverride(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with default agent settings
	testConfig := &config.Config{
		Name:    "Test Config",
		Version: "1.0.0",
		OpenAI: config.OpenAIConfig{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  "sk-test-key",
		},
		Model: config.ModelConfig{
			Model:       "glm-4.7",
			Temperature: 0.7,
			MaxTokens:   128000,
			TopP:        1.0,
		},
		Execution: config.ExecutionConfig{
			Engine:  "",
			Timeout: 30 * time.Second,
		},
		Agent: config.AgentConfig{
			MaxTurns: 10,
		},
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with max turns flag
	args := []string{
		"--max-turns", "15",
		"test prompt",
	}

	cfg, err := loadAgentConfig(mockClient, args)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify max turns from flag is used
	if cfg.maxTurnsLimit != 15 {
		t.Errorf("Expected maxTurnsLimit 15, got %d", cfg.maxTurnsLimit)
	}
}
