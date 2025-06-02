package fake

import (
	"context"
	"fmt"
	"testing"
	"time"

	"nyxze/fayth/model"
)

func TestFakeModel_NonStreaming(t *testing.T) {
	tests := map[string]struct {
		input      []model.Message
		output     model.Message
		shouldFail bool
	}{
		"Empty slice": {
			input:      []model.Message{},
			shouldFail: true,
		},
		"Simple message": {
			input:  []model.Message{model.NewTextMessage(model.User, "Hello")},
			output: model.NewTextMessage(model.Assistant, "Test response"),
		},
		"Multiple inputs, but still returns fixed response": {
			input: []model.Message{
				model.NewTextMessage(model.User, "Hello"),
				model.NewTextMessage(model.Assistant, "Hi"),
			},
			output: model.NewTextMessage(model.Assistant, "Test response"),
		},
		"Unset response should return error": {
			input:      []model.Message{model.NewTextMessage(model.User, "Hello")},
			shouldFail: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			model := NewModel("FakeModel", tt.output)
			ctx := context.Background()
			resp, err := model.Generate(ctx, tt.input)

			if tt.shouldFail {
				if err == nil {
					t.Error("Generate() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Generate() unexpected error: %v", err)
				return
			}

			if len(resp.Messages) != 1 {
				t.Errorf("Generate() expected 1 message, got %d", len(resp.Messages))
				return
			}

			got := resp.Messages[0]
			if got.Role != tt.output.Role {
				t.Errorf("Generate() wrong role, got %v, want %v", got.Role, tt.output.Role)
			}
			if got.Text() != tt.output.Text() {
				t.Errorf("Generate() wrong text, got %q, want %q", got.Text(), tt.output.Text())
			}
		})
	}
}

func TestFakeModel_Streaming(t *testing.T) {
	tests := map[string]struct {
		input         model.Message
		chunkSize     int
		chunkDelay    time.Duration
		expectChunks  int
		expectError   error
		cancelContext bool
	}{
		"Basic streaming": {
			input:        model.NewTextMessage(model.Assistant, "Hello world"),
			chunkSize:    5,
			chunkDelay:   10 * time.Millisecond,
			expectChunks: 3, // "Hello" " worl" "d"
		},
		"Single chunk": {
			input:        model.NewTextMessage(model.Assistant, "Hi"),
			chunkSize:    10,
			chunkDelay:   10 * time.Millisecond,
			expectChunks: 1,
		},
		"Context cancellation": {
			input:         model.NewTextMessage(model.Assistant, "This is a long message that should be interrupted"),
			chunkSize:     5,
			chunkDelay:    50 * time.Millisecond,
			cancelContext: true,
			expectError:   context.Canceled,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create model with test configuration
			fakeModel := NewModel("FakeModel", tt.input).(*fakeModel)
			WithChunkSize(tt.chunkSize)(fakeModel)
			WithChunkDelay(tt.chunkDelay)(fakeModel)

			// Create context
			ctx := context.Background()
			if tt.cancelContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				// Cancel after first chunk
				go func() {
					time.Sleep(tt.chunkDelay)
					cancel()
				}()
			}

			// Run streaming generation
			gen, err := fakeModel.Generate(ctx, []model.Message{tt.input}, model.WithStream(true))

			finalMsg := ""
			for m := range gen.All() {
				fmt.Println(m)
				finalMsg += m.Text()
			}
			// Check generation result
			if tt.expectError != nil {
				// For cancelContext, error is at Generation layer
				if tt.cancelContext {
					err = gen.Err
				}
				if err == nil {
					t.Error("Generate() expected error, got nil")
				} else if tt.expectError.Error() != err.Error() {
					t.Errorf("Generate() wrong error, got %v, want %v", err.Error(), tt.expectError)
				}
				return
			}

			// Verify chunks
			if tt.expectChunks > 0 && len(gen.Messages) != tt.expectChunks {
				t.Errorf("Generate() wrong number of chunks, got %d, want %d", len(gen.Messages), tt.expectChunks)
			}

			// Verify final message matches input
			if !tt.cancelContext {
				if finalMsg != tt.input.Text() {
					t.Errorf("Generate() wrong accumulated text, got %q, want %q", finalMsg, tt.input.Text())
				}
			}
		})
	}
}

func TestFakeModel_Configuration(t *testing.T) {
	fakeModel := NewModel("FakeModel", model.NewTextMessage(model.Assistant, "Test")).(*fakeModel)

	t.Run("Default configuration", func(t *testing.T) {
		if fakeModel.ChunkSize != 10 {
			t.Errorf("Default chunk size should be 10, got %d", fakeModel.ChunkSize)
		}
		if fakeModel.ChunkDelay != 100*time.Millisecond {
			t.Errorf("Default chunk delay should be 100ms, got %v", fakeModel.ChunkDelay)
		}
	})

	t.Run("Custom configuration", func(t *testing.T) {
		WithChunkSize(20)(fakeModel)
		WithChunkDelay(50 * time.Millisecond)(fakeModel)

		if fakeModel.ChunkSize != 20 {
			t.Errorf("Chunk size should be 20, got %d", fakeModel.ChunkSize)
		}
		if fakeModel.ChunkDelay != 50*time.Millisecond {
			t.Errorf("Chunk delay should be 50ms, got %v", fakeModel.ChunkDelay)
		}
	})
}
