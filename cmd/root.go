package cmd

import (
	"fmt"
	"os"

	"cli_agent/internal/config"
	"cli_agent/internal/logger"
	"cli_agent/internal/openai"

	"github.com/spf13/pflag"
)

var (
	// configFile is the path to the configuration file
	configFile string
)

// Execute runs the root command
func Execute() error {
	// Parse command-line flags
	flagSet := parseFlags()

	// Check if a subcommand is provided
	if flagSet.NArg() > 0 {
		subcommand := flagSet.Arg(0)
		// If a config file is specified, load it first
		if configFile != "" {
			client, err := initializeClient()
			if err != nil {
				return err
			}

			// Execute the subcommand
			return executeSubcommand(client, subcommand, flagSet.Args()[1:])
		}
		return fmt.Errorf("configuration file is required for subcommands. Use --config or -c to specify a config file")
	}

	// If a config file is specified, load it
	if configFile != "" {
		_, err := initializeClient()
		if err != nil {
			return err
		}

		fmt.Println("CLI Agent OpenAI Client is running!")
		fmt.Println("Available commands:")
		fmt.Println("  chat <prompt>    - Send a chat completion request")
		fmt.Println("  models --list    - List all available models")
		fmt.Println("  models --get <id> - Get details for a specific model")
		fmt.Println("  parse <text>     - Extract bash command from text using <do>...</do> tags")
		fmt.Println("  execute <command> - Execute a bash command")
		return nil
	}

	fmt.Println("No configuration file specified. Use --config or -c to specify a config file.")
	fmt.Println("Available commands:")
	fmt.Println("  chat <prompt>    - Send a chat completion request")
	fmt.Println("  models --list    - List all available models")
	fmt.Println("  models --get <id> - Get details for a specific model")
	fmt.Println("  parse <text>     - Extract bash command from text using <do>...</do> tags")
	fmt.Println("  execute <command> - Execute a bash command")
	return nil
}

// initializeClient loads the configuration and initializes the OpenAI client
func initializeClient() (openai.CLIClient, error) {
	fmt.Printf("Loading configuration from: %s\n", configFile)
	cfg, err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate the configuration
	if err := config.ValidateAndSetDefaults(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	fmt.Println("Configuration loaded successfully!")

	// Create logger
	log := logger.NewLogger(cfg.GetVerbose(), cfg.GetDebug())

	// Initialize OpenAI client
	client, err := openai.NewClient(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenAI client: %w", err)
	}

	return client, nil
}

// executeSubcommand executes the specified subcommand
func executeSubcommand(client openai.CLIClient, subcommand string, args []string) error {
	switch subcommand {
	case "chat":
		return ExecuteChat(client, args)
	case "models":
		return ExecuteModels(client, args)
	case "parse":
		return ExecuteParse(args)
	case "execute":
		return ExecuteExecute(args)
	default:
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}

// parseFlags sets up and parses command-line flags
func parseFlags() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("cli-agent", pflag.ExitOnError)

	// Define the config file flag
	flagSet.StringVarP(&configFile, "config", "c", "", "Path to the configuration file (supports YAML, JSON, TOML)")

	// Disable interspersed flags so parsing stops at the first positional argument (subcommand)
	// This allows subcommands to have their own flags
	flagSet.SetInterspersed(false)

	// Parse flags from command line
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	return flagSet
}
