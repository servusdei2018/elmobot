package nim

import "fmt"

type CompletionRequest struct {
	Model            string         `json:"model"`
	Messages         []Message      `json:"messages"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	Stream           bool           `json:"stream"`
	TopK             int            `json:"top_k,omitempty"`
	RepetitionPenalty float32       `json:"repetition_penalty,omitempty"`
}

func (r *CompletionRequest) Validate() error {
	if r.Model == "" {
		return ErrMissingModel
	}
	if len(r.Messages) == 0 {
		return ErrMissingMessages
	}
	if r.Temperature < 0 || r.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2, got %f", r.Temperature)
	}
	if r.TopP < 0 || r.TopP > 1 {
		return fmt.Errorf("top_p must be between 0 and 1, got %f", r.TopP)
	}
	return nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	ID                string                `json:"id"`
	Object            string                `json:"object"`
	Created           int64                 `json:"created"`
	Model             string                `json:"model"`
	Choices           []Choice              `json:"choices"`
	Usage             Usage                 `json:"usage"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	LogProbs     *LogProbs `json:"logprobs"`
	FinishReason string   `json:"finish_reason"`
	MatchedStop  *string  `json:"matched_stop"`
}

type LogProbs struct {
	Content []LogProbContent `json:"content,omitempty"`
}

type LogProbContent struct {
	Token   string                   `json:"token"`
	LogProb float64                  `json:"logprob"`
	Bytes   []int                    `json:"bytes,omitempty"`
	TopLogProbs []map[string]interface{} `json:"top_logprobs,omitempty"`
}

type Usage struct {
	PromptTokens           int `json:"prompt_tokens"`
	CompletionTokens       int `json:"completion_tokens"`
	TotalTokens            int `json:"total_tokens"`
	PromptTokensDetails    *TokenDetails `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *TokenDetails `json:"completion_tokens_details,omitempty"`
	ReasoningTokens        int `json:"reasoning_tokens,omitempty"`
}

type TokenDetails struct {
	CachedTokens int `json:"cached_tokens,omitempty"`
	AudioTokens  int `json:"audio_tokens,omitempty"`
}

type StreamEvent struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []StreamChoice `json:"choices"`
	Usage   *Usage        `json:"usage,omitempty"`
}

type StreamChoice struct {
	Index        int        `json:"index"`
	Delta        Delta      `json:"delta"`
	LogProbs     *LogProbs  `json:"logprobs,omitempty"`
	FinishReason *string    `json:"finish_reason"`
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}
