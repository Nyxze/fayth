package openai

import (
	"context"
	"errors"
	"fmt"
	"nyxze/fayth/model"
	"nyxze/fayth/model/openai/internal"
	"os"
)

var (
	ErrNoContentInResponse = errors.New("no content in generation response")
	ErrModelGen            = errors.New("failed to convert to generation type")
	ErrInvalidMimeType     = errors.New("invalid mime type on content")
)

type chatModel struct {
	client *internal.Client
}

// Compile type interface assertion
var _ model.Model = (*chatModel)(nil)

// Return a New OpenAI [model.Model]
func New(configs ...ModelOptions) (*chatModel, error) {
	// Default options
	config := &options{
		ApiKey: os.Getenv(API_KEY_ENV),
		Model:  os.Getenv(MODEL_NAME_ENV),
	}

	// Apply overrides
	for _, conf := range configs {
		conf(config)
	}

	// Validate config
	if config.ApiKey == "" {
		return nil, fmt.Errorf("missing API key")
	}
	if config.Model == "" {
		return nil, fmt.Errorf("missing model")
	}
	client, err := internal.NewClient(config.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to creae openai client %v", err)
	}
	return &chatModel{
		client: client,
	}, nil
}

// Type Alias
type chatMessage = internal.ChatMessage
type chatRequest = internal.ChatCompletionRequest

// [model.Model] implementation
func (o chatModel) Generate(ctx context.Context, messages []model.Message) (*model.Generation, error) {

	// TOOO: Add options call
	// Apply options call

	// Convert model.Message => ChatMessage
	chatMsg := make([]*chatMessage, len(messages))
	for _, m := range messages {
		msg := &chatMessage{}
		switch m.Role {
		case model.User:
			msg.Role = internal.RoleUser
		case model.Assistant:
			msg.Role = internal.RoleAssistant
		case model.System:
			msg.Role = internal.RoleDev
		case model.Tool:
			msg.Role = internal.RoleTool
		}

	}
	request := &chatRequest{
		Messages: chatMsg,
	}
	// Create request for intenal client
	response, err := o.client.ChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}
	// Validate response
	if len(response.Choices) == 0 {
		return nil, ErrNoContentInResponse
	}

	gen, err := toModelGeneration(response)
	if err != nil {
		return nil, ErrModelGen
	}
	return gen, nil
}

func toModelGeneration(response internal.ChatCompletionResponse) (*model.Generation, error) {
	gen := &model.Generation{}
	for _, v := range response.Choices {

	}
	return gen, nil
}
