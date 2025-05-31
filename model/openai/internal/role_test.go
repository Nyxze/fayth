package internal

import (
	"nyxze/fayth/model"
	"testing"
)

func TestToModelRole(t *testing.T) {
	tests := []struct {
		name       string
		openAIRole Role
		want       model.Role
	}{
		{
			name:       "Assistant role",
			openAIRole: AssistantRole,
			want:       model.Assistant,
		},
		{
			name:       "System role",
			openAIRole: SystemRole,
			want:       model.System,
		},
		{
			name:       "Developer role maps to System",
			openAIRole: DevRole,
			want:       model.System,
		},
		{
			name:       "User role",
			openAIRole: UserRole,
			want:       model.User,
		},
		{
			name:       "Tool role",
			openAIRole: ToolRole,
			want:       model.Tool,
		},
		{
			name:       "Function role maps to Tool",
			openAIRole: FuncRole,
			want:       model.Tool,
		},
		{
			name:       "Unknown role defaults to Assistant",
			openAIRole: "unknown",
			want:       model.Assistant,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToModelRole(tt.openAIRole); got != tt.want {
				t.Errorf("ToModelRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToOpenAIRole(t *testing.T) {
	tests := []struct {
		name      string
		modelRole model.Role
		want      Role
	}{
		{
			name:      "Assistant role",
			modelRole: model.Assistant,
			want:      AssistantRole,
		},
		{
			name:      "System role",
			modelRole: model.System,
			want:      SystemRole,
		},
		{
			name:      "User role",
			modelRole: model.User,
			want:      UserRole,
		},
		{
			name:      "Tool role",
			modelRole: model.Tool,
			want:      ToolRole,
		},
		{
			name:      "Unknown role defaults to Assistant",
			modelRole: "unknown",
			want:      AssistantRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToOpenAIRole(tt.modelRole); got != tt.want {
				t.Errorf("ToOpenAIRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
