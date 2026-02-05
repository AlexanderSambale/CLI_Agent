package openai

import (
	"fmt"
	"net/http"
	"time"

	"cli_agent/internal/config"
	"cli_agent/internal/logger"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Client wraps the OpenAI client with additional functionality
type Client struct {
	openai.Client
	config *config.OpenAIConfig
	logger *logger.Logger
}

// NewClient creates a new OpenAI client with the given configuration
func NewClient(cfg *config.OpenAIConfig, log *logger.Logger) (*Client, error) {
	// Validate required fields
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base_url is required")
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	// Configure HTTP client
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.HTTPClient.Timeout) * time.Second,
	}

	// Build client options
	options := []option.RequestOption{
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
		option.WithHTTPClient(httpClient),
	}

	// Create the OpenAI client
	client := openai.NewClient(
		options...,
	)

	log.Verbose("OpenAI client initialized successfully")
	log.Verbosef("Base URL: %s", cfg.BaseURL)

	return &Client{
		Client: client,
		config: cfg,
		logger: log,
	}, nil
}

// GetConfig returns the client's configuration
func (c *Client) GetConfig() *config.OpenAIConfig {
	return c.config
}