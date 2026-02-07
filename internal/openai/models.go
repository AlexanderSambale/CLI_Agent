package openai

import (
	"context"
	"fmt"

	openaiapi "github.com/openai/openai-go/v3"
)

// ListModels lists all available models
func (c *Client) ListModels(ctx context.Context) ([]openaiapi.Model, error) {
	c.GetLogger().Verbosef("Listing available models")

	models, err := c.Client.Models.List(ctx)
	if err != nil {
		c.logger.Errorf("Failed to list models: %v", err)
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	c.logger.Verbosef("Found %d models", len(models.Data))
	return models.Data, nil
}

// GetModel retrieves a specific model
func (c *Client) GetModel(ctx context.Context, modelID string) (*openaiapi.Model, error) {
	c.GetLogger().Verbosef("Retrieving model: %s", modelID)

	model, err := c.Models.Get(ctx, modelID)
	if err != nil {
		c.logger.Errorf("Failed to retrieve model %s: %v", modelID, err)
		return nil, fmt.Errorf("failed to retrieve model: %w", err)
	}

	c.logger.Verbosef("Model retrieved: %s (owned by: %s)", model.ID, model.OwnedBy)
	return model, nil
}
