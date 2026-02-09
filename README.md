# CLI Agent

A CLI-driven client using a simple configuration system. This application provides a lightweight, efficient interface for interacting with the OpenAI API through command-line commands.

## Building

```bash
go build -o cli-agent
```

## Usage

### Running without a configuration file

```bash
./cli-agent
```

Output:

```bash
No configuration file specified. Use --config or -c to specify a config file.
Available commands:
  chat <prompt>    - Send a chat completion request
  models --list    - List all available models
  models --get <id> - Get details for a specific model
```

### Running with a configuration file

```bash
./cli-agent --config example.yaml
```

Or using the short flag:

```bash
./cli-agent -c example.yaml
```

## Commands

### Chat Completions

Send a chat completion request to the OpenAI API:

```bash
./cli-agent -c example.yaml chat "What is the capital of France?"
```

With custom parameters:

```bash
./cli-agent -c example.yaml chat --model gpt-4o --temperature 0.5 --max-tokens 1000 "Explain quantum computing"
```

With a system message:

```bash
./cli-agent -c example.yaml chat --system "You are a helpful assistant." "Hello!"
```

#### Chat Options

- `-m, --model`: Model to use for chat completion (default: from config)
- `-t, --temperature`: Sampling temperature (0-2, default: from config)
- `-n, --max-tokens`: Maximum tokens to generate (default: from config)
- `-p, --top-p`: Nucleus sampling threshold (0-1, default: from config)
- `-s, --system`: System message to set context

### Models

List all available models:

```bash
./cli-agent -c example.yaml models --list
```

Get details for a specific model:

```bash
./cli-agent -c example.yaml models --get gpt-4o
```

#### Models Options

- `-l, --list`: List all available models
- `-g, --get`: Get details for a specific model

## Configuration

### Supported Configuration Formats

The CLI agent supports the following configuration file formats:

- YAML (`.yaml`, `.yml`)
- JSON (`.json`)
- TOML (`.toml`)

### Configuration File Structure

An example configuration file is provided at [`example.yaml`](example.yaml):

```yaml
name: "CLI Agent OpenAI Client"
version: "1.0.0"

settings:
  debug: false
  verbose: true

openai:
  # Required: OpenAI API base URL
  base_url: "https://api.openai.com/v1"

  # Required: Your OpenAI API key
  api_key: "sk-..."

  # HTTP client configuration
  http_client:
    timeout: 120              # Request timeout in seconds
    max_retries: 3            # Maximum number of retries
    retry_delay: 1000         # Delay between retries in milliseconds

  # Default parameters for API requests
  defaults:
    model: "gpt-4o"
    temperature: 0.7
    max_tokens: 2048
    top_p: 1.0
```

### Configuration Fields

#### General Settings

- `name`: Application name
- `version`: Application version
- `settings.debug`: Enable debug logging
- `settings.verbose`: Enable verbose logging

#### OpenAI Configuration

- `openai.base_url`: OpenAI API base URL (required)
- `openai.api_key`: Your OpenAI API key (required)
- `openai.http_client.timeout`: Request timeout in seconds (default: 120)
- `openai.http_client.max_retries`: Maximum number of retries (default: 3)
- `openai.http_client.retry_delay`: Delay between retries in milliseconds (default: 1000)
- `openai.defaults.model`: Default model to use (default: gpt-4o)
- `openai.defaults.temperature`: Default sampling temperature (default: 0.7)
- `openai.defaults.max_tokens`: Default maximum tokens (default: 2048)
- `openai.defaults.top_p`: Default nucleus sampling threshold (default: 1.0)

## Project Structure

