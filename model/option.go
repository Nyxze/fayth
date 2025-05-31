package model

// ModelOptions contains configuration options for model inference
type ModelOptions struct {
	// Model specifies which model to use (e.g., "gpt-4", "gpt-3.5-turbo")
	Model string `json:"model"`
	
	// Temperature controls randomness in the output (0.0 to 2.0)
	// Higher values make output more random, lower values more deterministic
	Temperature float64 `json:"temperature"`
	
	// MaxTokens sets the maximum number of tokens to generate
	// If 0, uses the model's default maximum
	MaxTokens int `json:"max_tokens,omitzero"`
	
	// TopP controls nucleus sampling (0.0 to 1.0)
	// Alternative to temperature for controlling randomness
	TopP float64 `json:"top_p,omitzero"`
	
	// FrequencyPenalty penalizes frequent tokens (-2.0 to 2.0)
	// Positive values decrease likelihood of repeating tokens
	FrequencyPenalty float64 `json:"frequency_penalty,omitzero"`
	
	// PresencePenalty penalizes tokens that have appeared (-2.0 to 2.0)
	// Positive values encourage talking about new topics
	PresencePenalty float64 `json:"presence_penalty,omitzero"`
	
	// Stop sequences where the model will stop generating
	Stop []string `json:"stop,omitzero"`
	
	// Seed for deterministic sampling (if supported by model)
	Seed int64 `json:"seed,omitzero"`
	
	// User identifier for abuse monitoring
	User string `json:"user,omitzero"`
	
	// ResponseFormat specifies the format of the response
	// Can be "text" or "json_object" for JSON mode
	ResponseFormat ResponseFormat `json:"response_format,omitzero"`
	
	// Stream enables streaming responses (not yet implemented)
	Stream bool `json:"stream,omitzero"`
	
	// LogProbs enables log probabilities in response
	LogProbs bool `json:"logprobs,omitzero"`
	
	// TopLogProbs specifies number of top log probabilities to return (0-20)
	TopLogProbs int `json:"top_logprobs,omitzero"`
}

// ResponseFormat specifies the format of the model's output
type ResponseFormat struct {
	Type string `json:"type,omitzero"` // "text" or "json_object"
}

type ModelOption func(*ModelOptions)

// WithModel sets the model to use
func WithModel(model string) ModelOption {
	return func(mo *ModelOptions) {
		mo.Model = model
	}
}

// WithTemperature sets the temperature for randomness control
func WithTemperature(temp float64) ModelOption {
	return func(mo *ModelOptions) {
		mo.Temperature = temp
	}
}

// WithMaxTokens sets the maximum number of tokens to generate
func WithMaxTokens(maxTokens int) ModelOption {
	return func(mo *ModelOptions) {
		mo.MaxTokens = maxTokens
	}
}

// WithTopP sets the nucleus sampling parameter
func WithTopP(topP float64) ModelOption {
	return func(mo *ModelOptions) {
		mo.TopP = topP
	}
}

// WithFrequencyPenalty sets the frequency penalty
func WithFrequencyPenalty(penalty float64) ModelOption {
	return func(mo *ModelOptions) {
		mo.FrequencyPenalty = penalty
	}
}

// WithPresencePenalty sets the presence penalty
func WithPresencePenalty(penalty float64) ModelOption {
	return func(mo *ModelOptions) {
		mo.PresencePenalty = penalty
	}
}

// WithStop sets stop sequences
func WithStop(stop ...string) ModelOption {
	return func(mo *ModelOptions) {
		mo.Stop = stop
	}
}

// WithSeed sets the seed for deterministic sampling
func WithSeed(seed int64) ModelOption {
	return func(mo *ModelOptions) {
		mo.Seed = seed
	}
}

// WithUser sets the user identifier
func WithUser(user string) ModelOption {
	return func(mo *ModelOptions) {
		mo.User = user
	}
}

// WithJSONMode enables JSON response format
func WithJSONMode() ModelOption {
	return func(mo *ModelOptions) {
		mo.ResponseFormat = ResponseFormat{Type: "json_object"}
	}
}

// WithTextMode sets text response format (default)
func WithTextMode() ModelOption {
	return func(mo *ModelOptions) {
		mo.ResponseFormat = ResponseFormat{Type: "text"}
	}
}

// WithStream enables streaming responses
func WithStream(stream bool) ModelOption {
	return func(mo *ModelOptions) {
		mo.Stream = stream
	}
}

// WithLogProbs enables log probabilities
func WithLogProbs(enabled bool) ModelOption {
	return func(mo *ModelOptions) {
		mo.LogProbs = enabled
	}
}

// WithTopLogProbs sets the number of top log probabilities to return
func WithTopLogProbs(count int) ModelOption {
	return func(mo *ModelOptions) {
		mo.TopLogProbs = count
	}
}

func MergeOptions(base ModelOptions, overrides ...ModelOption) ModelOptions {
	opts := base
	for _, o := range overrides {
		o(&opts)
	}
	return opts
}
