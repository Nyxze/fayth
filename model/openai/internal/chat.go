package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"nyxze/choco-go"
	"strings"
)

const (
	completionsAPI = "chat/completions"
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
func (c *ChatService) Completion(ctx context.Context, chatRequest ChatCompletionRequest, opts ...CallOption) (*ChatCompletionResponse, error) {

	// Append method CallOption at the end
	opts = append(c.Options[:], opts...)

	// Apply config
	config := &CallConfig{}
	for i := range opts {
		opts[i](config)
	}

	if config.APIKey == "" {
		return nil, ErrMissingToken
	}

	req, err := newRequest(ctx, http.MethodPost, chatRequest)
	if err != nil {
		return nil, err
	}
	res, err := sendRequest(req, config)
	if err != nil {
		return nil, err
	}
	// Convert API Response to an error
	if res.StatusCode >= 400 {
		apiError := NewErrorFromResponse(res.Body)
		apiError.Request = req.Raw()
		apiError.Response = res
		apiError.StatusCode = res.StatusCode
		return nil, apiError
	}

	var chatResponse ChatCompletionResponse
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&chatResponse)
	if err != nil {
		return nil, err
	}
	return &chatResponse, nil
}

func newRequest(ctx context.Context, method string, v any) (*choco.Request, error) {

	// Create choco request from ChatCompletionsRequest
	req, err := choco.NewRequest(ctx, http.MethodPost, completionsAPI)
	if err != nil {
		return nil, err
	}
	err = choco.MarshalAsJSON(req, v)
	return req, err
}

func sendRequest(req *choco.Request, config *CallConfig) (*http.Response, error) {

	funcs := []choco.PipelineStepFunc{
		applyHeaders(config), // Set Auth
		applyBaseUrl(config.BaseUrl),
	}
	// Apply Config to Pipeline
	if config.BaseUrl != nil {
		funcs = append(funcs, applyBaseUrl(config.BaseUrl))
	}

	steps := choco.WithStepFuncs(funcs...)
	pipeline, err := choco.NewPipeline(steps)
	if err != nil {
		return nil, err
	}

	res, err := pipeline.Execute(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func applyHeaders(config *CallConfig) choco.PipelineStepFunc {
	return func(req *choco.Request, next choco.RequestHandlerFunc) (*http.Response, error) {
		req.SetAuthorization(choco.AuthSchemeBearer, config.APIKey)
		return next(req)
	}
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