```bash
CLI_Agent/
├── cmd/                    # Command-line interface entry points
│   ├── root.go            # Main CLI command definition
│   ├── chat.go            # Chat completions command
│   └── models.go          # Models API command
├── internal/              # Private application code
│   ├── config/           # Configuration parsing and validation
│   │   ├── config.go     # Config loading logic
│   │   └── config_test.go # Config unit tests
│   ├── openai/           # OpenAI client wrapper
│   │   ├── client.go     # Client initialization
│   │   ├── chat.go       # Chat completions API
│   │   ├── models.go     # Models API
│   │   ├── errors.go     # Error handling
│   │   ├── client_test.go # Client unit tests
│   │   └── chat_test.go  # Chat completion unit tests
│   ├── logger/           # Verbose logging utilities
│   │   └── logger.go     # Structured logging implementation
│   └── mocks/            # Generated mock clients for testing
│       └── client_mock.go # GoMock generated mock
├── tests/                # Integration tests and test helpers
│   ├── helpers.go        # Test helper functions
│   └── integration/      # Integration tests
│       ├── chat_test.go  # Chat integration tests
│       └── models_test.go # Models integration tests
├── testdata/             # Test fixtures and constants
│   ├── config/           # Test configuration files
│   │   ├── valid.yaml    # Valid YAML config
│   │   ├── valid.json    # Valid JSON config
│   │   ├── valid.toml    # Valid TOML config
│   │   └── invalid.yaml  # Invalid config for testing
│   └── test_constants/   # Test constants
│       └── constants.go  # Shared test constants
├── example.yaml          # Example configuration file
├── go.mod                # Go module definition
└── main.go               # Application entry point
```

## Development

### Dependencies

