package openai_test

import (
	"nyxze/fayth/model"
	"nyxze/fayth/model/openai"
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
				return openai.New(openai.WithModel("fake-model"))
			},
			ShouldFail: true,
		},
		"Uses API key from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(openai.API_KEY_ENV, "fake-key")
				return openai.New(openai.WithModel("fake-model"))
			},
			ShouldFail: false,
		},
		"Uses model name from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(openai.MODEL_NAME_ENV, "fake-model")
				return openai.New(openai.WithAPIKey("Fake-key"))
			},
			ShouldFail: false,
		},

		"Fails when model name is missing": {
			Func: func(inner *testing.T) (model.Model, error) {
				// Don't set key from env or options
				inner.Helper()
				return openai.New(openai.WithAPIKey("fake-key"))
			},
			ShouldFail: true,
		},
		"Passes when both API key and model name are provided explicitly": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				return openai.New(openai.WithAPIKey("fake-key"),
					openai.WithModel("model"))
			},
			ShouldFail: false,
		},
		"Uses both API key and model name from environment": {
			Func: func(inner *testing.T) (model.Model, error) {
				inner.Helper()
				inner.Setenv(openai.API_KEY_ENV, "fake-key")
				inner.Setenv(openai.MODEL_NAME_ENV, "fake-model")
				return openai.New()
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
