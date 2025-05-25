package model

import "context"

type Generation struct {
	Results []*Message
}
type Model interface {
	Generate(ctx context.Context, m []Message, opts ...ModelOption) (*Generation, error)
}
