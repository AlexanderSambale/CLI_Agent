package cmd

import (
	"context"
	"fmt"

	"cli_agent/internal/executor"
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
)

func GetAgentCmdFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("agent", pflag.ExitOnError)

	flagSet.StringVarP(&agentModel, "model", "m", "", "Model to use for agent")
	flagSet.Float64VarP(&agentTemperature, "temperature", "t", 0, "Sampling temperature (0-2)")
	flagSet.IntVarP(&agentMaxTokens, "max-tokens", "n", 0, "Maximum tokens to generate")
	flagSet.Float64VarP(&agentTopP, "top-p", "p", 0, "Nucleus sampling threshold (0-1)")
	flagSet.StringVarP(&agentSystem, "system", "s", "", "System message to set context")
	flagSet.IntVarP(&agentMaxTurns, "max-turns", "x", 10, "Maximum number of agent turns before stopping")

	return flagSet

}

// ExecuteAgent runs the agent command
func ExecuteAgent(client openai.CLIClient, args []string) error {
	flagSet := GetAgentCmdFlagSet()
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// Get user prompt
	if flagSet.NArg() < 1 {
		return fmt.Errorf("prompt is required")
	}
	prompt := flagSet.Arg(0)

	// Get config for execution settings
	cfg := client.GetCLIConfig()
	execConfig := cfg.GetExecutionConfig()

	// Get logger
	log := client.GetLogger()

	// Use config defaults if not specified
	model := agentModel
	if model == "" {
		model = cfg.GetOpenAIConfig().Defaults.Model
	}

	temperature := agentTemperature
	if temperature == 0 {
		temperature = cfg.GetOpenAIConfig().Defaults.Temperature
	}

	maxTokens := agentMaxTokens
	if maxTokens == 0 {
		maxTokens = cfg.GetOpenAIConfig().Defaults.MaxTokens
	}

	topP := agentTopP
	if topP == 0 {
		topP = cfg.GetOpenAIConfig().Defaults.TopP
	}

	// Build initial messages
	messages := []openaiapi.ChatCompletionMessageParamUnion{}
	if agentSystem != "" {
		messages = append(messages, openaiapi.SystemMessage(agentSystem))
	}
	messages = append(messages, openaiapi.UserMessage(prompt))

	// Create executor
	exec := executor.NewExecutor(&execConfig)

	// Agent loop
	ctx := context.Background()
	turnCount := 0
	maxTurnsLimit := agentMaxTurns

	for {
		// Check turn limit
		if turnCount >= maxTurnsLimit {
			log.Verbosef("Reached maximum turns limit (%d), stopping", maxTurnsLimit)
			fmt.Println("\n[Agent reached maximum turns limit]")
			break
		}
		turnCount++

		log.Verbosef("Agent turn %d", turnCount)
		log.Verbosef("Sending request to OpenAI with %d messages", len(messages))

		// Create request
		req := &openai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			Temperature: &temperature,
			MaxTokens:   &maxTokens,
			TopP:        &topP,
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
		log.Verbosef("LLM response length: %d characters", len(assistantMessage))

		// Print the LLM output
		fmt.Printf("\n--- (Turn %d) ---\n%s\n", turnCount, assistantMessage)

		// Add assistant message to history
		messages = append(messages, openaiapi.AssistantMessage(assistantMessage))

		// Try to parse a bash command from the response
		command, err := parser.ExtractBashCommand(assistantMessage)
		if err != nil {
			// No command found - agent is done
			switch err {
			case parser.ErrNoBashAction:
				log.Verbosef("No bash action found, agent finished")
			case parser.ErrMultipleBashActions:
				fmt.Println("\n[Agent error: multiple bash actions found]")
				return fmt.Errorf("multiple bash actions found in response")
			case parser.ErrEmptyBashAction:
				fmt.Println("\n[Agent error: empty bash action]")
				return fmt.Errorf("empty bash action found")
			default:
				return fmt.Errorf("error parsing response: %w", err)
			}
			break
		}

		// Execute the command
		fmt.Printf("\n--- Executing ---\n%s\n", command)
		result, err := exec.Execute(ctx, command)
		if err != nil {
			log.Verbosef("Command execution error: %v", err)
		}

		// Build output message
		outputMessage := buildOutputMessage(result, err)

		// Print the output
		fmt.Printf("\n--- Output ---\n%s\n", outputMessage)

		// Add output to message history
		messages = append(messages, openaiapi.UserMessage(outputMessage))
	}

	return nil
}

// buildOutputMessage creates a formatted message from the execution result
func buildOutputMessage(result *executor.Result, execErr error) string {
	var output string

	if execErr != nil {
		output = fmt.Sprintf("Command execution error: %v\n", execErr)
	} else {
		output = "Command executed successfully.\n"
	}

	if result != nil {
		output += fmt.Sprintf("Exit code: %d\n", result.ExitCode)
		output += fmt.Sprintf("Duration: %v\n", result.Duration)

		if result.Stdout != "" {
			output += fmt.Sprintf("\n--- stdout ---\n%s\n", result.Stdout)
		}
		if result.Stderr != "" {
			output += fmt.Sprintf("\n--- stderr ---\n%s\n", result.Stderr)
		}
	}

	return output
}
