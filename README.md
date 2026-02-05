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
│   │   └── config.go     # Config loading logic
│   ├── openai/           # OpenAI client wrapper
│   │   ├── client.go     # Client initialization
│   │   ├── chat.go       # Chat completions API
│   │   ├── models.go     # Models API
│   │   └── errors.go     # Error handling
│   └── logger/           # Verbose logging utilities
│       └── logger.go     # Structured logging implementation
├── example.yaml          # Example configuration file
├── go.mod                # Go module definition
└── main.go               # Application entry point
```

## Development

### Dependencies

- [viper](https://github.com/spf13/viper) - Configuration management
- [pflag](https://github.com/spf13/pflag) - POSIX-compliant command-line flag parsing
- [openai-go/v3](https://github.com/openai/openai-go) - OpenAI API client library

### Adding Configuration Fields

To add new configuration fields:

1. Update the [`Config`](internal/config/config.go:11) struct in [`internal/config/config.go`](internal/config/config.go)
2. Add validation logic in the [`Validate()`](internal/config/config.go:68) method
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
