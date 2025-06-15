package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"nyxze/fayth/model"
	"nyxze/fayth/model/openai/internal"
)

// mockRoundTripper implements http.RoundTripper for testing
type mockRoundTripper struct {
	response     *http.Response
	responseFunc func(*http.Request) (*http.Response, error)
	requests     []*http.Request
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.requests = append(m.requests, req)
	if m.responseFunc != nil {
		return m.responseFunc(req)
	}
	return m.response, nil
}

// Helper function to create a mock response
func mockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// Helper function to create a streaming response
func mockStreamResponse(responses []string) *http.Response {
	var buffer bytes.Buffer
	for _, resp := range responses {
		buffer.WriteString("data: ")
		buffer.WriteString(resp)
		buffer.WriteString("\n\n")
	}
	buffer.WriteString("data: [DONE]\n\n")
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&buffer),
		Header:     make(http.Header),
	}
}

func TestNewModel(t *testing.T) {
	type CreateFunc func(inner *testing.T) (model.Model, error)
	tests := map[string]struct {
		Func       CreateFunc
		ShouldFail bool
	}{
		"Uses API key from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(internal.API_KEY_ENV, "fake-key")
				return New()
			},
			ShouldFail: false,
		},
		"Uses model name from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(internal.MODEL_NAME_ENV, "fake-model")
				return New(WithAPIKey("Fake-key"))
			},
			ShouldFail: false,
		},
		"Passes when both API key and model name are provided explicitly": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				return New(WithAPIKey("fake-api"), WithModel("Hello"))
			},
			ShouldFail: false,
		},
		"Uses both API key and model name from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(internal.API_KEY_ENV, "fake-key")
				inner.Setenv(internal.MODEL_NAME_ENV, "fake-model")
				return New()
			},
			ShouldFail: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			model, err := tt.Func(t)
			if tt.ShouldFail {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			}
			if err == nil && model == nil {
				t.Errorf("expected model, got nil")
			}
		})
	}
}
func TestOpenAI_NonStreaming(t *testing.T) {
	tests := []struct {
		name           string
		input          []model.Message
		mockResponse   string
		expectedError  bool
		expectedOutput string
	}{
		{
			name: "successful completion",
			input: []model.Message{
				model.NewTextMessage(model.User, "Hello"),
			},
			mockResponse: `{
				"id": "test-id",
				"object": "chat.completion",
				"created": 1700000000,
				"model": "gpt-4",
				"choices": [
					{
						"index": 0,
						"message": {
							"role": "assistant",
							"content": "Hi there!"
						},
						"finish_reason": "stop"
					}
				]
			}`,
			expectedOutput: "Hi there!",
		},
		{
			name: "API error",
			input: []model.Message{
				model.NewTextMessage(model.User, "Hello"),
			},
			mockResponse: `{
				"error": {
					"message": "Invalid API key",
					"type": "invalid_request_error",
					"code": "invalid_api_key"
				}
			}`,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(internal.API_KEY_ENV, "fake")
			// Create mock transport
			status := http.StatusOK
			if tt.expectedError {
				status = http.StatusUnauthorized
			}
			mock := &mockRoundTripper{
				response: mockResponse(
					status,
					tt.mockResponse,
				),
			}

			// Create mod with mock transport
			mod, err := New(WithHTTPClient(&http.Client{Transport: mock}))
			if err != nil {
				t.Fatalf("Failed to create model: %v", err)
			}

			// Make the request
			resp, err := mod.Generate(context.Background(), tt.input)

			// Check error cases
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			// Check success cases
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			messages := []model.Message{}
			for m := range resp.Messages() {
				messages = append(messages, m)
			}
			if len(messages) != 1 {
				t.Errorf("Expected 1 message, got %d", len(messages))
				return
			}

			if messages[0].Text() != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, messages[0].Text())
			}

			// Verify request
			if len(mock.requests) != 1 {
				t.Errorf("Expected 1 request, got %d", len(mock.requests))
				return
			}

			req := mock.requests[0]
			if req.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", req.Method)
			}

			if !strings.HasSuffix(req.URL.Path, "chat/completions") {
				t.Errorf("Expected path to end with chat/completions, got %s", req.URL.Path)
			}
		})
	}
}

