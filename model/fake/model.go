package fake

import (
	"context"
	"errors"
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
	if options.Stream {
		// Kinda wacky
		gen := &model.Generation{}
		gen.MsgIter = f.fakeIter(ctx, gen)
		return gen, nil
	}

	// Non-streaming case: return complete response immediately
	return model.NewGeneration([]model.Message{f.Response}), nil
}

func (f *fakeModel) fakeIter(ctx context.Context, gen *model.Generation) model.MessageIter {
	return func(yield func(model.Message) bool) {
		text := f.Response.Text()
		for i := 0; i < len(text); i += f.ChunkSize {
			select {
			case <-ctx.Done():
				gen.Err = ctx.Err()
				return
			default:
				end := min(len(text), i+f.ChunkSize)
				chunk := text[i:end]
				msg := model.NewTextMessage(model.Assistant, chunk)
				if !yield(msg) {
					return
				}
				time.Sleep(f.ChunkDelay)
			}
		}
	}
}
