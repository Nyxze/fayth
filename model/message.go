package model

import (
	"encoding/json"
	"fmt"
)

type Role string

const (
	User      Role = "user"
	Assistant Role = "assistant"
	Tool      Role = "tool"
	System    Role = "system"
)
const (
	TextKind  string = "text"
	ImageKind string = "image"
)

// ContentPart represents a generic content element of a message.
// Each content type (e.g text, image) implements this interface
// Mostly used with type switch assertion
type ContentPart interface {
	// Kind returns the discriminator type of the content part.
	// Used to distinguish between different content types.
	Kind() string
}

// Message is the main data structure used to communicate with a model.
type Message struct {
	Role       Role              `json:"role"`
	Contents   []ContentPart     `json:"contents"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Properties map[string]any    `json:"properties,omitempty"`
}

// Convinient function for creating a new Message
func NewMessage(role Role, contents ...ContentFunc) Message {
	msg := Message{
		Role: role,
	}
	for _, f := range contents {
		f(&msg)
	}
	return msg
}

func NewTextMessage(role Role, texts ...string) Message {
	return NewMessage(role, WithTextContent(texts...))
}

func (m *Message) UnmarshalJSON(b []byte) error {
	var schema struct {
		Role     Role              `json:"role"`
		Contents []json.RawMessage `json:"contents"`
	}
	if err := json.Unmarshal(b, &schema); err != nil {
		return err
	}
	m.Role = schema.Role
	size := len(schema.Contents)
	m.Contents = make([]ContentPart, 0, size)
	for i := range size {
		cp, err := unmarshalContentPart(schema.Contents[i])
		if err != nil {
			return err
		}
		m.Contents = append(m.Contents, cp)
	}
	return nil
}

// Return all content TextContent
func (m Message) Text() []string {
	var texts []string
	for _, c := range m.Contents {
		if t, ok := c.(TextContent); ok {
			texts = append(texts, t.Text)
		}
	}
	return texts
}

// Represent plain text content in a message
type TextContent struct {
	Text string
}

func (tc TextContent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}{
		Type: tc.Kind(),
		Text: tc.Text,
	})
}

// Kind returns the type of content, which is "text" for TextContent.
func (TextContent) Kind() string { return TextKind }

// Represent image content, either using base64 or form an url
type ImageContent struct {
	SourceType string `json:"source_type"`
	MIMEType   string `json:"mime_type"`
	Data       []byte `json:"data"`
}

func (ImageContent) Kind() string { return ImageKind }

func unmarshalContentPart(data []byte) (ContentPart, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}

	// Dispatch given discriminator type
	switch probe.Type {
	case TextKind:
		var t TextContent
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		return t, nil
	case ImageKind:
		var i ImageContent
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}
		return i, nil
	default:
		return nil, fmt.Errorf("unknown content kind: %s", probe.Type)
	}
}
