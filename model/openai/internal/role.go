package internal

import "nyxze/fayth/model"

// ToModelRole converts an OpenAI role to a model.Role
func ToModelRole(role Role) model.Role {
	switch role {
	case AssistantRole:
		return model.Assistant
	case SystemRole, DevRole:
		return model.System
	case UserRole:
		return model.User
	case ToolRole, FuncRole:
		return model.Tool
	default:
		return model.Assistant // Default to assistant for unknown roles
	}
}

// ToOpenAIRole converts a model.Role to an OpenAI role
func ToOpenAIRole(role model.Role) Role {
	switch role {
	case model.Assistant:
		return AssistantRole
	case model.System:
		return SystemRole
	case model.User:
		return UserRole
	case model.Tool:
		return ToolRole
	default:
		return AssistantRole // Default to assistant for unknown roles
	}
}
