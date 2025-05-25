package openai

type Option func(*options)

type options struct {
	apikey       string
	model        string
	baseURL      string
	organization string
}

// WithAPIKey sets the API key to authenticate requests.
func WithAPIKey(key string) Option {
	return func(o *options) {
		o.apikey = key
	}
}

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(base string) Option {
	return func(o *options) {
		o.baseURL = base
	}
}

// WithModel sets the model name to use for generation.
func WithModel(model string) Option {
	return func(o *options) {
		o.model = model
	}
}

// WithOrganization sets the organization ID or name.
func WithOrganization(org string) Option {
	return func(o *options) {
		o.organization = org
	}
}
