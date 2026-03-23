package cmd

import (
	"context"
	"fmt"
	"time"

	"cli_agent/internal/executor"
	"cli_agent/internal/logger"
	"cli_agent/internal/openai"
	"cli_agent/internal/parser"

	openaiapi "github.com/openai/openai-go/v3"
	"github.com/spf13/pflag"
)

var (
	agentModel       string
	agentTemperature float64
	agentMaxTokens   int
	agentTopP        float64
	agentSystem      string
	agentMaxTurns    int
	agentEngine      string
	agentTimeout     int
)

func GetAgentCmdFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("agent", pflag.ExitOnError)

	flagSet.StringVarP(&agentModel, "model", "m", "", "Model to use for agent")
	flagSet.Float64VarP(&agentTemperature, "temperature", "t", 0, "Sampling temperature (0-2)")
	flagSet.IntVarP(&agentMaxTokens, "max-tokens", "n", 0, "Maximum tokens to generate")
	flagSet.Float64VarP(&agentTopP, "top-p", "p", 0, "Nucleus sampling threshold (0-1)")
	flagSet.StringVarP(&agentSystem, "system", "s", "", "System message to set context")
	flagSet.IntVarP(&agentMaxTurns, "max-turns", "x", 10, "Maximum number of agent turns before stopping")
	flagSet.StringVarP(&agentEngine, "engine", "e", "", "Command execution engine prefix (e.g., 'docker run --rm ubuntu bash -c')")
	flagSet.IntVarP(&agentTimeout, "timeout", "T", 0, "Command execution timeout in seconds")

	return flagSet

}

// agentRuntimeConfig holds the configuration for the agent execution
type agentRuntimeConfig struct {
	model         string
	temperature   float64
	maxTokens     int
	topP          float64
	systemMessage string
	maxTurnsLimit int
	messages      []openaiapi.ChatCompletionMessageParamUnion
	executor      executor.Executor
	logger        logger.CLILogger
}

// loadAgentConfig loads and validates the agent configuration from flags and config
func loadAgentConfig(client openai.CLIClient, args []string) (*agentRuntimeConfig, error) {
	flagSet := GetAgentCmdFlagSet()
	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	// Get user prompt
	prompt, err := readInput(flagSet)
	if err != nil {
		return nil, err
	}
	if prompt == "" {
		return nil, fmt.Errorf("no input provided")
	}

	// Get config for execution settings
	cfg := client.GetCLIConfig()
	execConfig := cfg.GetExecutionConfig()

	// Get logger
	log := client.GetLogger()

	// Use config defaults if not specified
	modelConfig := cfg.GetModelConfig()
	agentConfig := cfg.GetAgentConfig()

	// Use flag if changed, otherwise use config, otherwise use default
	model := agentModel
	if !flagSet.Changed("model") {
		model = modelConfig.Model
	}

	temperature := agentTemperature
	if !flagSet.Changed("temperature") {
		temperature = modelConfig.Temperature
	}

	maxTokens := agentMaxTokens
	if !flagSet.Changed("max-tokens") {
		maxTokens = modelConfig.MaxTokens
	}

	topP := agentTopP
	if !flagSet.Changed("top-p") {
		topP = modelConfig.TopP
	}

	// Use command-line flag or config default for system message
	systemMessage := agentSystem
	if !flagSet.Changed("system") {
		systemMessage = modelConfig.System
	}

	// Use command-line flag or config default for max turns
	maxTurnsLimit := agentMaxTurns
	if !flagSet.Changed("max-turns") {
		maxTurnsLimit = agentConfig.MaxTurns
	}

	// Use command-line flag or config default for engine
	engine := agentEngine
	if !flagSet.Changed("engine") {
		engine = execConfig.Engine
	}

	// Use command-line flag or config default for timeout
	timeout := agentTimeout
	if !flagSet.Changed("timeout") {
		timeout = int(execConfig.Timeout.Seconds())
	}

	// Build initial messages
	messages := []openaiapi.ChatCompletionMessageParamUnion{}
	if systemMessage != "" {
		messages = append(messages, openaiapi.SystemMessage(systemMessage))
	}
	messages = append(messages, openaiapi.UserMessage(prompt))

	// Create executor with overridden values
	overriddenExecConfig := execConfig
	overriddenExecConfig.Engine = engine
	overriddenExecConfig.Timeout = time.Duration(timeout) * time.Second
	exec := executor.NewExecutor(&overriddenExecConfig)

	return &agentRuntimeConfig{
		model:         model,
		temperature:   temperature,
		maxTokens:     maxTokens,
		topP:          topP,
		systemMessage: systemMessage,
		maxTurnsLimit: maxTurnsLimit,
		messages:      messages,
		executor:      exec,
		logger:        log,
	}, nil
}

