package openai

import (
	"context"
	"fmt"
	"nyxze/fayth/model"
	"nyxze/fayth/model/openai/internal"
	"testing"
)

func TestNewModel(t *testing.T) {
	type CreateFunc func(inner *testing.T) (model.Model, error)
	tests := map[string]struct {
		Func       CreateFunc
		ShouldFail bool
	}{
		"Fails when API key is missing": {
			// Don't set model from env or options
			Func: func(innerT *testing.T) (model.Model, error) {
				innerT.Helper()
				return New()
			},
			ShouldFail: true,
		},
		"Uses API key from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(internal.API_KEY_ENV, "fake-key")
				return New()
			},
			ShouldFail: false,
		},
		"Uses model name from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(internal.MODEL_NAME_ENV, "fake-model")
				return New(WithAPIKey("Fake-key"))
			},
			ShouldFail: false,
		},
		"Passes when both API key and model name are provided explicitly": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				return New(WithAPIKey("fake-api"), WithModel("Hello"))
			},
			ShouldFail: false,
		},
		"Uses both API key and model name from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(internal.API_KEY_ENV, "fake-key")
				inner.Setenv(internal.MODEL_NAME_ENV, "fake-model")
				return New()
			},
			ShouldFail: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			model, err := tt.Func(t)
			if tt.ShouldFail {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			}
			if err == nil && model == nil {
				t.Errorf("expected model, got nil")
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	t.Setenv(internal.MODEL_NAME_ENV, ChatModelGPT4)
	client, err := New(WithAPIKey("hello"))
	if err != nil {
		panic(1)
	}
	ctx := context.Background()
	msg := model.NewTextMessage(model.User, "Hello")
	resp, err := client.Generate(ctx, []model.Message{msg})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
