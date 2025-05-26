package model

type ModelOptions struct {
	Temperature float64
}
type ModelOption func(*ModelOptions)
