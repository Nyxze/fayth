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

const RawJSON string = `
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
	]`

type message struct {
	Role    string
	Content []Content
}
type Content interface {
	Kind() string
}
type TextContent struct {
	Text string
}

func (TextContent) Kind() string { return "text" }

func (m *message) UnmarshalJSON(data []byte) error {
	// // JSON Schema
	type content struct {
		Type     string
		Text     string
		ImageURL struct {
			Url    string
			Detail string
		}
		InputAudio struct {
			Data   string
			Format string
		}
		File struct {
			FileData string
			Id       string
			Name     string
		}
	}
	ret := struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{}

	// Simple case
	err := json.Unmarshal(data, &ret)
	if err == nil {
		*m = message{Role: ret.Role, Content: []Content{TextContent{Text: ret.Content}}}
		return nil
	}
	ret2 := struct {
		Role    string
		Content []content
	}{}
	err = json.Unmarshal(data, &ret2)
	if err != nil {
		return err
	}
	size := len(ret2.Content)
	contents := make([]Content, size)
	for i := 0; i < size; i++ {
		b := ret2.Content[i]
		switch b.Type {
		case "text":
		case "input_audio":
		case "file":
		}
	}
	return nil
}

func BenchmarkJSON(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		data := []message{}
		err := json.Unmarshal([]byte(RawJSON), &data)
		if err != nil {
			fmt.Println("Error")
		}
	}
}
