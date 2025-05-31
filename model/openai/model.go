package openai

import (
	"context"
	"errors"
	"strings"

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
	for i := range len(chatMsg) {
		chatMsg[i] = toOpenAIMessages(messages[i])
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
		Stream:           options.StreamHandler != nil,
	}

	// Handle non-streaming case
	if !req.Stream {
		resp, err := m.client.Chat.Completion(ctx, req)
		if err != nil {
			return nil, err
		}
		return toGeneration(resp)
	}
	// Handle streaming case
	streamChan, err := m.client.Chat.CompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}

	var currentContent strings.Builder
	var lastRole string

	for {
		select {
		case <-ctx.Done():
			return &model.Generation{Error: ctx.Err()}, nil
		case streamResp, ok := <-streamChan:
			if !ok {
				// Stream completed, return final generation
				if currentContent.Len() == 0 {
					return nil, ErrNoContentInResponse
				}
				return &model.Generation{
					Messages: []model.Message{
						model.NewTextMessage(internal.ToModelRole(lastRole), currentContent.String()),
					},
				}, nil
			}

			// Handle errors in stream response
			if streamResp.Object == "error" {
				err := errors.New(streamResp.Choices[0].Delta.Content)
				return &model.Generation{Error: err}, nil
			}

			// Process each choice in the response
			for _, choice := range streamResp.Choices {
				if choice.FinishReason != "" {
					continue
				}

				// Update role if provided
				if choice.Delta.Role != "" {
					lastRole = choice.Delta.Role
				}

				// Accumulate content if provided
				if choice.Delta.Content != "" {
					currentContent.WriteString(choice.Delta.Content)
					// Create message from current content and call handler
					role := internal.ToModelRole(lastRole)
					msg := model.NewTextMessage(role, currentContent.String())

					// Pass message to handler
					if err := options.StreamHandler(msg); err != nil {
						return &model.Generation{Error: err}, nil
					}
				}
			}
		}
	}
}

func (c *llm) String() string {
	return "OpenAI"
}

func toGeneration(resp *internal.ChatCompletionResponse) (*model.Generation, error) {
	gen := &model.Generation{}
	for _, v := range resp.Choices {
		role := internal.ToModelRole(v.Message.Role)
		msg := model.NewTextMessage(role, v.Message.Content)
		gen.Messages = append(gen.Messages, msg)
	}
	return gen, nil
}

func toOpenAIMessages(message model.Message) ChatMessage {
	return internal.ChatMessage{
		Role:     internal.ToOpenAIRole(message.Role),
		Contents: internal.ToChatContent(message.Contents),
	}
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