func TestOpenAI_Streaming(t *testing.T) {
	tests := []struct {
		name            string
		input           []model.Message
		mockResponses   []string
		expectedChunks  []string
		expectedError   bool
		cancelAfterRead int // Number of chunks to read before canceling
	}{
		{
			name: "successful streaming",
			input: []model.Message{
				model.NewTextMessage(model.User, "Count to 3"),
			},
			mockResponses: []string{
				`{"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"role":"assistant","content":"1"},"finish_reason":null}]}`,
				`{"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":", 2"},"finish_reason":null}]}`,
				`{"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":", 3"},"finish_reason":"stop"}]}`,
			},
			expectedChunks: []string{"1", ", 2", ", 3"},
		},
		{
			name: "context cancellation",
			input: []model.Message{
				model.NewTextMessage(model.User, "Long response"),
			},
			mockResponses: []string{
				`{"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"role":"assistant","content":"Part 1"},"finish_reason":null}]}`,
				`{"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":" Part 2"},"finish_reason":null}]}`,
				`{"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":" Part 3"},"finish_reason":null}]}`,
			},
			cancelAfterRead: 1,
			expectedError:   true,
			expectedChunks:  []string{"Part 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(internal.API_KEY_ENV, "fake")
			// Create mock transport with streaming response
			mock := &mockRoundTripper{
				response: mockStreamResponse(tt.mockResponses),
			}

			// Create llm with mock transport
			llm, err := New(WithHTTPClient(&http.Client{Transport: mock}))
			if err != nil {
				t.Fatalf("Failed to create model: %v", err)
			}

			// Create context (with cancellation if specified)
			var ctx context.Context
			var cancel context.CancelFunc
			if tt.cancelAfterRead > 0 {
				ctx, cancel = context.WithCancel(context.Background())
			} else {
				ctx = context.Background()
				cancel = func() {}
			}
			defer cancel()

			// Track received chunks
			var chunks []string

			// Make the streaming request
			gen, err := llm.Generate(ctx, tt.input, model.WithStream(true))

			for m := range gen.Messages() {
				fmt.Println(m.Text())
				chunks = append(chunks, m.Text())
				if len(chunks) == tt.cancelAfterRead {
					cancel()
				}
			}
			// Check success cases
			if err != nil && !tt.expectedError {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify chunks
			if tt.cancelAfterRead == 0 && len(chunks) != len(tt.expectedChunks) {
				t.Errorf("Expected %d chunks, got %d", len(tt.expectedChunks), len(chunks))
				return
			}

			// Verify chunk contents
			for i := range chunks {
				if i >= len(tt.expectedChunks) {
					break
				}
				if chunks[i] != tt.expectedChunks[i] {
					t.Errorf("Chunk %d: expected %q, got %q", i, tt.expectedChunks[i], chunks[i])
				}
			}

			// Verify request
			if len(mock.requests) != 1 {
				t.Errorf("Expected 1 request, got %d", len(mock.requests))
				return
			}

			req := mock.requests[0]
			if req.Method != http.MethodPost {
				t.Errorf("Expected POST request, got %s", req.Method)
			}

			// Verify streaming was enabled in request
			var reqBody internal.ChatCompletionRequest
			if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
				t.Errorf("Failed to decode request body: %v", err)
				return
			}
			if !reqBody.Stream {
				t.Error("Expected stream to be true in request")
			}
		})
	}
}

func TestOpenAI_ValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		input         []model.Message
		options       []model.ModelOption
		expectedError string
	}{
		{
			name:          "empty messages",
			input:         []model.Message{},
			expectedError: "empty messages",
		},
		{
			name: "invalid temperature",
			input: []model.Message{
				model.NewTextMessage(model.User, "Hello"),
			},
			options: []model.ModelOption{
				model.WithTemperature(2.5),
			},
			expectedError: "temperature must be between 0.0 and 2.0",
		},
		{
			name: "empty model",
			input: []model.Message{
				model.NewTextMessage(model.User, "Hello"),
			},
			options: []model.ModelOption{
				model.WithModel(""),
			},
			expectedError: "no model provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := New()
			if err != nil {
				t.Fatalf("Failed to create model: %v", err)
			}

			_, err = model.Generate(context.Background(), tt.input, tt.options...)
			if err == nil {
				t.Error("Expected error but got none")
				return
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing %q, got %q", tt.expectedError, err.Error())
			}
		})
	}
}
