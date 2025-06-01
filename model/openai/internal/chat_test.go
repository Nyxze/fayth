package internal

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

func TestChatService_Completion(t *testing.T) {
	tests := []struct {
		name           string
		request        ChatCompletionRequest
		mockResponse   string
		mockStatus     int
		expectedError  bool
		validateOutput func(*testing.T, *ChatResponse)
	}{
		{
			name: "successful completion",
			request: ChatCompletionRequest{
				Model: "gpt-4",
				Messages: []ChatMessage{
					{Role: UserRole, Contents: []ChatContent{{Type: TextContent, Text: "Hello"}}},
				},
			},
			mockStatus: http.StatusOK,
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
			validateOutput: func(t *testing.T, resp *ChatResponse) {
				if resp.Response.ID != "test-id" {
					t.Errorf("Expected ID 'test-id', got %q", resp.Response.ID)
				}
				if len(resp.Response.Choices) != 1 {
					t.Fatalf("Expected 1 choice, got %d", len(resp.Response.Choices))
				}
				if resp.Response.Choices[0].Message.Content != "Hi there!" {
					t.Errorf("Expected content 'Hi there!', got %q", resp.Response.Choices[0].Message.Content)
				}
			},
		},
		{
			name: "API error response",
			request: ChatCompletionRequest{
				Model: "invalid-model",
				Messages: []ChatMessage{
					{Role: UserRole, Contents: []ChatContent{{Type: TextContent, Text: "Hello"}}},
				},
			},
			mockStatus: http.StatusBadRequest,
			mockResponse: `{
				"error": {
					"message": "Invalid model",
					"type": "invalid_request_error",
					"code": "invalid_model"
				}
			}`,
			expectedError: true,
		},
		{
			name: "missing API key",
			request: ChatCompletionRequest{
				Model: "gpt-4",
				Messages: []ChatMessage{
					{Role: UserRole, Contents: []ChatContent{{Type: TextContent, Text: "Hello"}}},
				},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock transport
			transport := &mockTransport{
				response: &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
					Header:     make(http.Header),
				},
			}

			// Create service with mock client
			service := CreateTestChatService(transport)

			// Make request
			resp, err := service.Completion(context.Background(), tt.request)

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

			if tt.validateOutput != nil {
				tt.validateOutput(t, resp)
			}
		})
	}
}

func TestChatService_CompletionStream(t *testing.T) {
	tests := []struct {
		name          string
		request       ChatCompletionRequest
		mockResponse  string
		mockStatus    int
		expectedError bool
		expectedData  []string
	}{
		{
			name: "successful stream",
			request: ChatCompletionRequest{
				Model:  "gpt-4",
				Stream: true,
				Messages: []ChatMessage{
					{Role: UserRole, Contents: []ChatContent{{Type: TextContent, Text: "Count to 3"}}},
				},
			},
			mockStatus: http.StatusOK,
			mockResponse: `data: {"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"role":"assistant","content":"1"},"finish_reason":null}]}

data: {"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":", 2"},"finish_reason":null}]}

data: {"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":", 3"},"finish_reason":"stop"}]}

data: [DONE]
`,
			expectedData: []string{"1", ", 2", ", 3"},
		},
		{
			name: "API error response",
			request: ChatCompletionRequest{
				Model:  "invalid-model",
				Stream: true,
				Messages: []ChatMessage{
					{Role: UserRole, Contents: []ChatContent{{Type: TextContent, Text: "Hello"}}},
				},
			},
			mockStatus: http.StatusBadRequest,
			mockResponse: `{
				"error": {
					"message": "Invalid model",
					"type": "invalid_request_error",
					"code": "invalid_model"
				}
			}`,
			expectedError: true,
		},
		{
			name: "context cancellation",
			request: ChatCompletionRequest{
				Model:  "gpt-4",
				Stream: true,
				Messages: []ChatMessage{
					{Role: UserRole, Contents: []ChatContent{{Type: TextContent, Text: "Long response"}}},
				},
			},
			mockStatus: http.StatusOK,
			mockResponse: `data: {"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"role":"assistant","content":"Part 1"},"finish_reason":null}]}

data: {"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":" Part 2"},"finish_reason":null}]}

data: {"id":"1","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":" Part 3"},"finish_reason":"stop"}]}

data: [DONE]
`,
			expectedData: []string{"Part 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock transport
			transport := &mockTransport{
				response: &http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(strings.NewReader(tt.mockResponse)),
					Header:     make(http.Header),
				},
			}

			// Create service with mock client

			// Create context (with cancellation for the cancellation test)
			ctx := context.Background()
			if tt.name == "context cancellation" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 100*time.Millisecond)
				defer cancel()
			}

			service := CreateTestChatService(transport)
			// Make request
			cmpl, err := service.Completion(ctx, tt.request)

			// Check error cases
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Read from stream
			var received []string
			for response := range cmpl.StreamIter {
				if response.Object == "error" {
					t.Errorf("Received error in stream: %v", response.Choices[0].Delta.Content)
					break
				}
				if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
					received = append(received, response.Choices[0].Delta.Content)
				}
				if tt.name == "context cancellation" && len(received) == 1 {
					break
				}
			}

			// Verify received data
			if len(received) != len(tt.expectedData) {
				t.Errorf("Expected %d chunks, got %d", len(tt.expectedData), len(received))
				return
			}

			for i, chunk := range received {
				if chunk != tt.expectedData[i] {
					t.Errorf("Chunk %d: expected %q, got %q", i, tt.expectedData[i], chunk)
				}
			}
		})
	}
}