// runAgentLoop executes the agent loop with the given configuration
func runAgentLoop(client openai.CLIClient, cfg *agentRuntimeConfig) error {
	ctx := context.Background()
	turnCount := 0
	consecutiveNoCommandTurns := 0

	for {
		// Check turn limit
		if turnCount >= cfg.maxTurnsLimit {
			cfg.logger.Verbosef("Reached maximum turns limit (%d), stopping", cfg.maxTurnsLimit)
			fmt.Println("\n[Agent reached maximum turns limit]")
			break
		}
		turnCount++

		cfg.logger.Verbosef("Agent turn %d", turnCount)
		cfg.logger.Verbosef("Sending request to OpenAI with %d messages", len(cfg.messages))

		// Create request
		req := &openai.ChatCompletionRequest{
			Model:       cfg.model,
			Messages:    cfg.messages,
			Temperature: &cfg.temperature,
			MaxTokens:   &cfg.maxTokens,
			TopP:        &cfg.topP,
		}

		// Query the LLM
		resp, err := openai.CreateChatCompletion(client, ctx, req)
		if err != nil {
			return fmt.Errorf("chat completion failed: %w", err)
		}

		// Get assistant response
		if len(resp.Choices) == 0 {
			return fmt.Errorf("no response choices returned from API")
		}

		assistantMessage := resp.Choices[0].Message.Content
		cfg.logger.Verbosef("LLM response length: %d characters", len(assistantMessage))

		// Print the LLM output
		cfg.logger.Verbosef("\n--- (Turn %d) ---\n%s\n", turnCount, assistantMessage)

		// Add assistant message to history
		cfg.messages = append(cfg.messages, openaiapi.AssistantMessage(assistantMessage))

		// Try to parse a bash command from the response
		command, err := parser.ExtractBashCommand(assistantMessage)
		if err != nil {
			// Handle parsing errors
			feedbackMessage := ""
			switch err {
			case parser.ErrNoBashAction:
				// No bash action found - increment counter
				consecutiveNoCommandTurns++
				feedbackMessage = fmt.Sprintf("No bash command was found in your response.")

				// Check if this is the second consecutive turn with no command
				if consecutiveNoCommandTurns >= 2 {
					cfg.logger.Verbosef("No bash command found in two consecutive turns, stopping")
					fmt.Println("\n[Agent stopped: No bash command found in two consecutive turns]")
					return nil
				}
			case parser.ErrMultipleBashActions:
				feedbackMessage = fmt.Sprintf("Agent error: multiple bash actions found in response")
			case parser.ErrEmptyBashAction:
				feedbackMessage = fmt.Sprintf("Agent error: empty bash action found")
			default:
				feedbackMessage = fmt.Sprintf("error parsing response: %w", err)
			}
			cfg.logger.Verbosef(feedbackMessage)
			cfg.messages = append(cfg.messages, openaiapi.UserMessage(feedbackMessage))
			continue
		}

		// Reset counter when a command is found
		consecutiveNoCommandTurns = 0

		// Execute the command
		fmt.Printf("--- Executing ---\n%s\n", command)
		result, err := cfg.executor.Execute(ctx, command)
		if err != nil {
			cfg.logger.Verbosef("Command execution error: %v", err)
		}

		// Build output message
		outputMessage := buildOutputMessage(result, err, cfg.logger)

		// Print the output
		fmt.Printf("\n--- Output ---\n%s\n", outputMessage)

		// Add output to message history
		cfg.messages = append(cfg.messages, openaiapi.UserMessage(outputMessage))
	}

	return nil
}

// ExecuteAgent runs the agent command
func ExecuteAgent(client openai.CLIClient, args []string) error {
	cfg, err := loadAgentConfig(client, args)
	if err != nil {
		return err
	}

	return runAgentLoop(client, cfg)
}

// buildOutputMessage creates a formatted message from the execution result
func buildOutputMessage(result *executor.Result, execErr error, log logger.CLILogger) string {
	var output string

	if execErr != nil {
		output = fmt.Sprintf("Command execution error: %v\n", execErr)
	} else {
		output = ""
	}

	if result != nil {
		log.Verbosef("Duration: %v\n", result.Duration)

		if result.Stdout != "" {
			output += result.Stdout
		}
		if result.Stderr != "" {
			output += fmt.Sprintf("\n--- stderr ---\n%s\n", result.Stderr)
		}
	}

	return output
}
