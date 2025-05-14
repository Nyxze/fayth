package internal

import (
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
func (c *Client) ChatCompletion(prompt string) (ChatCompletionResponse, error) {
	// TODO
	return ChatCompletionResponse{}, nil
}

// OpenAI Completion API
//
// Deprecated: This API has been flagged as Legacy by OpenAI
//
// See: https://platform.openai.com/docs/api-reference/completions
func (c *Client) Completion(prompt string) (string, error) {
	// Fallback to ChatCompletion implementation
	response, err := c.ChatCompletion(prompt)
	if err != nil {

	}
	return response.Choices[0].Message.Content, nil
}

// OpenAI Responses API
//
// https://platform.openai.com/docs/api-reference/responses
func (c *Client) Response() {
	// TODO
}
