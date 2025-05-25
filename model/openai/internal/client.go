package internal

import (
	"os"
)

// Represent the underlying client that
// communicate with OpenAI services
// Aggregate all OpenAI service in a single client
// Each subclients correspond to a given service, rather than providing all operations on a single client.
// See github.com/openai/openai-go for references
type Client struct {
	Options []CallOption
	Chat    ChatService
}

func NewClient(opts ...CallOption) (c Client) {
	opts = append(DefaultClientOptions(), opts...)
	c.Chat = NewChatService(opts...)
	return
}

// Set
func DefaultClientOptions() []CallOption {
	defaults := []CallOption{}
	if o, ok := os.LookupEnv(BASE_URL_ENV); ok {
		defaults = append(defaults, WithBaseURL(o))
	}
	if o, ok := os.LookupEnv(API_KEY_ENV); ok {
		defaults = append(defaults, WithAPIKey(o))
	}
	if o, ok := os.LookupEnv(ORG_ID_ENV); ok {
		defaults = append(defaults, WithOrganization(o))
	}
	if o, ok := os.LookupEnv(PROJECT_ID_ENV); ok {
		defaults = append(defaults, WithProject(o))
	}
	return defaults
}
