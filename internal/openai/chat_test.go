package openai

import (
	"cli_agent/internal/logger"
	mocks "cli_agent/internal/mocks"
	tc "cli_agent/testdata/test_constants"
	"context"
	"testing"

	"github.com/openai/openai-go/v3"
	openaiapi "github.com/openai/openai-go/v3"
	gm "go.uber.org/mock/gomock"
)

const (
	errExpectedInvalidCreds = "Expected error with invalid credentials, got nil"
)

func TestCreateChatCompletionSuccess(t *testing.T) {
	ctrl := gm.NewController(t)
	mockClient := mocks.NewMockCLIClient(ctrl)
	mockClient.EXPECT().
		NewCompletion(gm.Any(), gm.Any()).
		Return(MockChatCompletionResponse(tc.TestResponseContent), nil)
	l := logger.NewLogger(false, false)
	mockClient.EXPECT().
		GetLogger().
		Return(l).
		AnyTimes()

	req := &ChatCompletionRequest{
		Model:       tc.TestModel,
		Messages:    []openaiapi.ChatCompletionMessageParamUnion{openaiapi.UserMessage("test")},
		Temperature: f64(tc.TestTemperature),
		MaxTokens:   intP(tc.TestMaxTokens),
		TopP:        f64(tc.TestTopP),
	}

	resp, err := CreateChatCompletion(mockClient, context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify content
	if resp.Choices[0].Message.Content != tc.TestResponseContent {
		t.Errorf("Expected content 'Test response content', got '%s'", resp.Choices[0].Message.Content)
	}
}

func f64(v float64) *float64 {
	return &v
}

func intP(v int) *int {
	return &v
}

func MockChatCompletionResponse(content string) *openai.ChatCompletion {
	return &openai.ChatCompletion{
		ID:      "chatcmpl-test",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   "gpt-4",
		Choices: []openai.ChatCompletionChoice{
			MockChatCompletionChoice(content),
		},
		Usage: openai.CompletionUsage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}
}

func MockChatCompletionChoice(content string) openai.ChatCompletionChoice {
	return openai.ChatCompletionChoice{
		Message: openai.ChatCompletionMessage{
			Role:    "assistant",
			Content: content,
		},
	}
}
