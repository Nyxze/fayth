package openai

import (
	"context"
	"errors"

	"nyxze/fayth/model"
	"nyxze/fayth/model/openai/internal"
)

var (
	ErrNoContentInResponse = errors.New("no content in generation response")
	ErrModelGen            = errors.New("failed to convert to generation type")
	ErrInvalidMimeType     = errors.New("invalid mime type on content")
)

// Type Alias
type chatMessage = internal.ChatMessage
type chatRequest = internal.ChatCompletionRequest

// Use an underlying [openai.Client] for doing inference
type llm struct {
	client *internal.Client
}

// Compile type interface assertion
var _ model.Model = (*llm)(nil)

// Return a New OpenAI [model.Model]
func New(opts ...CallOption) (*llm, error) {

	client, err := newClient(opts...)
	if err != nil {
		return nil, err
	}
	return &llm{
		client: client,
	}, nil
}

func newClient(opts ...CallOption) (*internal.Client, error) {
	// Create config with default values
	options := CallOptions{
		Model: defaultChatModel,
	}

	// Apply
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}
	client := internal.NewClient(options.internalOpts...)
	return &client, nil
}

// Generate implements [model.Model] interface
func (o llm) Generate(ctx context.Context, messages []model.Message, opts ...model.ModelOption) (*model.Generation, error) {

	// Create request
	req := internal.ChatCompletionRequest{}

	resp, err := o.client.Chat.Completion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert to Generation
	return toModelGeneration(resp)
}

func toModelGeneration(response *internal.ChatCompletionResponse) (*model.Generation, error) {
	gen := &model.Generation{}
	return gen, nil
}
