package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"nyxze/fayth/model"
	"nyxze/fayth/model/openai"
)

func main() {
	// Create a new OpenAI model instance
	llm, err := openai.New(
		openai.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create a message
	messages := []model.Message{
		model.NewTextMessage(model.User, "Tell me a joke about a chicken"),
	}

	// Example 1: Streaming with message handler
	fmt.Println("Streaming response:")
	fmt.Print("Assistant: ")

	gen, err := llm.Generate(context.Background(), messages,
		model.WithStream(func(msg model.Message) error {
			// Print each chunk as it arrives
			fmt.Print(msg.Text()[0])
			return nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gen.Messages)
}
