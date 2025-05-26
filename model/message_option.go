package model

type ContentFunc func(*Message)

// Appends new TextContent to the message's contents
func WithTextContent(content ...string) ContentFunc {
	return func(m *Message) {
		for _, c := range content {
			m.Contents = append(m.Contents, TextContent{Text: c})
		}
	}
}
