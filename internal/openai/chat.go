package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v3"
)

// ChatCompletionRequest represents a chat completion request
type ChatCompletionRequest struct {
	Model       string
	Messages    []openai.ChatCompletionMessageParamUnion
	Temperature *float64
	MaxTokens   *int
	TopP        *float64
}

// ChatCompletionResponse represents a chat completion response
type ChatCompletionResponse struct {
	ID      string
	Object  string
	Created int64
	Model   string
	Choices []openai.ChatCompletionChoice
	Usage   openai.CompletionUsage
}

// CreateChatCompletion creates a chat completion
func CreateChatCompletion(c CLIClient, ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	c.GetLogger().Verbosef("Creating chat completion with model: %s", req.Model)
	c.GetLogger().Verbosef("Number of messages: %d", len(req.Messages))

	// Build request parameters
	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(req.Model),
		Messages: req.Messages,
	}

	// Add optional parameters
	if req.Temperature != nil {
		params.Temperature = openai.Float(*req.Temperature)
	}
	if req.MaxTokens != nil {
		params.MaxTokens = openai.Int(int64(*req.MaxTokens))
	}
	if req.TopP != nil {
		params.TopP = openai.Float(*req.TopP)
	}

	// Make the API call
	completion, err := c.NewCompletion(ctx, params)
	if err != nil {
		c.GetLogger().Errorf("Chat completion failed: %v", err)
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	c.GetLogger().Verbosef("Chat completion created successfully")
	c.GetLogger().Verbosef("Response ID: %s", completion.ID)
	c.GetLogger().Verbosef("Tokens used: %d (prompt) + %d (completion) = %d (total)",
		completion.Usage.PromptTokens,
		completion.Usage.CompletionTokens,
		completion.Usage.TotalTokens)

	return &ChatCompletionResponse{
		ID:      completion.ID,
		Object:  string(completion.Object),
		Created: completion.Created,
		Model:   completion.Model,
		Choices: completion.Choices,
		Usage:   completion.Usage,
	}, nil
}
