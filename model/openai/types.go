package openai

import (
	"nyxze/fayth/model"
	"nyxze/fayth/model/openai/internal"
)

var InternalRoleMap = map[internal.Role]model.Role{
	internal.RoleUser:      model.User,
	internal.RoleAssistant: model.Assistant,

	internal.RoleDev:    model.System,
	internal.RoleSystem: model.System,

	internal.RoleFunction: model.Tool,
	internal.RoleTool:     model.Tool,
}

var RoleMap = map[model.Role]internal.Role{}

const (
	API_KEY_ENV    = "OPENAI_API_KEY" //nolint:gosec
	MODEL_NAME_ENV = "OPENAI_MODEL"   //nolint:gosec
)

const (
	defaultChatModel = "gpt-3.5-turbo"
)
