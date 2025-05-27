package internal

import (
	"encoding/json"
	"testing"
)

func TestMultiPartContent(t *testing.T) {
	jsonSample := `
{
    "model": "gpt-4o",
    "messages": [
        {
            "role": "developer",
            "content": "You are a helpful assistant."
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "Hello"
                },
                {
                    "type": "text",
                    "text": "How are you?"
                }
            ]
        }
    ]
}
`
	var msg ChatCompletionRequest
	err := json.Unmarshal([]byte(jsonSample), &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON into ChatCompletionRequest: %v", err)
	}

	// Check model
	if msg.Model != "gpt-4o" {
		t.Errorf("Invalid model name. Got: %v, Expected: %v", msg.Model, "gpt-4.1")
	}

	// Check messages count
	if len(msg.Messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(msg.Messages))
	}

	// Check first message (simple string)
	m0 := msg.Messages[0]
	if m0.Role != RoleDev {
		t.Errorf("Expected role 'developer', got %v", m0.Role)
	}
	if len(m0.Contents) != 1 {
		t.Errorf("Expected 1 content item in first message, got %d", len(m0.Contents))
	} else if m0.Contents[0].Text != "You are a helpful assistant." {
		t.Errorf("Unexpected content text: %v", m0.Contents[0].Text)
	}

	// Check second message (multi content)
	m1 := msg.Messages[1]
	if m1.Role != RoleUser {
		t.Errorf("Expected role 'user', got %v", m1.Role)
	}
	if len(m1.Contents) != 2 {
		t.Errorf("Expected 2 content items in second message, got %d", len(m1.Contents))
	}
	expected := []string{"Hello", "How are you?"}
	for i, content := range m1.Contents {
		if content.Type != TEXT {
			t.Errorf("Expected type 'text' at index %d, got %v", i, content.Type)
		}
		if content.Text != expected[i] {
			t.Errorf("Expected text %q at index %d, got %q", expected[i], i, content.Text)
		}
	}
}
