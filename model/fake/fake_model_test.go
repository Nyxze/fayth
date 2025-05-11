package fake_test

import (
	"context"
	"nyxze/fayth/model"
	"nyxze/fayth/model/fake"
	"reflect"
	"testing"
)

func TestFakeModel(t *testing.T) {
	tests := map[string]struct {
		Input      []model.Message
		Output     model.Message
		ShouldFail bool
	}{
		"Empty slice": {
			Input:      []model.Message{},
			ShouldFail: true,
		},
		"Simple message": {
			Input:  []model.Message{model.NewTextMessage(model.User, "Hello")},
			Output: model.NewTextMessage(model.Assistant, "Dummy"),
		},
		"Multiple inputs, but still returns fixed response": {
			Input:  []model.Message{model.NewTextMessage(model.User, "Hello")},
			Output: model.NewTextMessage(model.Assistant, "Dummy"),
		},
		"Unset response should return error": {
			Input:      []model.Message{model.NewTextMessage(model.User, "Hello")},
			ShouldFail: true,
		},
	}
	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			model := fake.NewModel("FakeModel", tt.Output)
			ctx := context.Background()
			resp, err := model.Generate(ctx, tt.Input)
			if err != nil {
				if !tt.ShouldFail {
					t.Errorf("Generate() unexpected error = %v", err)
				}
				return // skip output check
			}
			if tt.ShouldFail {
				t.Errorf("Generate() expected failure but got success")
			}

			// Compare value
			if !reflect.DeepEqual(resp, tt.Output) {
				t.Errorf("Generate() error = %v, Expected = %v, Got  %v ", err, tt.Output, resp)
			}
		})
	}
}
