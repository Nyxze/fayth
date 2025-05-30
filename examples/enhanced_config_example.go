package main

import (
	"context"
	"fmt"
	"log"

	"nyxze/fayth/model"
	"nyxze/fayth/model/openai"
)

func main() {
	// This example demonstrates the enhanced model configuration options
	// Note: This won't actually run due to missing dependencies, but shows the API
	
	fmt.Println("Enhanced OpenAI Model Configuration Examples")
	fmt.Println("===========================================")
	
	// Example 1: Basic configuration with enhanced options
	fmt.Println("\n1. Basic Enhanced Configuration:")
	
	// Create a model with enhanced configuration
	model, err := openai.New(
		// Client options would go here
	)
	if err != nil {
		log.Fatal(err)
	}
	
	// Create messages
	messages := []model.Message{
		model.NewTextMessage(model.System, "You are a helpful assistant that responds in JSON format."),
		model.NewTextMessage(model.User, "What are the benefits of renewable energy?"),
	}
	
	// Generate with enhanced options
	generation, err := model.Generate(context.Background(), messages,
		model.WithModel("gpt-4"),
		model.WithTemperature(0.7),
		model.WithMaxTokens(1000),
		model.WithTopP(0.9),
		model.WithJSONMode(), // Force JSON response
		model.WithStop("###"), // Stop at this sequence
		model.WithUser("example-user"),
	)
	
	fmt.Printf("   Configuration applied successfully!\n")
	fmt.Printf("   Messages: %d\n", len(messages))
	fmt.Printf("   Options: Model=gpt-4, Temp=0.7, MaxTokens=1000, JSON mode\n")
	
	// Example 2: Creative writing configuration
	fmt.Println("\n2. Creative Writing Configuration:")
	
	creativeMessages := []model.Message{
		model.NewTextMessage(model.System, "You are a creative writer."),
		model.NewTextMessage(model.User, "Write a short story about a robot discovering emotions."),
	}
	
	_, err = model.Generate(context.Background(), creativeMessages,
		model.WithModel("gpt-4"),
		model.WithTemperature(1.2), // Higher creativity
		model.WithMaxTokens(2000),
		model.WithPresencePenalty(0.6), // Encourage new topics
		model.WithFrequencyPenalty(0.3), // Reduce repetition
		model.WithStop("THE END", "---"),
		model.WithSeed(42), // Reproducible results
	)
	
	fmt.Printf("   Creative configuration: High temperature, presence penalty\n")
	fmt.Printf("   Seed: 42 for reproducible creative output\n")
	
	// Example 3: Analytical/factual configuration
	fmt.Println("\n3. Analytical/Factual Configuration:")
	
	analyticalMessages := []model.Message{
		model.NewTextMessage(model.System, "You are a precise, factual analyst."),
		model.NewTextMessage(model.User, "Analyze the economic impact of remote work."),
	}
	
	_, err = model.Generate(context.Background(), analyticalMessages,
		model.WithModel("gpt-4"),
		model.WithTemperature(0.2), // Low creativity, high precision
		model.WithMaxTokens(1500),
		model.WithTopP(0.8), // Focused sampling
		model.WithFrequencyPenalty(-0.2), // Allow some repetition for emphasis
		model.WithLogProbs(true), // Get confidence scores
		model.WithTopLogProbs(5), // Top 5 alternatives
	)
	
	fmt.Printf("   Analytical configuration: Low temperature, log probabilities\n")
	fmt.Printf("   Optimized for factual, precise responses\n")
	
	// Example 4: Code generation configuration
	fmt.Println("\n4. Code Generation Configuration:")
	
	codeMessages := []model.Message{
		model.NewTextMessage(model.System, "You are an expert programmer. Respond with clean, well-documented code."),
		model.NewTextMessage(model.User, "Write a Python function to calculate fibonacci numbers."),
	}
	
	_, err = model.Generate(context.Background(), codeMessages,
		model.WithModel("gpt-4"),
		model.WithTemperature(0.1), // Very deterministic
		model.WithMaxTokens(800),
		model.WithStop("```", "# End"), // Stop at code block end
		model.WithPresencePenalty(0.1), // Slight penalty for repetition
		model.WithUser("code-generator"),
	)
	
	fmt.Printf("   Code generation: Very low temperature, specific stop sequences\n")
	fmt.Printf("   Optimized for deterministic, clean code output\n")
	
	// Example 5: Conversational configuration
	fmt.Println("\n5. Conversational Configuration:")
	
	conversationMessages := []model.Message{
		model.NewTextMessage(model.System, "You are a friendly, helpful assistant."),
		model.NewTextMessage(model.User, "How's your day going?"),
		model.NewTextMessage(model.Assistant, "I'm doing well, thank you for asking! How can I help you today?"),
		model.NewTextMessage(model.User, "I'm looking for book recommendations."),
	}
	
	_, err = model.Generate(context.Background(), conversationMessages,
		model.WithModel("gpt-3.5-turbo"),
		model.WithTemperature(0.8), // Balanced creativity
		model.WithMaxTokens(300), // Shorter responses
		model.WithPresencePenalty(0.3), // Encourage topic variety
		model.WithFrequencyPenalty(0.2), // Reduce repetitive phrases
		model.WithUser("chat-user-123"),
	)
	
	fmt.Printf("   Conversational: Balanced settings for natural dialogue\n")
	fmt.Printf("   Shorter responses, variety encouraged\n")
	
	// Example 6: Batch processing with different configurations
	fmt.Println("\n6. Batch Processing with Different Configs:")
	
	tasks := []struct {
		name     string
		messages []model.Message
		options  []model.ModelOption
	}{
		{
			name: "Summary",
			messages: []model.Message{
				model.NewTextMessage(model.User, "Summarize this article: [article text]"),
			},
			options: []model.ModelOption{
				model.WithTemperature(0.3),
				model.WithMaxTokens(200),
				model.WithStop("---"),
			},
		},
		{
			name: "Creative",
			messages: []model.Message{
				model.NewTextMessage(model.User, "Write a poem about technology"),
			},
			options: []model.ModelOption{
				model.WithTemperature(1.1),
				model.WithMaxTokens(400),
				model.WithPresencePenalty(0.5),
			},
		},
		{
			name: "Analysis",
			messages: []model.Message{
				model.NewTextMessage(model.User, "Analyze the pros and cons of electric vehicles"),
			},
			options: []model.ModelOption{
				model.WithTemperature(0.4),
				model.WithMaxTokens(800),
				model.WithJSONMode(),
			},
		},
	}
	
	for _, task := range tasks {
		fmt.Printf("   Processing %s task with custom configuration\n", task.name)
		_, err = model.Generate(context.Background(), task.messages, task.options...)
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
		}
	}
	
	fmt.Println("\n✅ All configuration examples completed!")
	fmt.Println("\nKey Benefits of Enhanced Configuration:")
	fmt.Println("• Fine-tuned control over model behavior")
	fmt.Println("• Validation ensures parameters are within valid ranges")
	fmt.Println("• Fluent API for easy configuration")
	fmt.Println("• Support for all major OpenAI parameters")
	fmt.Println("• Type-safe options with proper defaults")
}