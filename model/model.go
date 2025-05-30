package model

import "context"

type Generation struct {
	Messages []Message `json:"messages"`
}
type Model interface {
	Generate(ctx context.Context, m []Message, opts ...ModelOption) (*Generation, error)
}
