package main

import (
	"encoding/json"
	"fmt"
	"log"

	"nyxze/fayth/model"
)

func main() {
	fmt.Println("=== OpenAI Model Configuration with omitzero Demo ===\n")

	// Example 1: Basic configuration with only required fields
	fmt.Println("1. Basic Configuration (only model and temperature):")
	basicOptions := model.ModelOptions{
		Model:       "gpt-4",
		Temperature: 0.7,
	}
	
	data, err := json.MarshalIndent(basicOptions, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("JSON output:\n%s\n\n", data)

	// Example 2: Configuration with some optional parameters
	fmt.Println("2. Partial Configuration (with some optional parameters):")
	partialOptions := model.MergeOptions(model.ModelOptions{},
		model.WithModel("gpt-4"),
		model.WithTemperature(0.8),
		model.WithMaxTokens(1000),
		model.WithTopP(0.9),
		model.WithUser("demo-user"),
	)
	
	data, err = json.MarshalIndent(partialOptions, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("JSON output:\n%s\n\n", data)

	// Example 3: Full configuration with all parameters
	fmt.Println("3. Full Configuration (all parameters set):")
	fullOptions := model.MergeOptions(model.ModelOptions{},
		model.WithModel("gpt-4"),
		model.WithTemperature(1.2),
		model.WithMaxTokens(2000),
		model.WithTopP(0.95),
		model.WithFrequencyPenalty(0.5),
		model.WithPresencePenalty(-0.2),
		model.WithStop("END", "STOP"),
		model.WithSeed(42),
		model.WithUser("full-demo"),
		model.WithJSONMode(),
		model.WithStream(true),
		model.WithLogProbs(true),
		model.WithTopLogProbs(10),
	)
	
	data, err = json.MarshalIndent(fullOptions, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("JSON output:\n%s\n\n", data)

	// Example 4: Show zero value omission
	fmt.Println("4. Zero Value Demonstration:")
	zeroOptions := model.ModelOptions{
		Model:            "gpt-4",
		Temperature:      0.7,
		MaxTokens:        0,     // Will be omitted
		TopP:             0.0,   // Will be omitted
		FrequencyPenalty: 0.0,   // Will be omitted
		PresencePenalty:  0.0,   // Will be omitted
		Stop:             nil,   // Will be omitted
		Seed:             0,     // Will be omitted
		User:             "",    // Will be omitted
		ResponseFormat:   model.ResponseFormat{}, // Will be omitted (empty Type)
		Stream:           false, // Will be omitted
		LogProbs:         false, // Will be omitted
		TopLogProbs:      0,     // Will be omitted
	}
	
	data, err = json.MarshalIndent(zeroOptions, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("JSON output (notice only non-zero values):\n%s\n\n", data)

	fmt.Println("=== Benefits of omitzero ===")
	fmt.Println("✓ Cleaner JSON output - only relevant parameters are sent")
	fmt.Println("✓ No need for pointer types - zero values are naturally omitted")
	fmt.Println("✓ Simpler API - direct value assignment instead of pointer manipulation")
	fmt.Println("✓ Better performance - smaller JSON payloads")
	fmt.Println("✓ Modern Go idioms - leverages Go 1.24+ features")
}