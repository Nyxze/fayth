package internal

import (
	"context"
	"net/http"
)

// Represent the underlying client that
// communicate with OpenAI services
type Client struct {
	ApiKey       string
	Model        string
	organization string
	projectId    string
	httpClient   *http.Client
	baseUrl      string
}

func NewClient(apiKey string, opts ...ClientOptions) (*Client, error) {
	// Default client options
	client := &Client{
		httpClient: http.DefaultClient,
		baseUrl:    defaultBaseURL,
	}
	// Overrides
	for _, option := range opts {
		option(client)
	}
	return client, nil
}

// OpenAI ChatCompletion API
//
// https://platform.openai.com/docs/api-reference/chat
func (c *Client) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (response *ChatCompletionResponse, err error) {
	// Validate request
	// Create httpRequest

	// Execute

	// Unmarshal Response

	return
}

// OpenAI Completion API
//
// Deprecated: This API has been flagged as Legacy by OpenAI
//
// See: https://platform.openai.com/docs/api-reference/completions
func (c *Client) Completion(ctx context.Context, prompt string) (string, error) {
	// Fallback to ChatCompletion implementation
	return "", nil
}

// OpenAI Responses API
//
// https://platform.openai.com/docs/api-reference/responses
func (c *Client) Response() {
	// TODO
}
