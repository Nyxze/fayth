package internal_test

import (
	"encoding/json"
	"fmt"
	"nyxze/fayth/model/openai/internal"
	"testing"
)

func TestMultiPartContent(t *testing.T) {

	rawJSON := `{
    "model": "gpt-4.1",
    "messages": [
      {
        "role": "developer",
        "content": "You are a helpful assistant."
      },
      {
        "role": "user",
        "content": "Hello!"
      }
   	]}`
	var msg internal.ChatCompletionRequest
	err := json.Unmarshal([]byte(rawJSON), &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON into ChatMessage: %v", err)
	}
	if msg.Model != "gpt-4.1" {
		t.Fatalf("Invalid model name. Got : %v, Expected: %v", msg.Model, "gpt-4.1")
	}
}

const RawJSON string = `{
	"messages":
		[
			{
				"role": "developer",
				"content": [
					{
						"type":"text",
						"text":"Hello"
					}
				]
			},
			{
				"role": "user",
				"content": "Hello!"
			}
	]}`

func TestUnmarsh(t *testing.T) {
	msg := internal.ChatCompletionRequest{}
	err := json.Unmarshal([]byte(RawJSON), &msg)
	if err != nil {
		fmt.Println(err)
	}
}
