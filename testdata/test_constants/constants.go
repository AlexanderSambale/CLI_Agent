package test_constants

import "time"

// Test constants
const (
	TestBaseURL             = "https://api.openai.com/v1"
	TestAPIKey              = "sk-test-key"
	TestConfigName          = "Test Config"
	TestVersion             = "1.0.0"
	TestName                = "Test"
	TestModel               = "glm-4.7"
	TestTemperature         = 0.7
	TestMaxTokens           = 128000
	TestTopP                = 1.0
	TestTimeout             = 120
	TestMaxRetries          = 3
	TestRetryDelay          = 1000
	TestAltModel            = "gpt-3.5-turbo"
	TestAltTimeout          = 60
	TestResponseContent     = "Test response content"
	ErrFailedToCreateClient = "Failed to create client: %v"

	// Execution config test constants
	TestExecutionEngine     = ""
	TestExecutionTimeout    = 30 * time.Second
	TestAltExecutionEngine  = "docker run --rm alpine:latest bash -c"
	TestAltExecutionTimeout = 2 * time.Second

	// Agent config test constants
	TestAgentSystem   = "You are a helpful coding assistant. Use <do>...</do> tags to wrap bash commands you want to execute."
	TestAgentMaxTurns = 5

	// Config file paths
	TestValidYAMLConfig = "../../testdata/config/valid.yaml"
	TestValidJSONConfig = "../../testdata/config/valid.json"
	TestValidTOMLConfig = "../../testdata/config/valid.toml"
	TestInvalidConfig   = "../../testdata/config/invalid.yaml"
)
