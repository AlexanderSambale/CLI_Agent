package cmd

import (
	"cli_agent/internal/config"
	"cli_agent/internal/logger"
	mock_openai "cli_agent/internal/mocks"
	"os"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

// Test constants
const (
	testPrompt                     = "test prompt"
	errExpectedNoError             = "Expected no error, got: %v"
	errFailedToLoadTestConfig      = "Failed to load test config: %v"
	errExpectedExecutorToBeCreated = "Expected executor to be created"
	TestValidYAMLConfig            = "../testdata/config/valid.yaml"
)

// TestLoadAgentConfig_Success tests successful loading of agent configuration
func TestLoadAgentConfigSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Load test config from file
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with command-line argument
	args := []string{testPrompt}

	cfg, err := loadAgentConfig(mockClient, args, os.Stdin)
	if err != nil {
		t.Fatalf(errExpectedNoError, err)
	}

	// Verify configuration
	if cfg.model != testConfig.GetModelConfig().Model {
		t.Errorf("Expected model '%s', got '%s'", testConfig.GetModelConfig().Model, cfg.model)
	}
	if cfg.temperature != testConfig.GetModelConfig().Temperature {
		t.Errorf("Expected temperature %f, got %f", testConfig.GetModelConfig().Temperature, cfg.temperature)
	}
	if cfg.maxTokens != testConfig.GetModelConfig().MaxTokens {
		t.Errorf("Expected maxTokens %d, got %d", testConfig.GetModelConfig().MaxTokens, cfg.maxTokens)
	}
	if cfg.topP != testConfig.GetModelConfig().TopP {
		t.Errorf("Expected topP %f, got %f", testConfig.GetModelConfig().TopP, cfg.topP)
	}
	if cfg.systemMessage != testConfig.GetModelConfig().System {
		t.Errorf("Expected system message '%s', got '%s'", testConfig.GetModelConfig().System, cfg.systemMessage)
	}
	if cfg.maxTurnsLimit != testConfig.GetAgentConfig().MaxTurns {
		t.Errorf("Expected maxTurnsLimit %d, got %d", testConfig.GetAgentConfig().MaxTurns, cfg.maxTurnsLimit)
	}

	// Verify messages count (system + user)
	if len(cfg.messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(cfg.messages))
	}

	// Verify executor is created
	if cfg.executor == nil {
		t.Error(errExpectedExecutorToBeCreated)
	}

	// Verify logger is set
	if cfg.logger == nil {
		t.Error("Expected logger to be set")
	}
}

// TestLoadAgentConfig_WithFlags tests loading agent configuration with command-line flags
func TestLoadAgentConfigWithFlags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Load test config from file
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
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
		testPrompt,
	}

	cfg, err := loadAgentConfig(mockClient, args, os.Stdin)
	if err != nil {
		t.Fatalf(errExpectedNoError, err)
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
		t.Error(errExpectedExecutorToBeCreated)
	}
}

// TestLoadAgentConfig_NoInput tests error when no input is provided
func TestLoadAgentConfigNoInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Load test config from file
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with no input (no arguments and stdin is a terminal)
	args := []string{}

	_, err = loadAgentConfig(mockClient, args, os.Stdin)
	if err == nil {
		t.Error("Expected error for no input, got nil")
	}
	if err.Error() != "no input provided" {
		t.Errorf("Expected error 'no input provided', got '%v'", err)
	}
}

// TestLoadAgentConfig_NoSystemMessage tests configuration without system message
func TestLoadAgentConfigNoSystemMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Load test config from file
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
	}

	// Modify config to have no system message
	testConfig.(*config.Config).Model.System = ""

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with command-line argument
	args := []string{testPrompt}

	cfg, err := loadAgentConfig(mockClient, args, os.Stdin)
	if err != nil {
		t.Fatalf(errExpectedNoError, err)
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
func TestLoadAgentConfigWithSystemMessageFlag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Load test config from file
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with system message flag
	args := []string{
		"--system", "Custom system message from flag",
		testPrompt,
	}

	cfg, err := loadAgentConfig(mockClient, args, os.Stdin)
	if err != nil {
		t.Fatalf(errExpectedNoError, err)
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
func TestLoadAgentConfigExecutionConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Load test config from file
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
	}

	// Modify config to have custom execution settings
	testConfig.(*config.Config).Execution.Engine = "docker run --rm alpine bash -c"
	testConfig.(*config.Config).Execution.Timeout = 60 * time.Second

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test without overriding execution config
	args := []string{testPrompt}

	cfg, err := loadAgentConfig(mockClient, args, os.Stdin)
	if err != nil {
		t.Fatalf(errExpectedNoError, err)
	}

	// Verify executor is created
	if cfg.executor == nil {
		t.Error(errExpectedExecutorToBeCreated)
	}
}

// TestLoadAgentConfig_ExecutionConfigOverride tests execution configuration override with flags
func TestLoadAgentConfigExecutionConfigOverride(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with default execution settings
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
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
		testPrompt,
	}

	cfg, err := loadAgentConfig(mockClient, args, os.Stdin)
	if err != nil {
		t.Fatalf(errExpectedNoError, err)
	}

	// Verify executor is created
	if cfg.executor == nil {
		t.Error(errExpectedExecutorToBeCreated)
	}
}

// TestLoadAgentConfig_AgentConfig tests agent configuration
func TestLoadAgentConfigAgentConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with agent settings
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
	}
	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test without overriding agent config
	args := []string{testPrompt}

	cfg, err := loadAgentConfig(mockClient, args, os.Stdin)
	if err != nil {
		t.Fatalf(errExpectedNoError, err)
	}

	// Verify max turns from config is used
	if cfg.maxTurnsLimit != 5 {
		t.Errorf("Expected maxTurnsLimit 5, got %d", cfg.maxTurnsLimit)
	}
}

// TestLoadAgentConfig_AgentConfigOverride tests agent configuration override with flags
func TestLoadAgentConfigAgentConfigOverride(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	mockClient := mock_openai.NewMockCLIClient(ctrl)

	// Create test config with default agent settings
	testConfig, err := config.Load(TestValidYAMLConfig)
	if err != nil {
		t.Fatalf(errFailedToLoadTestConfig, err)
	}

	// Create test logger
	testLogger := logger.NewLogger(false, false)

	// Set up mock expectations
	mockClient.EXPECT().GetCLIConfig().Return(testConfig)
	mockClient.EXPECT().GetLogger().Return(testLogger)

	// Test with max turns flag
	args := []string{
		"--max-turns", "15",
		testPrompt,
	}

	cfg, err := loadAgentConfig(mockClient, args, os.Stdin)
	if err != nil {
		t.Fatalf(errExpectedNoError, err)
	}

	// Verify max turns from flag is used
	if cfg.maxTurnsLimit != 15 {
		t.Errorf("Expected maxTurnsLimit 15, got %d", cfg.maxTurnsLimit)
	}
}
