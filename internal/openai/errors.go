package openai

import (
	"errors"
	"fmt"
)

const (
	missingRequiredConfiguration = "missing required configuration"
	invalidApiKey                = "invalid API key"
	rateLimitExceeded            = "rate limit exceeded"
	invalidRequest               = "invalid request"
	modelNotFound                = "model not found"
	requestTimeout               = "request timeout"
)

var (
	// ErrMissingConfig is returned when required configuration is missing
	ErrMissingConfig = errors.New(missingRequiredConfiguration)

	// ErrInvalidAPIKey is returned when the API key is invalid
	ErrInvalidAPIKey = errors.New(invalidApiKey)

	// ErrRateLimited is returned when rate limit is exceeded
	ErrRateLimited = errors.New(rateLimitExceeded)

	// ErrInvalidRequest is returned for invalid requests
	ErrInvalidRequest = errors.New(invalidRequest)

	// ErrModelNotFound is returned when a model is not found
	ErrModelNotFound = errors.New(modelNotFound)

	// ErrTimeout is returned when a request times out
	ErrTimeout = errors.New(requestTimeout)
)

// APIError represents an OpenAI API error
type APIError struct {
	Type    string
	Message string
	Code    string
	Err     error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error [type=%s, code=%s]: %s", e.Type, e.Code, e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// NewAPIError creates a new API error
func NewAPIError(errType, message, code string, err error) *APIError {
	return &APIError{
		Type:    errType,
		Message: message,
		Code:    code,
		Err:     err,
	}
}
