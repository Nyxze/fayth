package openai

import "nyxze/fayth/model/openai/internal"

type CallOption func(*CallOptions) error

type CallOptions struct {
	Model        string
	internalOpts []internal.CallOption
}

// WithAPIKey sets the API key to authenticate requests.
func WithAPIKey(key string) CallOption {
	return func(opts *CallOptions) error {
		opts.internalOpts = append(opts.internalOpts, internal.WithAPIKey(key))
		return nil
	}
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(base string) CallOption {
	return func(opts *CallOptions) error {
		opts.internalOpts = append(opts.internalOpts, internal.WithBaseURL(base))
		return nil
	}
}

// WithModel sets the model name to use for generation.
func WithModel(model string) CallOption {
	return func(opts *CallOptions) error {
		opts.Model = model
		return nil
	}
}

// WithOrganization sets the organization ID or name.
func WithOrganization(org string) CallOption {
	return func(opts *CallOptions) error {
		opts.internalOpts = append(opts.internalOpts, internal.WithOrganization(org))
		return nil
	}
}
