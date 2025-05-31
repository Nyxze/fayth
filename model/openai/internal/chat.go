package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nyxze/choco-go"
	choco_json "nyxze/choco-go/json"
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
		apiError := NewErrorFromResponse(res)
		apiError.Request = req.Raw()
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

// CompletionStream handles streaming responses from the OpenAI Chat API
func (c *ChatService) CompletionStream(ctx context.Context, chatRequest ChatCompletionRequest, opts ...CallOption) (<-chan *ChatCompletionStreamResponse, error) {

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
		apiError := NewErrorFromResponse(res)
		apiError.Request = req.Raw()
		return nil, apiError
	}

	responseChan := make(chan *ChatCompletionStreamResponse)

	go func() {
		defer close(responseChan)
		defer res.Body.Close()

		reader := bufio.NewReader(res.Body)
		for {
			// Check if context is cancelled
			select {
			case <-ctx.Done():
				return
			default:
			}

			// Read the next line
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					responseChan <- &ChatCompletionStreamResponse{
						Object: "error",
						Choices: []ChatStreamingChoice{{
							Delta: ChatCompletionMessage{
								Content: fmt.Sprintf("Error reading stream: %v", err),
							},
						}},
					}
				}
				return
			}

			// Skip empty lines
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			// Remove "data: " prefix
			if !bytes.HasPrefix(line, []byte("data: ")) {
				continue
			}
			line = bytes.TrimPrefix(line, []byte("data: "))

			// Check for stream end
			if bytes.Equal(line, []byte("[DONE]")) {
				return
			}

			// Parse the response
			var streamResponse ChatCompletionStreamResponse
			if err := json.Unmarshal(line, &streamResponse); err != nil {
				responseChan <- &ChatCompletionStreamResponse{
					Object: "error",
					Choices: []ChatStreamingChoice{{
						Delta: ChatCompletionMessage{
							Content: fmt.Sprintf("Error parsing stream response: %v", err),
						},
					}},
				}
				continue
			}

			select {
			case <-ctx.Done():
				return
			case responseChan <- &streamResponse:
			}
		}
	}()

	return responseChan, nil
}

func newRequest(ctx context.Context, method string, cReq ChatCompletionRequest) (*choco.Request, error) {
	// Create choco request from ChatCompletionsRequest
	req, err := choco.NewRequest(ctx, method, completionsAPI)
	if err != nil {
		return nil, err
	}
	err = choco_json.MarshalAsJSON(req, cReq)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func sendRequest(req *choco.Request, config *CallConfig) (*http.Response, error) {

	funcs := []choco.PipelineStepFunc{
		applyHeaders(config), // Set Auth
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
