# Product

## Why This Project Exists

This project aims to create a CLI-driven coding agent that helps developers automate coding tasks through a simple configuration system. The agent is written in Go (Golang) and designed to be lightweight, efficient, and easy to use from the command line. It is inspired by the mini-SWE-agent project and follows similar principles of using bash-only tools, linear history, and stateless execution.

## Problems It Solves

- **Automated Code Generation**: Reduces repetitive coding tasks by automating them through a CLI interface
- **Simple Configuration**: Provides an easy-to-understand configuration system for defining coding tasks
- **Developer Productivity**: Helps developers focus on high-level logic while the agent handles boilerplate and repetitive code
- **OpenAI API Integration**: Provides a clean CLI interface to interact with OpenAI's chat completions and models API

## How It Should Work

The CLI agent currently:

1. Accepts configuration files that define OpenAI API settings and default parameters
2. Processes these configurations to initialize the OpenAI client
3. Executes through a command-line interface with clear options and feedback
4. Supports chat completions with customizable parameters (model, temperature, max tokens, top_p)
5. Provides models listing and retrieval capabilities

Future capabilities (planned):

1. Accept configuration files that define coding tasks and patterns
2. Process these configurations to generate or modify code
3. Execute bash commands for file operations and system interactions
4. Maintain conversation history for context-aware interactions
5. Support various code generation and modification patterns

## User Experience Goals

- **Intuitive CLI**: Clear, well-documented command-line interface
- **Fast Execution**: Leverage Go's performance for quick code generation
- **Flexible Configuration**: Support multiple configuration formats and patterns
- **Reliable Output**: Generate consistent, high-quality code
- **Easy Integration**: Simple to incorporate into existing development workflows
- **Stateless Execution**: Each command runs independently without persistent state
- **Bash-Only Tools**: Use standard bash commands for file operations and system interactions

## Target Users

- Developers who want to automate repetitive coding tasks
- Teams looking to standardize code generation patterns
- Individuals who prefer CLI tools over GUI-based IDE features
- Users who need a simple interface to OpenAI's API for coding assistance

## Current Capabilities

- **Chat Completions**: Send prompts to OpenAI models and receive responses
- **Models Management**: List available models and retrieve model details
- **Configuration Management**: Support for YAML, JSON, and TOML configuration files
- **Customizable Parameters**: Override default model, temperature, max tokens, and top_p via command-line flags
- **System Messages**: Set context with system messages for chat completions
- **Verbose Logging**: Configurable logging levels for debugging and monitoring

## Planned Capabilities

- **File Operations**: Read, write, and edit files through the CLI
- **Bash Tool Execution**: Execute bash commands for system interactions
- **Conversation History**: Maintain context across multiple interactions
- **Code Analysis**: Analyze existing code for context and modifications
- **Template-Based Generation**: Use templates for consistent code generation patterns
