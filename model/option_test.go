package model

import (
	"testing"
)

func TestModelOptions(t *testing.T) {
	t.Run("WithModel", func(t *testing.T) {
		opts := ModelOptions{}
		WithModel("gpt-4")(&opts)
		if opts.Model != "gpt-4" {
			t.Errorf("Expected model 'gpt-4', got '%s'", opts.Model)
		}
	})

	t.Run("WithTemperature", func(t *testing.T) {
		opts := ModelOptions{}
		WithTemperature(0.7)(&opts)
		if opts.Temperature != 0.7 {
			t.Errorf("Expected temperature 0.7, got %f", opts.Temperature)
		}
	})

	t.Run("WithMaxTokens", func(t *testing.T) {
		opts := ModelOptions{}
		WithMaxTokens(1000)(&opts)
		if opts.MaxTokens != 1000 {
			t.Errorf("Expected max_tokens 1000, got %d", opts.MaxTokens)
		}
	})

	t.Run("WithTopP", func(t *testing.T) {
		opts := ModelOptions{}
		WithTopP(0.9)(&opts)
		if opts.TopP != 0.9 {
			t.Errorf("Expected top_p 0.9, got %f", opts.TopP)
		}
	})

	t.Run("WithFrequencyPenalty", func(t *testing.T) {
		opts := ModelOptions{}
		WithFrequencyPenalty(0.5)(&opts)
		if opts.FrequencyPenalty != 0.5 {
			t.Errorf("Expected frequency_penalty 0.5, got %f", opts.FrequencyPenalty)
		}
	})

	t.Run("WithPresencePenalty", func(t *testing.T) {
		opts := ModelOptions{}
		WithPresencePenalty(-0.3)(&opts)
		if opts.PresencePenalty != -0.3 {
			t.Errorf("Expected presence_penalty -0.3, got %f", opts.PresencePenalty)
		}
	})

	t.Run("WithStop", func(t *testing.T) {
		opts := ModelOptions{}
		WithStop("END", "STOP")(&opts)
		if len(opts.Stop) != 2 || opts.Stop[0] != "END" || opts.Stop[1] != "STOP" {
			t.Errorf("Expected stop sequences [END, STOP], got %v", opts.Stop)
		}
	})

	t.Run("WithSeed", func(t *testing.T) {
		opts := ModelOptions{}
		WithSeed(12345)(&opts)
		if opts.Seed != 12345 {
			t.Errorf("Expected seed 12345, got %d", opts.Seed)
		}
	})

	t.Run("WithUser", func(t *testing.T) {
		opts := ModelOptions{}
		WithUser("test-user")(&opts)
		if opts.User != "test-user" {
			t.Errorf("Expected user 'test-user', got '%s'", opts.User)
		}
	})

	t.Run("WithJSONMode", func(t *testing.T) {
		opts := ModelOptions{}
		WithJSONMode()(&opts)
		if opts.ResponseFormat.Type != "json_object" {
			t.Errorf("Expected JSON response format, got %s", opts.ResponseFormat.Type)
		}
	})

	t.Run("WithTextMode", func(t *testing.T) {
		opts := ModelOptions{}
		WithTextMode()(&opts)
		if opts.ResponseFormat.Type != "text" {
			t.Errorf("Expected text response format, got %s", opts.ResponseFormat.Type)
		}
	})

	t.Run("WithStream", func(t *testing.T) {
		opts := ModelOptions{}
		WithStream(true)(&opts)
		if !opts.Stream {
			t.Errorf("Expected stream true, got %t", opts.Stream)
		}
	})

	t.Run("WithLogProbs", func(t *testing.T) {
		opts := ModelOptions{}
		WithLogProbs(true)(&opts)
		if !opts.LogProbs {
			t.Errorf("Expected logprobs true, got %t", opts.LogProbs)
		}
	})

	t.Run("WithTopLogProbs", func(t *testing.T) {
		opts := ModelOptions{}
		WithTopLogProbs(10)(&opts)
		if opts.TopLogProbs != 10 {
			t.Errorf("Expected top_logprobs 10, got %d", opts.TopLogProbs)
		}
	})
}

func TestMergeOptions(t *testing.T) {
	t.Run("MergeBasicOptions", func(t *testing.T) {
		base := ModelOptions{
			Model:       "gpt-3.5-turbo",
			Temperature: 1.0,
		}

		merged := MergeOptions(base,
			WithModel("gpt-4"),
			WithTemperature(0.7),
			WithMaxTokens(1000),
		)

		if merged.Model != "gpt-4" {
			t.Errorf("Expected merged model 'gpt-4', got '%s'", merged.Model)
		}
		if merged.Temperature != 0.7 {
			t.Errorf("Expected merged temperature 0.7, got %f", merged.Temperature)
		}
		if merged.MaxTokens != 1000 {
			t.Errorf("Expected merged max_tokens 1000, got %d", merged.MaxTokens)
		}
	})

	t.Run("MergeComplexOptions", func(t *testing.T) {
		base := ModelOptions{
			Model:       "gpt-4",
			Temperature: 0.5,
		}

		merged := MergeOptions(base,
			WithTopP(0.9),
			WithFrequencyPenalty(0.3),
			WithPresencePenalty(-0.2),
			WithStop("END"),
			WithSeed(42),
			WithUser("test"),
			WithJSONMode(),
			WithLogProbs(true),
			WithTopLogProbs(5),
		)

		// Check all options were applied
		if merged.Model != "gpt-4" {
			t.Errorf("Expected model preserved")
		}
		if merged.Temperature != 0.5 {
			t.Errorf("Expected temperature preserved")
		}
		if merged.TopP != 0.9 {
			t.Errorf("Expected top_p 0.9")
		}
		if merged.FrequencyPenalty != 0.3 {
			t.Errorf("Expected frequency_penalty 0.3")
		}
		if merged.PresencePenalty != -0.2 {
			t.Errorf("Expected presence_penalty -0.2")
		}
		if len(merged.Stop) != 1 || merged.Stop[0] != "END" {
			t.Errorf("Expected stop sequence [END]")
		}
		if merged.Seed != 42 {
			t.Errorf("Expected seed 42")
		}
		if merged.User != "test" {
			t.Errorf("Expected user 'test'")
		}
		if merged.ResponseFormat.Type != "json_object" {
			t.Errorf("Expected JSON response format")
		}
		if !merged.LogProbs {
			t.Errorf("Expected logprobs true")
		}
		if merged.TopLogProbs != 5 {
			t.Errorf("Expected top_logprobs 5")
		}
	})

	t.Run("OverrideOptions", func(t *testing.T) {
		base := ModelOptions{
			Model:       "gpt-3.5-turbo",
			Temperature: 1.0,
			MaxTokens:   500,
		}

		merged := MergeOptions(base,
			WithModel("gpt-4"),
			WithMaxTokens(1000),
		)

		if merged.Model != "gpt-4" {
			t.Errorf("Expected model override to 'gpt-4'")
		}
		if merged.MaxTokens != 1000 {
			t.Errorf("Expected max_tokens override to 1000")
		}
		if merged.Temperature != 1.0 {
			t.Errorf("Expected temperature to remain unchanged")
		}
	})
}