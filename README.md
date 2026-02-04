# CLI Agent

A CLI-driven coding agent using a simple configuration system.

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
```
No configuration file specified. Use --config or -c to specify a config file.
CLI Agent is running!
```

### Running with a configuration file

```bash
./cli-agent --config configs/example.yaml
```

Or using the short flag:

```bash
./cli-agent -c configs/example.yaml
```

Output:
```
Loading configuration from: configs/example.yaml
Configuration loaded successfully!
CLI Agent is running!
```

### Supported Configuration Formats

The CLI agent supports the following configuration file formats:
- YAML (`.yaml`, `.yml`)
- JSON (`.json`)
- TOML (`.toml`)

### Example Configuration File

An example configuration file is provided at [`configs/example.yaml`](configs/example.yaml):

```yaml
name: "CLI Agent"
version: "0.1.0"

settings:
  debug: false
  verbose: true
```

## Project Structure

```
CLI_Agent/
├── cmd/                    # Command-line interface entry points
│   └── root.go            # Main CLI command definition
├── internal/              # Private application code
│   └── config/           # Configuration parsing and validation
│       └── config.go     # Config loading logic
├── configs/              # Example configuration files
│   └── example.yaml      # Example configuration
├── go.mod                # Go module definition
└── main.go               # Application entry point
```

## Development

### Dependencies

- [viper](https://github.com/spf13/viper) - Configuration management
- [pflag](https://github.com/spf13/pflag) - POSIX-compliant command-line flag parsing

### Adding Configuration Fields

To add new configuration fields:

1. Update the [`Config`](internal/config/config.go:11) struct in [`internal/config/config.go`](internal/config/config.go)
2. Add validation logic in the [`Validate()`](internal/config/config.go:68) method
3. Update the example configuration file in [`configs/`](configs/)