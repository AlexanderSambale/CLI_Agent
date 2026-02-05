package cmd

import (
	"context"
	"fmt"

	"cli_agent/internal/openai"

	openaiapi "github.com/openai/openai-go/v3"
	"github.com/spf13/pflag"
)

var (
	chatModel       string
	chatTemperature float64
	chatMaxTokens   int
	chatTopP        float64
	chatSystem      string
)

// ChatCmd represents the chat command
func ChatCmd(client *openai.Client) *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("chat", pflag.ExitOnError)

	flagSet.StringVarP(&chatModel, "model", "m", "", "Model to use for chat completion")
	flagSet.Float64VarP(&chatTemperature, "temperature", "t", 0, "Sampling temperature (0-2)")
	flagSet.IntVarP(&chatMaxTokens, "max-tokens", "n", 0, "Maximum tokens to generate")
	flagSet.Float64VarP(&chatTopP, "top-p", "p", 0, "Nucleus sampling threshold (0-1)")
	flagSet.StringVarP(&chatSystem, "system", "s", "", "System message to set context")

	return flagSet
}

// ExecuteChat runs the chat command
func ExecuteChat(client *openai.Client, args []string) error {
	flagSet := ChatCmd(client)
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	// Get user prompt
	if flagSet.NArg() < 1 {
		return fmt.Errorf("prompt is required")
	}
	prompt := flagSet.Arg(0)

	// Build messages
	messages := []openaiapi.ChatCompletionMessageParamUnion{}
	if chatSystem != "" {
		messages = append(messages, openaiapi.SystemMessage(chatSystem))
	}
	messages = append(messages, openaiapi.UserMessage(prompt))

	// Use config defaults if not specified
	cfg := client.GetConfig()
	model := chatModel
	if model == "" {
		model = cfg.Defaults.Model
	}

	temperature := chatTemperature
	if temperature == 0 {
		temperature = cfg.Defaults.Temperature
	}

	maxTokens := chatMaxTokens
	if maxTokens == 0 {
		maxTokens = cfg.Defaults.MaxTokens
	}

	topP := chatTopP
	if topP == 0 {
		topP = cfg.Defaults.TopP
	}

	// Create request
	req := &openai.ChatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: &temperature,
		MaxTokens:   &maxTokens,
		TopP:        &topP,
	}

	// Execute
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}

	// Print response
	for _, choice := range resp.Choices {
		fmt.Println(choice.Message.Content)
	}

	return nil
}