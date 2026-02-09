# Tech

## Technologies Used

### Core Language

- **Go (Golang)**: Primary programming language
  - Version: 1.25.5
  - Reasoning: Performance, simplicity, strong standard library, excellent for CLI tools

### CLI Framework

- **pflag** (github.com/spf13/pflag v1.0.10): POSIX-compliant command-line flag parsing
  - Used for parsing command-line arguments and flags
  - Supports short and long flags, default values, and help text

### Configuration Management

- **viper** (github.com/spf13/viper v1.21.0): Configuration management library
  - Supports multiple configuration formats: YAML, JSON, TOML
  - Automatic format detection based on file extension
  - Environment variable support (not currently used)
  - Configuration validation and default value setting

### OpenAI API Client

- **openai-go/v3** (github.com/openai/openai-go/v3 v3.17.0): Official OpenAI Go client library
  - Type-safe API client for OpenAI services
  - Supports chat completions, models API, and other endpoints
  - Built-in error handling and response parsing

### Logging

- **Custom Logger** ([`internal/logger/logger.go`](internal/logger/logger.go:1)): Structured logging implementation
  - Multiple log levels: INFO, VERBOSE, DEBUG, ERROR
  - Timestamp formatting
  - Configurable verbosity levels
  - Output to stderr

## Development Setup

### Prerequisites

- Go 1.25.5 or later
- Git for version control
- Text editor or IDE with Go support

### Build System

- **go build**: Standard Go build tool

  ```bash
  go build -o cli-agent
  ```

- **go test**: Standard Go testing framework (not yet implemented)
- **go mod**: Go module management
  - Dependencies managed via [`go.mod`](go.mod:1)
  - Vendor directory not used (go.sum for checksums)

### Development Tools

- **gofmt**: Standard Go formatter (recommended)
- **golangci-lint**: Linting tool (recommended, not yet configured)
- **Testing**: go test with comprehensive test suite
- **GoMock** (go.uber.org/mock/gomock): Mock generation framework for unit testing

## Technical Constraints

- Must be cross-platform compatible (Linux, macOS, Windows)
- Minimal external dependencies (only 3 direct dependencies)
- Configuration files should be human-readable and editable
- CLI should follow standard Unix conventions
- Error messages should be clear and actionable
- No external template engine currently used
- No code analysis tools currently integrated

## Dependencies

### Direct Dependencies

- `github.com/spf13/pflag v1.0.10` - Command-line flag parsing
- `github.com/spf13/viper v1.21.0` - Configuration management
- `github.com/openai/openai-go/v3 v3.17.0` - OpenAI API client

### Indirect Dependencies

- `github.com/fsnotify/fsnotify v1.9.0` - File system notifications (viper)
- `github.com/go-viper/mapstructure/v2 v2.4.0` - Map structure decoding (viper)
- `github.com/pelletier/go-toml/v2 v2.2.4` - TOML parser (viper)
- `github.com/sagikazarmark/locafero v0.11.0` - File location utilities (viper)
- `github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8` - Concurrency utilities (viper)
- `github.com/spf13/afero v1.15.0` - File system abstraction (viper)
- `github.com/spf13/cast v1.10.0` - Type casting utilities (viper)
- `github.com/subosito/gotenv v1.6.0` - Environment variable parsing (viper)
- `github.com/tidwall/gjson v1.18.0` - JSON parsing (viper)
- `github.com/tidwall/match v1.2.1` - Pattern matching (viper)
- `github.com/tidwall/pretty v1.2.1` - JSON pretty printing (viper)
- `github.com/tidwall/sjson v1.2.5` - JSON setting (viper)
- `go.yaml.in/yaml/v3 v3.0.4` - YAML parser (viper)
- `golang.org/x/sys v0.29.0` - System call wrappers
- `golang.org/x/text v0.28.0` - Text processing utilities

## Tool Usage Patterns

### Version Control

- Git for source control
- Standard Git workflow (no specific branch strategy defined)
- Commit message conventions not yet established

### Code Organization

- Follows Go project layout standards
- Uses `internal/` package for private code
- Clear separation of concerns between packages
- Each package has a single responsibility

### Configuration Management

- Configuration files in project root (e.g., [`example.yaml`](example.yaml:1))
- Supports multiple formats via file extension detection
- Validation logic in [`config.Validate()`](internal/config/config.go:90)
- Default values set in validation method

### Error Handling

- Custom error types in [`internal/openai/errors.go`](internal/openai/errors.go:1)
- Error wrapping with context using `fmt.Errorf` and `%w`
- Structured error logging via logger
- API errors wrapped with custom [`APIError`](internal/openai/errors.go:29) type

### Testing Strategy

- Unit tests for core logic (config, client, chat packages)
- Integration tests for CLI commands (chat, models)
- Mock-based unit tests using GoMock
- Test fixtures and constants for consistent test data
- Test helpers for CLI command execution

## HTTP Client Configuration

The OpenAI client uses a custom HTTP client with configurable settings:

- **Timeout**: Configurable via `openai.http_client.timeout` (default: 120 seconds)
- **Max Retries**: Configurable via `openai.http_client.max_retries` (default: 3)
- **Retry Delay**: Configurable via `openai.http_client.retry_delay` (default: 1000ms)

Note: Retry logic is configured but not yet implemented in the client wrapper.

## Notes

The technology stack is stable and well-maintained. All dependencies are actively maintained and have good community support. The project uses minimal external dependencies to keep the binary size small and reduce attack surface. Future additions may include template engines, code analysis tools, and testing frameworks as the agent capabilities expand.
