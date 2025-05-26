package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"nyxze/choco-go"
	"strings"
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

	// Apply config
	config := &CallConfig{}
	for i := range opts {
		opts[i](config)
	}
	path := "chat/completions"

	// Create choco request from ChatCompletionsRequest
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

	// Convert API Response to an error
	if res.StatusCode >= 400 {
		fmt.Println("RESPONSE STATUS", res.StatusCode)
		apiError := Error{}
		err = json.NewDecoder(res.Body).Decode(&apiError)
		if err != nil {
			return nil, err
		}
		apiError.Request = req.Raw()
		apiError.Response = res
		apiError.StatusCode = res.StatusCode
		return nil, apiError
	}
	return nil, nil
}
func createPipeline(c *CallConfig) (choco.Pipeline, error) {
	steps := []choco.PipelineStepFunc{}
	if c.BaseUrl != nil {
		steps = append(steps, applyBaseUrl(c.BaseUrl))
	}
	return choco.NewPipeline(choco.WithStepFuncs(steps...))
}

func applyBaseUrl(u *url.URL) choco.PipelineStepFunc {
	return func(req *choco.Request, next choco.RequestHandlerFunc) (*http.Response, error) {
		raw := req.Raw()
		var err error
		raw.URL, err = u.Parse(strings.TrimLeft(raw.URL.String(), "/"))
		if err != nil {
			return nil, err
		}
		return next(req)
	}
}
