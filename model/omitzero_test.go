package model

import (
	"encoding/json"
	"testing"
)

func TestOmitZeroSerialization(t *testing.T) {
	t.Run("EmptyOptions", func(t *testing.T) {
		opts := ModelOptions{
			Model:       "gpt-4",
			Temperature: 0.7,
		}
		
		// Convert to JSON to simulate what would be sent to API
		data, err := json.Marshal(opts)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		
		// Parse back to verify omitzero behavior
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		
		// Check that zero values are omitted
		if _, exists := result["max_tokens"]; exists {
			t.Errorf("Expected max_tokens to be omitted when zero")
		}
		if _, exists := result["top_p"]; exists {
			t.Errorf("Expected top_p to be omitted when zero")
		}
		if _, exists := result["frequency_penalty"]; exists {
			t.Errorf("Expected frequency_penalty to be omitted when zero")
		}
		if _, exists := result["presence_penalty"]; exists {
			t.Errorf("Expected presence_penalty to be omitted when zero")
		}
		if _, exists := result["seed"]; exists {
			t.Errorf("Expected seed to be omitted when zero")
		}
		if _, exists := result["user"]; exists {
			t.Errorf("Expected user to be omitted when empty")
		}
		if _, exists := result["response_format"]; exists {
			t.Errorf("Expected response_format to be omitted when empty")
		}
		if _, exists := result["stream"]; exists {
			t.Errorf("Expected stream to be omitted when false")
		}
		if _, exists := result["logprobs"]; exists {
			t.Errorf("Expected logprobs to be omitted when false")
		}
		if _, exists := result["top_logprobs"]; exists {
			t.Errorf("Expected top_logprobs to be omitted when zero")
		}
		
		// Check that non-zero values are included
		if result["model"] != "gpt-4" {
			t.Errorf("Expected model to be included")
		}
		if result["temperature"] != 0.7 {
			t.Errorf("Expected temperature to be included")
		}
	})

	t.Run("SetOptions", func(t *testing.T) {
		opts := MergeOptions(ModelOptions{},
			WithModel("gpt-4"),
			WithTemperature(0.8),
			WithMaxTokens(1000),
			WithTopP(0.9),
			WithFrequencyPenalty(0.5),
			WithPresencePenalty(-0.2),
			WithStop("END"),
			WithSeed(42),
			WithUser("test-user"),
			WithJSONMode(),
			WithStream(true),
			WithLogProbs(true),
			WithTopLogProbs(5),
		)
		
		// Convert to JSON
		data, err := json.Marshal(opts)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		
		// Parse back to verify all values are included
		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		
		// Check that all set values are included
		expectedFields := []string{
			"model", "temperature", "max_tokens", "top_p", 
			"frequency_penalty", "presence_penalty", "stop", 
			"seed", "user", "response_format", "stream", 
			"logprobs", "top_logprobs",
		}
		
		for _, field := range expectedFields {
			if _, exists := result[field]; !exists {
				t.Errorf("Expected field %s to be included", field)
			}
		}
		
		// Verify specific values
		if result["model"] != "gpt-4" {
			t.Errorf("Expected model 'gpt-4', got %v", result["model"])
		}
		if result["max_tokens"] != float64(1000) { // JSON numbers are float64
			t.Errorf("Expected max_tokens 1000, got %v", result["max_tokens"])
		}
		if result["stream"] != true {
			t.Errorf("Expected stream true, got %v", result["stream"])
		}
	})
}