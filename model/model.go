package model

import (
	"context"
)

// Generation represents a complete response
type Generation struct {
	Messages []Message `json:"messages"`
	Error    error     `json:"error,omitempty"`
}

// MessageHandler is called for each message chunk during streaming
type MessageHandler func(Message) error

type Model interface {
	// Generate performs a generation with optional streaming
	// If a MessageHandler is provided via WithStream option, it will be called
	// for each message chunk as it arrives. The final complete response will
	// still be returned in the Generation.
	Generate(ctx context.Context, m []Message, opts ...ModelOption) (*Generation, error)
}
