package nim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://integrate.api.nvidia.com/v1"
	defaultTimeout = 2 * time.Minute
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type ClientConfig struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.APIKey == "" {
		return nil, ErrMissingAPIKey
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		baseURL:    baseURL,
		apiKey:     config.APIKey,
		httpClient: httpClient,
	}, nil
}

func (c *Client) CreateCompletion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq, false)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var completion CompletionResponse
	if err := json.Unmarshal(respBody, &completion); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &completion, nil
}

func (c *Client) CreateCompletionStream(ctx context.Context, req *CompletionRequest) (<-chan StreamEvent, <-chan error, error) {
	if err := req.Validate(); err != nil {
		return nil, nil, err
	}

	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq, true)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, c.handleErrorResponse(resp)
	}

	eventChan := make(chan StreamEvent, 1)
	errChan := make(chan error, 1)

	go c.handleStream(resp.Body, eventChan, errChan)

	return eventChan, errChan, nil
}

func (c *Client) setHeaders(req *http.Request, stream bool) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	if stream {
		req.Header.Set("Accept", "text/event-stream")
	} else {
		req.Header.Set("Accept", "application/json")
	}
	req.Header.Set("Content-Type", "application/json")
}

func (c *Client) handleStream(body io.ReadCloser, eventChan chan<- StreamEvent, errChan chan<- error) {
	defer close(eventChan)
	defer close(errChan)
	defer body.Close()

	decoder := NewStreamDecoder(body)

	for {
		event, err := decoder.Decode()
		if err != nil {
			if err == io.EOF {
				return
			}
			errChan <- err
			return
		}

		if event != nil {
			eventChan <- *event
		}
	}
}

func (c *Client) handleErrorResponse(resp *http.Response) error {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("http %d: failed to read error response body: %w", resp.StatusCode, err)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err != nil {
		return fmt.Errorf("http %d: %s", resp.StatusCode, string(respBody))
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    errResp.Error.Message,
		Type:       errResp.Error.Type,
	}
}
