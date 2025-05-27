package model

type ModelOptions struct {
	Model       string
	Temperature float64
}
type ModelOption func(*ModelOptions)

func WithModel(model string) ModelOption {
	return func(mo *ModelOptions) {
		mo.Model = model
	}
}