// ### Serialization test ###
func TestSimpleMessage_Unmarshal(t *testing.T) {
	input := `{
		"model": "gpt-4o",
		"messages": [{
			"role": "developer",
			"content": "You are a helpful assistant."
		}]
	}`

	var req ChatCompletionRequest
	err := json.Unmarshal([]byte(input), &req)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if req.Model != "gpt-4o" {
		t.Errorf("Expected model 'gpt-4o', got %q", req.Model)
	}

	if len(req.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(req.Messages))
	}

	msg := req.Messages[0]
	if msg.Role != DevRole {
		t.Errorf("Expected role 'developer', got %v", msg.Role)
	}
	if len(msg.Contents) != 1 || msg.Contents[0].Text != "You are a helpful assistant." {
		t.Errorf("Unexpected content: %+v", msg.Contents)
	}
}
func TestMultiPartMessage_Unmarshal(t *testing.T) {
	input := `{
		"model": "gpt-4o",
		"messages": [{
			"role": "user",
			"content": [
				{ "type": "text", "text": "Hello" },
				{ "type": "text", "text": "How are you?" }
			]
		}]
	}`

	var req ChatCompletionRequest
	err := json.Unmarshal([]byte(input), &req)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	msg := req.Messages[0]
	if msg.Role != UserRole {
		t.Errorf("Expected role 'user', got %v", msg.Role)
	}

	if len(msg.Contents) != 2 {
		t.Errorf("Expected 2 content parts, got %d", len(msg.Contents))
	}

	expected := []string{"Hello", "How are you?"}
	for i, content := range msg.Contents {
		if content.Type != TextContent {
			t.Errorf("Expected TEXT at index %d, got %v", i, content.Type)
		}
		if content.Text != expected[i] {
			t.Errorf("Expected text %q at index %d, got %q", expected[i], i, content.Text)
		}
	}
}
func TestEmptyContent_Unmarshal(t *testing.T) {
	input := `{
		"model": "gpt-4o",
		"messages": [{
			"role": "user",
			"content": []
		}]
	}`

	var req ChatCompletionRequest
	err := json.Unmarshal([]byte(input), &req)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(req.Messages) != 1 || len(req.Messages[0].Contents) != 0 {
		t.Errorf("Expected empty content, got %+v", req.Messages[0].Contents)
	}
}
func TestRoundTrip_MarshalUnmarshal(t *testing.T) {
	original := ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []ChatMessage{
			{Role: UserRole, Contents: []ChatContent{{Type: TextContent, Text: "Hi"}}},
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded ChatCompletionRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Model != original.Model || len(decoded.Messages) != 1 || decoded.Messages[0].Role != UserRole {
		t.Errorf("Round trip mismatch: %+v", decoded)
	}
}

func TestMarshal_ChatCompletionRequest(t *testing.T) {
	yolo := ChatMessage{}
	_, _ = json.Marshal(yolo)
	req := ChatCompletionRequest{
		Model: "gpt-4",
		Messages: []ChatMessage{
			{
				Role: UserRole,
				Contents: []ChatContent{
					{Type: TextContent, Text: "Hello"},
					{Type: TextContent, Text: "World"},
				},
			},
			{
				Role: AssistantRole,
				Contents: []ChatContent{
					{Type: TextContent, Text: "Hi there!"},
				},
			},
		},
		Stream:    true,
		MaxTokens: 100,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal into a map to verify the structure
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Verify top-level fields
	if result["model"] != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got %v", result["model"])
	}
	if result["stream"] != true {
		t.Errorf("Expected stream true, got %v", result["stream"])
	}
	if result["max_tokens"].(float64) != 100 {
		t.Errorf("Expected max_tokens 100, got %v", result["max_tokens"])
	}

	// Verify messages structure
	messages, ok := result["messages"].([]interface{})
	if !ok {
		t.Fatal("Messages field not found or not an array")
	}
	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	// Check first message (user with multiple contents)
	firstMsg := messages[0].(map[string]interface{})
	if firstMsg["role"] != "user" {
		t.Errorf("Expected first message role 'user', got %v", firstMsg["role"])
	}
	firstContents := firstMsg["content"].([]interface{})
	if len(firstContents) != 2 {
		t.Errorf("Expected 2 content items in first message, got %d", len(firstContents))
	}

	// Check second message (assistant with single content)
	secondMsg := messages[1].(map[string]interface{})
	if secondMsg["role"] != "assistant" {
		t.Errorf("Expected second message role 'assistant', got %v", secondMsg["role"])
	}
	// Simple message test
	secondContents := secondMsg["content"].(string)
	if secondContents != "Hi there!" {
		t.Errorf("Expected %s content item in second message, got %s", "Hi there!", secondContents)
	}
}

func CreateTestChatService(t *mockTransport) ChatService {
	service := NewChatService(
		WithAPIKey("test-key"),
		WithBaseURL("http://www.example.com"),
		WithHTTPClient(&http.Client{Transport: t}))

	return service
}
