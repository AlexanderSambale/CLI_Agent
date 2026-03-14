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

const (
	BaseURLIsRequired = "base_url is required"
	ApiKeyIsRequired  = "api_key is required"
)

type CLIClient interface {
	GetCLIConfig() config.CLIConfig
	GetLogger() logger.CLILogger
	GetModel(ctx context.Context, modelID string) (*openai.Model, error)
	ListModels(ctx context.Context) ([]openai.Model, error)
	NewCompletion(ctx context.Context, body openai.ChatCompletionNewParams, opts ...option.RequestOption) (res *openai.ChatCompletion, err error)
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

func (c *Client) NewCompletion(ctx context.Context, body openai.ChatCompletionNewParams, opts ...option.RequestOption) (res *openai.ChatCompletion, err error) {
	return c.Client.Chat.Completions.New(ctx, body, opts...)
}

func (c *Client) GetLogger() logger.CLILogger {
	return c.logger
}

// NewClient creates a new OpenAI client with the given configuration
func NewClient(cfg config.CLIConfig, log logger.CLILogger) (CLIClient, error) {
	OpenAIConfig := cfg.GetOpenAIConfig()
	BaseURL := OpenAIConfig.BaseURL
	APIKey := OpenAIConfig.APIKey
	// Validate required fields
	if BaseURL == "" {
		return nil, fmt.Errorf(BaseURLIsRequired)
	}
	if APIKey == "" {
		return nil, fmt.Errorf(ApiKeyIsRequired)
	}

	// Configure HTTP client
	httpClient := &http.Client{
		Timeout: time.Duration(OpenAIConfig.HTTPClient.Timeout) * time.Second,
	}

	// Build client options
	options := []option.RequestOption{
		option.WithAPIKey(APIKey),
		option.WithBaseURL(BaseURL),
		option.WithHTTPClient(httpClient),
	}

	// Create the OpenAI client
	client := openai.NewClient(
		options...,
	)

	log.Verbosef("OpenAI client initialized successfully")
	log.Verbosef("Base URL: %s", BaseURL)

	return &Client{
		Client: client,
		config: &config.Config{
			Name:     cfg.GetName(),
			Version:  cfg.GetVersion(),
			Settings: config.SettingsConfig{Debug: cfg.GetDebug(), Verbose: cfg.GetVerbose()},
			OpenAI:   OpenAIConfig,
			Execution: cfg.GetExecutionConfig(),
			Agent:     cfg.GetAgentConfig(),
		},
		logger: &logger.Logger{Verbose: log.GetVerbose(), Debug: log.GetDebug(), Output: log.GetOutput()},
	}, nil
}
