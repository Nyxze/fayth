package openai

const (
	RoleSystem    = "system"
	RoleAssistant = "assistant"
	RoleUser      = "user"
	RoleFunction  = "function"
	RoleTool      = "tool"
)
const (
	API_KEY_ENV    = "OPENAI_API_KEY" //nolint:gosec
	MODEL_NAME_ENV = "OPENAI_MODEL"   //nolint:gosec
)
