package openai

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cli_agent/internal/config"
	"cli_agent/internal/logger"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type CLIClient interface {
	GetChatService() openai.ChatService
	GetCLIConfig() config.CLIConfig
	GetLogger() logger.CLILogger
	CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)
	GetModel(ctx context.Context, modelID string) (*openai.Model, error)
	ListModels(ctx context.Context) ([]openai.Model, error)
}

// Client wraps the OpenAI client with additional functionality
type Client struct {
	openai.Client
	config *config.Config
	logger *logger.Logger
}

func (c *Client) GetCLIConfig() config.CLIConfig {
	return c.config
}

func (c *Client) GetChatService() openai.ChatService {
	return c.Client.Chat
}

func (c *Client) GetLogger() logger.CLILogger {
	return c.logger
}

// NewClient creates a new OpenAI client with the given configuration
func NewClient(cfg config.CLIConfig, log logger.CLILogger) (CLIClient, error) {
	// Validate required fields
	if cfg.GetOpenAIConfig().BaseURL == "" {
		return nil, fmt.Errorf("base_url is required")
	}
	if cfg.GetOpenAIConfig().APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	// Configure HTTP client
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.GetOpenAIConfig().HTTPClient.Timeout) * time.Second,
	}

	// Build client options
	options := []option.RequestOption{
		option.WithAPIKey(cfg.GetOpenAIConfig().APIKey),
		option.WithBaseURL(cfg.GetOpenAIConfig().BaseURL),
		option.WithHTTPClient(httpClient),
	}

	// Create the OpenAI client
	client := openai.NewClient(
		options...,
	)

	log.Verbosef("OpenAI client initialized successfully")
	log.Verbosef("Base URL: %s", cfg.GetOpenAIConfig().BaseURL)

	return &Client{
		Client: client,
		config: &config.Config{Name: cfg.GetName(), Version: cfg.GetVersion()},
		logger: &logger.Logger{Verbose: log.GetVerbose(), Debug: log.GetDebug(), Output: log.GetOutput()},
	}, nil
}
