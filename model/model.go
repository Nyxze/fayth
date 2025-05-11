package model

import "context"

type Model interface {
	Generate(ctx context.Context, messages []Message) (Message, error)
}
