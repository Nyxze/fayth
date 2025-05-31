package fake

import (
	"context"
	"errors"
	"strings"
	"time"

	"nyxze/fayth/model"
)

type fakeModel struct {
	Name     string
	Response model.Message
	// ChunkSize controls how many characters to send in each streaming chunk
	ChunkSize int
	// ChunkDelay controls the delay between chunks when streaming
	ChunkDelay time.Duration
}

// Create a fake model that responds with the given response
func NewModel(name string, resp model.Message) model.Model {
	return &fakeModel{
		Name:       name,
		Response:   resp,
		ChunkSize:  10,                     // Default to 10 characters per chunk
		ChunkDelay: 100 * time.Millisecond, // Default to 100ms delay between chunks
	}
}
func (f fakeModel) Generate(ctx context.Context, m []model.Message, opts ...model.ModelOption) (*model.Generation, error) {
	if len(m) == 0 {
		return nil, errors.New("empty message provided")
	}
	if f.Response.Contents == nil {
		return nil, errors.New("no content set in response")
	}

	// Apply options
	options := model.ModelOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	// Handle streaming case
	if options.StreamHandler != nil {
		// Get the text content to stream
		text := f.Response.Text()[0]
		var currentContent strings.Builder

		// Stream the content in chunks
		for i := 0; i < len(text); i += f.ChunkSize {
			select {
			case <-ctx.Done():
				return &model.Generation{Error: ctx.Err()}, nil
			default:
				// Calculate chunk end
				end := i + f.ChunkSize
				if end > len(text) {
					end = len(text)
				}

				// Append new chunk to accumulated content
				currentContent.WriteString(text[i:end])

				// Create message with accumulated content
				msg := model.NewTextMessage(f.Response.Role, currentContent.String())

				// Call the stream handler
				if err := options.StreamHandler(msg); err != nil {
					return &model.Generation{Error: err}, nil
				}

				// Simulate network delay between chunks
				time.Sleep(f.ChunkDelay)
			}
		}

		// Return the complete response
		return &model.Generation{Messages: []model.Message{f.Response}}, nil
	}

	// Non-streaming case: return complete response immediately
	return &model.Generation{Messages: []model.Message{f.Response}}, nil
}
