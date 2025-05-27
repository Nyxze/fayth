package internal

type Role = string
type ContentType = string

const (
	API_ENDPOINT = "https://api.openai.com/v1/"
)
const (
	RoleSystem    Role = "system"
	RoleDev       Role = "developer"
	RoleAssistant Role = "assistant"
	RoleUser      Role = "user"
	RoleFunction  Role = "function"
	RoleTool      Role = "tool"
)

const (
	API_KEY_ENV    = "OPENAI_API_KEY" //nolint:gosec
	MODEL_NAME_ENV = "OPENAI_MODEL"   //nolint:gosec
	BASE_URL_ENV   = "OPENAI_BASE_URL"
	ORG_ID_ENV     = "OPENAI_ORG_ID"
	PROJECT_ID_ENV = "OPENAI_PROJECT_ID"
)

const (
	TEXT        ContentType = "text"
	INPUT_AUDIO ContentType = "input_audio"
	REFUSLA     ContentType = "refusal"
)
