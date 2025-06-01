package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"
	"nyxze/choco-go"
	choco_json "nyxze/choco-go/json"
	"nyxze/choco-go/seqio"
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
type ChatResponse struct {
	Response   *ChatCompletionResponse
	StreamIter iter.Seq[ChatCompletionChunk]
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
func (c *ChatService) Completion(ctx context.Context, chatRequest ChatCompletionRequest, opts ...CallOption) (*ChatResponse, error) {

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
	// Stream path
	var chatResponse ChatCompletionResponse
	if !chatRequest.Stream {
		// Non-streaming path
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&chatResponse)
		if err != nil {
			return nil, err
		}
		return &ChatResponse{Response: &chatResponse}, nil
	}
	return &ChatResponse{
		StreamIter: readChunk(ctx, res.Body),
	}, nil
}

func readChunk(ctx context.Context, r io.ReadCloser) iter.Seq[ChatCompletionChunk] {
	return func(yield func(ChatCompletionChunk) bool) {
		for line := range seqio.Lines(ctx, r) {

			// Handle empty input
			if len(line) == 0 {
				continue
			}

			// Remove "data: " prefix that's required in SSE format
			b, found := strings.CutPrefix(line, "data: ")
			if !found {
				continue
			}

			// Check for stream termination message
			if strings.HasPrefix(b, "[DONE]") {
				// Using io.EOF is appropriate here as it signals normal end of stream
				break
			}
			var value ChatCompletionChunk
			// Parse JSON chunk into ChatCompletionChunk struct
			if err := json.Unmarshal([]byte(b), &value); err != nil {
				// Skip failing parsing
				fmt.Println("Failed to parse to ChatCompletion chunk")
				continue
			}
			if !yield(value) {
				return
			}
		}
	}
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
	// Create pipeline with headers and base URL
	funcs := []choco.PipelineStepFunc{
		applyHeaders(config),
	}

	// Apply base URL if provided
	if config.BaseUrl != nil {
		funcs = append(funcs, applyBaseUrl(config.BaseUrl))
	}

	// Create pipeline with custom client if provided
	opts := []choco.PipelineOption{
		choco.WithStepFuncs(funcs...),
	}

	// Add custom transport if client is provided
	if config.HTTPClient != nil {
		opts = append(opts, choco.WithCustomTransport(&customTransport{client: config.HTTPClient}))
	}

	// Create pipeline
	pipeline, err := choco.NewPipeline(opts...)
	if err != nil {
		return nil, err
	}

	// Send request
	return pipeline.Execute(req)
}

// customTransport implements choco.Transport using a custom http.Client
type customTransport struct {
	client *http.Client
}

func (t *customTransport) Send(req *http.Request) (*http.Response, error) {
	return t.client.Do(req)
}

func applyHeaders(config *CallConfig) choco.PipelineStepFunc {
	return func(req *choco.Request, next choco.RequestHandlerFunc) (*http.Response, error) {
		req.SetAuthorization(choco.AuthSchemeBearer, config.APIKey)
		if config.Organization != "" {
			req.Raw().Header.Set("OpenAI-Organization", config.Organization)
		}
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
