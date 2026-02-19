package nim

import (
	"bytes"
	"io"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
		errType error
	}{
		{
			name:    "valid config",
			config:  ClientConfig{APIKey: "test-key"},
			wantErr: false,
		},
		{
			name:    "missing API key",
			config:  ClientConfig{},
			wantErr: true,
			errType: ErrMissingAPIKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("NewClient() error = %v, want %v", err, tt.errType)
			}
			if !tt.wantErr && client == nil {
				t.Errorf("NewClient() returned nil client")
			}
		})
	}
}

func TestCompletionRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     *CompletionRequest
		wantErr bool
	}{
		{
			name:    "missing model",
			req:     &CompletionRequest{Messages: []Message{{Role: "user", Content: "test"}}},
			wantErr: true,
		},
		{
			name:    "missing messages",
			req:     &CompletionRequest{Model: "test-model"},
			wantErr: true,
		},
		{
			name:    "invalid temperature too high",
			req:     &CompletionRequest{Model: "test", Messages: []Message{{Role: "user", Content: "test"}}, Temperature: 3.0},
			wantErr: true,
		},
		{
			name:    "invalid temperature negative",
			req:     &CompletionRequest{Model: "test", Messages: []Message{{Role: "user", Content: "test"}}, Temperature: -0.5},
			wantErr: true,
		},
		{
			name:    "invalid top_p too high",
			req:     &CompletionRequest{Model: "test", Messages: []Message{{Role: "user", Content: "test"}}, TopP: 1.5},
			wantErr: true,
		},
		{
			name:    "invalid top_p negative",
			req:     &CompletionRequest{Model: "test", Messages: []Message{{Role: "user", Content: "test"}}, TopP: -0.1},
			wantErr: true,
		},
		{
			name: "valid request",
			req: &CompletionRequest{
				Model:       "test",
				Messages:    []Message{{Role: "user", Content: "test"}},
				Temperature: 1.0,
				TopP:        0.9,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStreamDecoder(t *testing.T) {
	testData := `data: {"id":"1","choices":[{"delta":{"content":"Hello"}}]}
data: {"id":"2","choices":[{"delta":{"content":" world"}}]}
data: [DONE]
`
	decoder := NewStreamDecoder(bytes.NewReader([]byte(testData)))

	// First event
	event, err := decoder.Decode()
	if err != nil {
		t.Fatalf("failed to decode first event: %v", err)
	}
	if event.ID != "1" {
		t.Errorf("unexpected ID: %s", event.ID)
	}

	// Second event
	event, err = decoder.Decode()
	if err != nil {
		t.Fatalf("failed to decode second event: %v", err)
	}
	if event.ID != "2" {
		t.Errorf("unexpected ID: %s", event.ID)
	}

	// [DONE] should return EOF
	event, err = decoder.Decode()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
	if event != nil {
		t.Errorf("expected nil event on EOF")
	}
}
