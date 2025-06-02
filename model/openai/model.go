package openai

import (
	"context"
	"errors"
	"fmt"

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
	Model: ChatModelGPT4,
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

// Generate implements the Model interface for both streaming and non-streaming responses
func (m llm) Generate(ctx context.Context, messages []model.Message, opts ...model.ModelOption) (*model.Generation, error) {
	if len(messages) == 0 {
		return nil, errors.New("empty messages")
	}

	options := model.MergeOptions(m.options, opts...)

	// Validate options
	if err := validateOptions(options); err != nil {
		return nil, err
	}

	chatMsg := make([]ChatMessage, len(messages))
	for i, msg := range messages {
		converted := toOpenAIMessages(msg)
		if err := validateChatMessage(converted); err != nil {
			return nil, fmt.Errorf("invalid message at position %d: %w", i, err)
		}
		chatMsg[i] = converted
	}

	// Create request from ModelOptions & Messages
	req := internal.ChatCompletionRequest{
		Messages:         chatMsg,
		Model:            options.Model,
		Temperature:      options.Temperature,
		TopP:             options.TopP,
		MaxTokens:        options.MaxTokens,
		FrequencyPenalty: options.FrequencyPenalty,
		PresencePenalty:  options.PresencePenalty,
		Stop:             options.Stop,
		Seed:             options.Seed,
		User:             options.User,
		ResponseFormat:   internal.ResponseFormat(options.ResponseFormat),
		LogProbs:         options.LogProbs,
		TopLogProbs:      options.TopLogProbs,
		Stream:           options.Stream,
	}

	resp, err := m.client.Chat.Completion(ctx, req)
	if err != nil {
		return nil, err
	}
	if req.Stream {
		stream := toMessageIter(resp, options.MessageHandler...)
		return model.NewGenerationWithStream(stream), nil
	}
	return toGeneration(resp.Response)
}

func (c *llm) String() string {
	return "OpenAI"
}

func toMessageIter(r *internal.ChatResponse, handlers ...model.MessageHandler) model.MessageIter {
	return func(yield func(model.Message) bool) {
		if r.StreamIter == nil {
			fmt.Println("ChatResponse stream is nil")
			return
		}
		for chunk := range r.StreamIter {
			for _, msg := range fromChunk(chunk) {

				// Raise message
				handleMessage(msg, handlers)

				// Forward
				if !yield(msg) {
					return
				}
			}
		}
	}
}

func handleMessage(msg model.Message, handlers []model.MessageHandler) {
	for _, h := range handlers {
		h(msg)
	}
}
func toGeneration(resp *internal.ChatCompletionResponse) (*model.Generation, error) {
	messages := make([]model.Message, 0, len(resp.Choices))
	for _, v := range resp.Choices {
		role := internal.ToModelRole(v.Message.Role)
		msg := model.NewTextMessage(role, v.Message.Content)
		messages = append(messages, msg)
	}
	return model.NewGeneration(messages), nil
}

func toOpenAIMessages(message model.Message) ChatMessage {
	return internal.ChatMessage{
		Role:     internal.ToOpenAIRole(message.Role),
		Contents: internal.ToChatContent(message.Contents),
	}
}

func fromChunk(c internal.ChatCompletionChunk) []model.Message {
	// Get the first choice from the chunk
	if len(c.Choices) == 0 {
		return nil
	}
	messages := make([]model.Message, 0, len(c.Choices))
	for _, c := range c.Choices {
		role := internal.ToModelRole(c.Delta.Role)
		content := c.Delta.Content
		msg := model.NewTextMessage(role, content)
		msg.Index = c.Index
		messages = append(messages, msg)
	}
	return messages
}

func validateChatMessage(msg ChatMessage) error {
	// Validate role
	if msg.Role == "" {
		return errors.New("message role cannot be empty")
	}

	// Validate contents
	if msg.Contents == nil {
		return errors.New("message contents cannot be nil")
	}

	// Additional validation could be added here based on content type
	return nil
}

// validateOptions validates the model options to ensure they're within acceptable ranges
func validateOptions(options model.ModelOptions) error {
	// Required fields
	if options.Model == "" {
		return errors.New("no model provided")
	}

	// Temperature validation (0.0 to 2.0)
	if options.Temperature < 0.0 || options.Temperature > 2.0 {
		return errors.New("temperature must be between 0.0 and 2.0")
	}

	// TopP validation (0.0 to 1.0) - only validate if non-zero
	if options.TopP != 0 && (options.TopP < 0.0 || options.TopP > 1.0) {
		return errors.New("top_p must be between 0.0 and 1.0")
	}

	// MaxTokens validation (must be positive if set)
	if options.MaxTokens != 0 && options.MaxTokens <= 0 {
		return errors.New("max_tokens must be positive")
	}

	// FrequencyPenalty validation (-2.0 to 2.0) - only validate if non-zero
	if options.FrequencyPenalty != 0 && (options.FrequencyPenalty < -2.0 || options.FrequencyPenalty > 2.0) {
		return errors.New("frequency_penalty must be between -2.0 and 2.0")
	}

	// PresencePenalty validation (-2.0 to 2.0) - only validate if non-zero
	if options.PresencePenalty != 0 && (options.PresencePenalty < -2.0 || options.PresencePenalty > 2.0) {
		return errors.New("presence_penalty must be between -2.0 and 2.0")
	}

	// TopLogProbs validation (0 to 20) - only validate if non-zero
	if options.TopLogProbs != 0 && (options.TopLogProbs < 0 || options.TopLogProbs > 20) {
		return errors.New("top_logprobs must be between 0 and 20")
	}

	// ResponseFormat validation - only validate if Type is set
	if options.ResponseFormat.Type != "" {
		if options.ResponseFormat.Type != "text" && options.ResponseFormat.Type != "json_object" {
			return errors.New("response_format type must be 'text' or 'json_object'")
		}
	}

	// Stop sequences validation (max 4 sequences)
	if len(options.Stop) > 4 {
		return errors.New("maximum of 4 stop sequences allowed")
	}

	return nil
}
