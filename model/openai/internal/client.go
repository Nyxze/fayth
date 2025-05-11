package internal

import "net/http"

// Represent the underlying client that
// communicate with OpenAI services
type Client struct {
	ApiKey       string
	organization string
	projectId    string
	client       *http.Client
}

func NewClient(apiKey string, opts ...ClientOptions) (*Client, error) {
	// Default client options
	client := &Client{
		client: http.DefaultClient,
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
func ChatCompletion(prompt string) (string, error) {
	// TODO
	return "", nil
}

// OpenAI Completion API
//
// Deprecated: This API has been flagged as Legacy by OpenAI
//
// See: https://platform.openai.com/docs/api-reference/completions
func Completion(prompt string) (string, error) {
	// TODO
	return "", nil
}

// OpenAI Responses API
//
// https://platform.openai.com/docs/api-reference/responses
func Response() {
	// TODO
}
