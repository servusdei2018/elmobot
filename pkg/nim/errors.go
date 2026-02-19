package nim

import (
	"errors"
	"fmt"
)

var (
	ErrMissingAPIKey  = errors.New("API key is required")
	ErrMissingModel   = errors.New("model is required")
	ErrMissingMessages = errors.New("at least one message is required")
)

type APIError struct {
	StatusCode int
	Message    string
	Type       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("NIM API error (status %d, type %s): %s", e.StatusCode, e.Type, e.Message)
}
