package model

import (
	"encoding/json"
	"testing"
)

func TestUnmarshal_Message(t *testing.T) {

	jsonData := `[
        {
            "role":"user",
            "contents" : [
                {
                    "type":"text",
                    "text":"Hello"
                }
            ]
        }
    ]`
	var messages []Message
	if err := json.Unmarshal([]byte(jsonData), &messages); err != nil {
		t.Fatal("failed to unmarshal")
	}

	if len(messages) != 1 {
		t.Fatalf("Unexpected len of message, expected : %v, go %v", 1, len(messages))
	}
	c := messages[0].Contents

	if len(c) != 1 {
		t.Fatalf("Unexpected len of content, expected : %v, go %v", 1, len(c))
	}
}
