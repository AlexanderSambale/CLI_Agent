# Architecture

## System Architecture

The CLI Agent follows a modular architecture with clear separation of concerns:

```
CLI_Agent/
├── cmd/                    # Command-line interface entry points
│   ├── root.go            # Main CLI command definition and execution
│   ├── chat.go            # Chat completions command
│   └── models.go          # Models API command
├── internal/              # Private application code
│   ├── config/           # Configuration parsing and validation
│   │   ├── config.go     # Config loading and validation logic
│   │   └── config_test.go # Config unit tests
│   ├── openai/           # OpenAI client wrapper
│   │   ├── client.go     # Client initialization and configuration
│   │   ├── chat.go       # Chat completions API
│   │   ├── models.go     # Models API
│   │   ├── errors.go     # Error handling and custom error types
│   │   ├── client_test.go # Client unit tests
│   │   └── chat_test.go  # Chat completion unit tests
│   ├── logger/           # Verbose logging utilities
│   │   └── logger.go     # Structured logging implementation
│   ├── parser/           # Bash command parser for agent
│   │   ├── parser.go     # Bash command extraction from LLM responses
│   │   └── parser_test.go # Parser unit tests
│   └── mocks/            # Generated mock clients for testing
│       └── client_mock.go # GoMock generated mock
├── tests/                # Integration tests and test helpers
│   ├── helpers.go        # Test helper functions
│   └── integration/      # Integration tests
│       ├── chat_test.go  # Chat command integration tests
│       └── models_test.go # Models command integration tests
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

## Component Relationships

1. **CLI Layer** ([`cmd/`](cmd/)): Handles user input, command parsing, and output formatting
   - [`root.go`](cmd/root.go): Main entry point, flag parsing, subcommand routing
   - [`chat.go`](cmd/chat.go): Chat completion command implementation
   - [`models.go`](cmd/models.go): Models listing and retrieval commands

2. **Configuration Layer** ([`internal/config/`](internal/config/)): Parses and validates configuration files
   - Supports YAML, JSON, and TOML formats
   - Validates required fields and sets defaults
   - Provides structured configuration objects

3. **OpenAI Client** ([`internal/openai/`](internal/openai/)): Wraps OpenAI API client library
   - [`client.go`](internal/openai/client.go): Client initialization with HTTP configuration
   - [`chat.go`](internal/openai/chat.go): Chat completion API wrapper
   - [`models.go`](internal/openai/models.go): Models API wrapper
   - [`errors.go`](internal/openai/errors.go): Custom error types and API error handling

4. **Logging** ([`internal/logger/`](internal/logger/)): Structured logging with multiple levels
   - INFO: General application information
   - VERBOSE: Detailed operation information (configurable)
   - DEBUG: Debug-level information for troubleshooting (configurable)
   - ERROR: Error messages and stack traces

5. **Parser** ([`internal/parser/`](internal/parser/)): Bash command extraction for agent capabilities
   - [`parser.go`](internal/parser/parser.go): Extracts bash commands from LLM responses using `<do>...</do>` tags
   - Supports multi-line commands and complex bash syntax
   - Comprehensive error handling for edge cases
   - 100+ test cases covering various bash command patterns

## Key Technical Decisions

- **Language**: Go (Golang) 1.25.5 for performance and simplicity
- **CLI Framework**: pflag (spf13/pflag) for POSIX-compliant flag parsing
- **Configuration**: viper (spf13/viper) for multi-format configuration support
- **OpenAI Client**: openai-go/v3 official library
- **Logging**: Custom structured logger with configurable verbosity
- **Error Handling**: Custom error types with wrapping for better error context
- **Bash Parser**: Custom regex-based parser using `<do>...</do>` tags for bash command extraction
- **Mocking**: GoMock (go.uber.org/mock) for advanced unit testing capabilities

## Design Patterns

- **Command Pattern**: CLI command structure with separate command files
- **Wrapper Pattern**: OpenAI client wrapper for additional functionality
- **Builder Pattern**: Configuration building with validation
- **Strategy Pattern**: Different configuration formats (YAML, JSON, TOML)
- **Parser Pattern**: Bash command extraction from LLM responses using custom tags

## Critical Implementation Paths

1. **Configuration Loading**: Config file path → [`config.Load()`](internal/config/config.go:48) → Format detection → Parsing → Validation → Config Object
2. **Client Initialization**: Config Object → [`openai.NewClient()`](internal/openai/client.go:23) → HTTP client setup → OpenAI client creation
3. **Chat Completion**: User input → [`ExecuteChat()`](cmd/chat.go:35) → Request building → [`CreateChatCompletion()`](internal/openai/chat.go:30) → API call → Response formatting
4. **Models Operations**: Command parsing → [`ExecuteModels()`](cmd/models.go:30) → [`ListModels()`](internal/openai/models.go:11) or [`GetModel()`](internal/openai/models.go:25) → API call → Output formatting

## Source Code Paths

- Main entry point: [`main.go`](main.go:1)
- CLI commands: [`cmd/`](cmd/)
- Configuration: [`internal/config/config.go`](internal/config/config.go:1)
- OpenAI client: [`internal/openai/`](internal/openai/)
- Logging: [`internal/logger/logger.go`](internal/logger/logger.go:1)
- Bash parser: [`internal/parser/parser.go`](internal/parser/parser.go:1)
- Example configuration: [`example.yaml`](example.yaml:1)

## Data Flow

```
User Input (CLI)
    ↓
cmd/root.go (flag parsing)
    ↓
cmd/chat.go or cmd/models.go (command execution)
    ↓
internal/openai/client.go (API interaction)
    ↓
internal/logger/logger.go (logging)
    ↓
Output to user
```

## Configuration Structure

```yaml
name: string              # Application name
version: string           # Application version
settings:
  debug: bool            # Enable debug logging
  verbose: bool          # Enable verbose logging
openai:
  base_url: string       # OpenAI API base URL (required)
  api_key: string        # OpenAI API key (required)
  http_client:
    timeout: int         # Request timeout in seconds
    max_retries: int     # Maximum number of retries
    retry_delay: int     # Delay between retries in milliseconds
  defaults:
    model: string        # Default model
    temperature: float64 # Default temperature
    max_tokens: int      # Default max tokens
    top_p: float64       # Default top_p
```

## Notes

The architecture is designed for extensibility. The OpenAI client wrapper can be extended to support additional API endpoints (embeddings, images, etc.). The configuration system can easily accommodate new fields. The logging system provides flexibility for different verbosity levels during development and production. The bash parser is ready for integration with agent capabilities to execute commands extracted from LLM responses.