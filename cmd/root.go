package cmd

import (
	"fmt"
	"os"

	"cli_agent/internal/config"

	"github.com/spf13/pflag"
)

var (
	// configFile is the path to the configuration file
	configFile string
)

// Execute runs the root command
func Execute() error {
	// Parse command-line flags
	parseFlags()

	// If a config file is specified, load it
	if configFile != "" {
		fmt.Printf("Loading configuration from: %s\n", configFile)
		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Validate the configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}

		fmt.Println("Configuration loaded successfully!")
	} else {
		fmt.Println("No configuration file specified. Use --config or -c to specify a config file.")
	}

	// For now, just print a message
	fmt.Println("CLI Agent is running!")

	return nil
}

// parseFlags sets up and parses command-line flags
func parseFlags() {
	flagSet := pflag.NewFlagSet("cli-agent", pflag.ExitOnError)

	// Define the config file flag
	flagSet.StringVarP(&configFile, "config", "c", "", "Path to the configuration file (supports YAML, JSON, TOML)")

	// Parse flags from command line
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}
}