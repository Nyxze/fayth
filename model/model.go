// Package model defines interfaces and data structures for working with LLM message generation.
package model

import (
	"context"
	"iter"
)

// Generation represents a complete or partial language model response.
// It may contain a static list of messages or a lazily-evaluated stream.
// If streamed, messages are appended incrementally to Messages during iteration.
type Generation struct {
	// Messages contains the full list of generated messages.
	// This may be populated immediately or during streaming.
	Messages []Message `json:"messages"`

	// MessageIter is an optional MessageIter of messages, lazily evaluated.
	MessageIter MessageIter

	// err stores any error that occurred during streaming.
	Err error
}

// MessageIter is an alias for an iterator that yields Message values.
type MessageIter = iter.Seq[Message]

// NewGeneration returns a Generation with a static list of messages.
func NewGeneration(m []Message) *Generation {
	return &Generation{
		Messages: m,
	}
}

// NewGenerationWithStream returns a Generation that streams its messages lazily from the given iterator.
func NewGenerationWithStream(stream MessageIter) *Generation {
	return &Generation{
		MessageIter: stream,
	}
}

// All returns an iterator over the generated messages.
// If the generation was streamed, this will lazily consume the stream and append results to Messages.
// Subsequent calls to All will yield from the fully populated Messages slice.
func (g *Generation) All() MessageIter {
	if g.MessageIter != nil {
		return g.iterStream()
	}
	return func(yield func(Message) bool) {
		for _, v := range g.Messages {
			if !yield(v) {
				return
			}
		}
	}
}

// Err returns any error that occurred during message streaming.
// It should be checked after exhausting the iterator returned by All.
func (g *Generation) Error() error {
	return g.Err
}

// iterStream returns a one-time iterator that consumes the underlying stream.
// As each message is received, it is appended to Messages.
// This function is only called once; subsequent calls to All will yield from the populated Messages slice.
func (g *Generation) iterStream() MessageIter {
	return func(yield func(Message) bool) {
		seq := g.MessageIter
		g.MessageIter = nil
		for v := range seq {
			g.Messages = append(g.Messages, v)
			if !yield(v) {
				return
			}
		}
	}
}

// MessageHandler defines a function that processes a single Message.
// This is typically used for handling streamed responses in real-time.
type MessageHandler func(Message)

// Model defines the interface for a language model capable of generating messages.
//
// Generate initiates a generation process based on the provided context, input messages,
// and optional model options. The returned Generation may contain all messages immediately
// or yield them over time if streaming is enabled.
type Model interface {
	Generate(ctx context.Context, m []Message, opts ...ModelOption) (*Generation, error)
}
