package openai

import (
	"context"
	"errors"

	"nyxze/fayth/model"
	"nyxze/fayth/model/openai/internal"
)

// Errors
var (
	ErrNoContentInResponse = errors.New("no content in generation response")
	ErrModelGen            = errors.New("failed to convert to generation type")
	ErrInvalidMimeType     = errors.New("invalid mime type on content")
)

// Default options
var DEFAULT_OPTIONS = model.ModelOptions{
	Model:       ChatModelGPT4,
	Temperature: 1,
}

// Type Alias
type ChatMessage = internal.ChatMessage
type ChatRequest = internal.ChatCompletionRequest

type llm struct {

	// Underlying [internal.Client] for inference call
	client *internal.Client

	// Global options
	options model.ModelOptions
}

// Compile type interface assertion
var _ model.Model = (*llm)(nil)

// Return a New OpenAI [model.Model]
func New(opts ...ClientOption) (*llm, error) {

	options := clientOptions{}

	// Apply options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}
	client := internal.NewClient(options.internalOpts...)

	model := &llm{
		client:  &client,
		options: DEFAULT_OPTIONS,
	}

	for _, opt := range options.modelOpts {
		opt(&model.options)
	}
	return model, nil
}

// Generate implements [model.Model] interface
func (m llm) Generate(ctx context.Context, messages []model.Message, opts ...model.ModelOption) (*model.Generation, error) {

	if len(messages) == 0 {
		return nil, errors.New("empty messages")
	}

	options := model.MergeOptions(m.options, opts...)

	// Validate
	if options.Model == "" {
		return nil, errors.New("no model provided")
	}
	chatMsg := make([]ChatMessage, len(messages))

	for i := range len(chatMsg) {
		chatMsg[i] = toOpenAIMessages(messages[i])
	}
	// Create request from ModelOptions & Messages
	req := internal.ChatCompletionRequest{
		Temperature: options.Temperature,
		Model:       options.Model,
		Messages:    chatMsg,
	}

	resp, err := m.client.Chat.Completion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert to Generation
	return toModelGeneration(resp)
}

func (c *llm) String() string {
	return "OpenAI"
}

func toModelGeneration(response *internal.ChatCompletionResponse) (*model.Generation, error) {
	gen := &model.Generation{}
	return gen, nil
}
func toOpenAIMessages(message model.Message) ChatMessage {
	var chatRole internal.Role
	// Set roles
	switch message.Role {
	case model.User:
		chatRole = internal.UserRole
	case model.Assistant:
		chatRole = internal.AssistantRole
	case model.System:
		chatRole = internal.DevRole
	case model.Tool:
		chatRole = internal.ToolRole
	}
	// Add Name ?
	return internal.ChatMessage{
		Role:     chatRole,
		Contents: internal.ToChatContent(message.Contents),
	}
}
