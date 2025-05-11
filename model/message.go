package model

type Role string

const (
	User      Role = "user"
	Assistant Role = "assistant"
	System    Role = "system"
)

// ContentPart represents a generic content element of a message.
// Each content type (e.g text, image) implements this interface
// Mostly used in conjonction with type switch assertion
type ContentPart interface {
	// Kind returns the discriminator type of the content part.
	// Used to distinguish between different content types.
	Kind() string
}

// Message is the main data structure used to communicate with a model.
type Message struct {
	Role       Role
	Contents   []ContentPart
	Metadata   map[string]string
	Properties map[string]any
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

// Convinient function for creating a new Message with TextContent
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

// Represent plain text content in a message
type TextContent struct {
	Text string `json:"text"`
}

// Kind returns the type of content, which is "text" for TextContent.
func (TextContent) Kind() string { return "text" }

// Represent image content, either using base64 or form an url
type ImageContent struct {
	SourceType string `json:"source_type"`
	MIMEType   string `json:"mime_type"`
	Data       []byte `json:"data"`
}

type ContentFunc func(*Message)

func (ImageContent) Kind() string { return "image" }

// Appends new TextContent to the message's contents
func WithTextContent(content ...string) ContentFunc {
	return func(m *Message) {
		for _, c := range content {
			m.Contents = append(m.Contents, TextContent{Text: c})
		}
	}
}
