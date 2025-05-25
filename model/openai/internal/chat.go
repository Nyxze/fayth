package internal

import (
	"context"
	"net/http"
	"nyxze/choco-go"
)

type ChatService struct {
	// CallOption at Service layer
	// See CallOption documentation
	Options []CallOption
}

func NewChatService(opts ...CallOption) ChatService {
	return ChatService{
		Options: opts,
	}
}

// OpenAI ChatCompletion API
// https://platform.openai.com/docs/api-reference/chat
// Endpoint
// https://api.openai.com/v1/chat/completions
func (c *ChatService) Completion(ctx context.Context, chatRequest ChatCompletionRequest, opts ...CallOption) (resp *ChatCompletionResponse, err error) {

	// Append method CallOption at the end
	opts = append(c.Options[:], opts...)
	path := "chat/completions"

	// Apply config
	config := &CallConfig{}
	for i := range opts {
		opts[i](config)
	}

	// Create choco request from
	req, err := choco.NewRequest(ctx, http.MethodPost, path)
	req.SetBody(nil, choco.ContentTypeAppJSON)
	if err != nil {
		return nil, err
	}
	pipeline, err := createPipeline(config)
	if err != nil {
		return nil, err
	}
	res, err := pipeline.Execute(req)
	if err != nil {
		return nil, err
	}
	// Do smthing with resp
	_ = res
	return nil, nil
}

func createPipeline(c *CallConfig) (choco.Pipeline, error) {
	steps := []choco.PipelineOption{}

	return choco.NewPipeline(steps...)
}
