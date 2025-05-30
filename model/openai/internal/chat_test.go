package internal

import (
	"encoding/json"
	"testing"
)

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