- [viper](https://github.com/spf13/viper) - Configuration management
- [pflag](https://github.com/spf13/pflag) - POSIX-compliant command-line flag parsing
- [openai-go/v3](https://github.com/openai/openai-go) - OpenAI API client library
- [GoMock](https://github.com/uber-go/mock) - Mock generation framework for testing

### Adding Configuration Fields

To add new configuration fields:

1. Update the [`Config`](internal/config/config.go:11) struct in [`internal/config/config.go`](internal/config/config.go)
2. Add validation logic in the [`ValidateAndSetDefaults()`](internal/config/config.go:68) method
3. Update the example configuration file

### Verbose Logging

The application provides comprehensive logging at multiple levels:

- **INFO**: General information about application state
- **VERBOSE**: Detailed information about operations (enabled by `verbose: true`)
- **DEBUG**: Debug-level information for troubleshooting (enabled by `debug: true`)
- **ERROR**: Error messages and stack traces

All logs include timestamps and log levels for easy filtering and analysis.

### Error Handling

The application includes comprehensive error handling:

- Configuration validation errors with clear messages
- API errors wrapped with custom error types
- Network error handling with retry logic
- Detailed error logging in verbose mode

## Testing

The project includes comprehensive unit and integration tests to ensure code quality and functionality.

### Running Tests

#### Run All Tests

```bash
go test ./...
```

#### Run Unit Tests Only

```bash
go test ./internal/... ./cmd/...
```

#### Run Integration Tests Only

Integration tests require valid API credentials and are tagged with the `integration` build tag:

```bash
go test -tags=integration ./tests/...
PROJECT_ROOT="$(pwd)" go test ./tests/integration -tags=integration
```

#### Run Tests with Coverage

```bash
go test -cover ./...
```

#### Run Tests with Coverage Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Run Specific Test

```bash
go test -run TestLoadValidYAML ./internal/config
```

#### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Test Structure

```bash
CLI_Agent/
├── internal/
│   ├── config/
│   │   └── config_test.go       # Config loading and validation tests
│   ├── openai/
│   │   ├── client_test.go       # Client initialization tests
│   │   └── chat_test.go         # Chat completion tests
│   └── mocks/
│       └── client_mock.go       # GoMock generated mock client
├── testdata/
│   ├── config/
│   │   ├── valid.yaml           # Valid YAML config
│   │   ├── valid.json           # Valid JSON config
│   │   ├── valid.toml           # Valid TOML config
│   │   └── invalid.yaml         # Invalid config for testing
│   └── test_constants/
│       └── constants.go         # Shared test constants
├── tests/
│   ├── helpers.go               # Test helper functions
│   └── integration/
│       ├── chat_test.go         # Chat integration tests
│       └── models_test.go       # Models integration tests
```

### Unit Tests

Unit tests cover individual components in isolation:

- **Config Package**: Tests configuration loading, validation, and default value setting
- **OpenAI Package**: Tests client initialization, request building, and error handling

#### Mock-Based Unit Tests

Unit tests use GoMock to test functionality without requiring API calls:

- **Mock Client**: [`internal/mocks/client_mock.go`](internal/mocks/client_mock.go) provides a GoMock-generated test double for the OpenAI client
- **Test Constants**: [`testdata/test_constants/constants.go`](testdata/test_constants/constants.go) provides shared test constants

Example of using mock client in tests:

```go
// Create mock controller and client
ctrl := gomock.NewController(t)
mockClient := mocks.NewMockCLIClient(ctrl)

// Set up expected behavior
mockClient.EXPECT().
    NewCompletion(gomock.Any(), gomock.Any()).
    Return(MockChatCompletionResponse("Test response content"), nil)

// Create request
req := &ChatCompletionRequest{
    Model:       "gpt-4",
    Messages:    []openaiapi.ChatCompletionMessageParamUnion{openaiapi.UserMessage("test")},
    Temperature: f64(0.7),
    MaxTokens:   intP(2048),
    TopP:        f64(1.0),
}

// Execute with mock response
resp, err := CreateChatCompletion(mockClient, context.Background(), req)

// Verify response
if err != nil {
    t.Fatalf("Expected no error, got: %v", err)
}
if resp.Choices[0].Message.Content != "Test response content" {
    t.Errorf("Expected content 'Test response content', got '%s'", resp.Choices[0].Message.Content)
}
```

Benefits of mock-based testing:

- **No API Dependencies**: Tests run without requiring valid API credentials
- **Faster Tests**: Mock responses eliminate network latency
- **Consistent Test Data**: Using test constants ensures consistent mock data
- **Better Test Coverage**: Can test error scenarios that are hard to reproduce with real API

### Integration Tests

Integration tests test the complete CLI workflow with actual API calls. These tests require:

1. A valid configuration file with API credentials
2. The CLI binary to be built

#### Setting Up Integration Tests

1. Create or update your configuration file (e.g., `config.yaml`) with valid API credentials:

    ```yaml
    openai:
      base_url: "https://your-api-endpoint/v1"
      api_key: "your-api-key"
      # ... other settings
    ```

2. Set the `PROJECT_ROOT` environment variable to point to your project root:

    ```bash
    export PROJECT_ROOT=$(pwd)
    ```

3. Run integration tests:

    ```bash
    go test -tags=integration ./tests/...
    ```

#### Integration Test Coverage

- **Models List**: Tests `./cli-agent -c config.yaml models --list`
  - Verifies successful execution
  - Checks for "Owned By" pattern in output

- **Models Get**: Tests `./cli-agent -c config.yaml models --get <model-id>`
  - Verifies successful execution
  - Checks for "ID:" in output
  - Validates model details are displayed

- **Chat**: Tests `./cli-agent -c config.yaml chat "What is the capital of France?"`
  - Verifies successful execution
  - Validates response is received
  - Tests with invalid configuration

### Test Fixtures

Test fixtures are provided in the `testdata/` directory:

- **Configuration Files**: Valid and invalid configurations for testing different scenarios
- **Test Constants**: Shared test constants for consistent test data

The [`testdata/test_constants/constants.go`](testdata/test_constants/constants.go) file provides the following test constants:

- `TestBaseURL`: Test API base URL
- `TestAPIKey`: Test API key
- `TestConfigName`: Test configuration name
- `TestVersion`: Test version
- `TestModel`: Test model identifier
- `TestTemperature`: Test temperature value
- `TestMaxTokens`: Test max tokens value
- `TestTopP`: Test top_p value
- `TestTimeout`: Test timeout value
- `TestMaxRetries`: Test max retries value
- `TestRetryDelay`: Test retry delay value
- `TestResponseContent`: Test response content

These constants provide consistent test data across all tests, ensuring reproducibility and reducing test maintenance.

### Test Helpers

The [`tests/helpers.go`](tests/helpers.go) file provides utility functions for testing:

- `ConfigPathIfExisting`: Returns config path if file exists
- `RunCLICommand`: Execute CLI commands and capture output
- `GetRootAndCLIAgent`: Get project root and CLI agent binary path
- `GetFirstModel`: Extract first model from models list output

### Skipping Integration Tests

Integration tests can be skipped by using the short test flag:

```bash
go test -short ./internal/...
```

## Security Considerations

1. **API Key Protection**
   - Never log API keys
   - Consider using environment variables for sensitive data
   - Do not commit configuration files with API keys to version control

2. **HTTPS Only**
   - Ensure `base_url` uses HTTPS
   - The application will warn about insecure HTTP connections

3. **Input Validation**
   - All user inputs are validated
   - Prompts are sanitized before sending to the API
   - Request sizes are limited to prevent abuse
