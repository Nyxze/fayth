package openai

import (
	"context"
	"fmt"
	"nyxze/fayth/model"
	"nyxze/fayth/model/openai/internal"
	"os"
)

type chatModel struct {
	client *internal.Client
}

// Compile type interface assertion
var _ model.Model = (*chatModel)(nil)

// Return a New OpenAI [model.Model]
func New(configs ...ModelOptions) (model.Model, error) {
	// Default options
	config := &options{
		ApiKey: os.Getenv(API_KEY_ENV),
		Model:  os.Getenv(MODEL_NAME_ENV),
	}

	// Apply overrides
	for _, conf := range configs {
		conf(config)
	}

	// Validate config
	if config.ApiKey == "" {
		return nil, fmt.Errorf("missing API key")
	}
	if config.Model == "" {
		return nil, fmt.Errorf("missing model")
	}
	client, err := internal.NewClient(config.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to creae openai client %v", err)
	}
	return &chatModel{
		client: client,
	}, nil
}

// [model.Model] implementation
func (o chatModel) Generate(ctx context.Context, messages []model.Message) (model.Message, error) {
	// use underlying client for performing the request
	return model.Message{}, nil
}
