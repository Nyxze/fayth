package model

import "context"

type Model interface {
	// Request model to generate a [Message] given m
	Generate(ctx context.Context, m []Message) (Message, error)
}
