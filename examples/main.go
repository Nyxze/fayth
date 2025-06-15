package examples

import (
	"context"
	"fmt"
	"log"
	"os"

	"nyxze/fayth/model"
	"nyxze/fayth/model/openai"
)

func Streaming() {
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

	gen, err := llm.Generate(context.Background(), messages, model.WithStream(true))
	if err != nil {
		log.Fatal(err)
	}

	for m := range gen.Messages() {
		fmt.Println(m)
	}
	if gen.Error() != nil {
		fmt.Println("Error when generating", gen)
	}
}
