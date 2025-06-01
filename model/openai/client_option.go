package openai

import (
	"net/http"
	"nyxze/fayth/model"
	"nyxze/fayth/model/openai/internal"
)

type ClientOption func(*clientOptions) error

// Wrapper around  [internal.CallOption] and [model.ModelOption]
type clientOptions struct {
	modelOpts    []model.ModelOption
	internalOpts []internal.CallOption
}

// WithAPIKey sets the API key to authenticate requests.
func WithAPIKey(key string) ClientOption {
	return func(opts *clientOptions) error {
		opts.internalOpts = append(opts.internalOpts, internal.WithAPIKey(key))
		return nil
	}
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(base string) ClientOption {
	return func(opts *clientOptions) error {
		opts.internalOpts = append(opts.internalOpts, internal.WithBaseURL(base))
		return nil
	}
}

// WithModel sets default model to use for generation.
func WithModel(name string) ClientOption {
	return func(opts *clientOptions) error {
		opts.modelOpts = append(opts.modelOpts, model.WithModel(name))
		return nil
	}
}

// WithOrganization sets the organization ID or name.
func WithOrganization(org string) ClientOption {
	return func(opts *clientOptions) error {
		opts.internalOpts = append(opts.internalOpts, internal.WithOrganization(org))
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client for making requests.
// This is primarily used for testing to inject mock transports.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(opts *clientOptions) error {
		opts.internalOpts = append(opts.internalOpts, internal.WithHTTPClient(client))
		return nil
	}
}
