package internal

type Role = string
type ContentType = string
type ResponsesModel = string
type ChatModel = string

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

const (
	ResponsesModelO1Pro                        ResponsesModel = "o1-pro"
	ResponsesModelO1Pro2025_03_19              ResponsesModel = "o1-pro-2025-03-19"
	ResponsesModelComputerUsePreview           ResponsesModel = "computer-use-preview"
	ResponsesModelComputerUsePreview2025_03_11 ResponsesModel = "computer-use-preview-2025-03-11"
)

const (
	ChatModelGPT4_1                           ChatModel = "gpt-4.1"
	ChatModelGPT4_1Mini                       ChatModel = "gpt-4.1-mini"
	ChatModelGPT4_1Nano                       ChatModel = "gpt-4.1-nano"
	ChatModelGPT4_1_2025_04_14                ChatModel = "gpt-4.1-2025-04-14"
	ChatModelGPT4_1Mini2025_04_14             ChatModel = "gpt-4.1-mini-2025-04-14"
	ChatModelGPT4_1Nano2025_04_14             ChatModel = "gpt-4.1-nano-2025-04-14"
	ChatModelO4Mini                           ChatModel = "o4-mini"
	ChatModelO4Mini2025_04_16                 ChatModel = "o4-mini-2025-04-16"
	ChatModelO3                               ChatModel = "o3"
	ChatModelO3_2025_04_16                    ChatModel = "o3-2025-04-16"
	ChatModelO3Mini                           ChatModel = "o3-mini"
	ChatModelO3Mini2025_01_31                 ChatModel = "o3-mini-2025-01-31"
	ChatModelO1                               ChatModel = "o1"
	ChatModelO1_2024_12_17                    ChatModel = "o1-2024-12-17"
	ChatModelO1Preview                        ChatModel = "o1-preview"
	ChatModelO1Preview2024_09_12              ChatModel = "o1-preview-2024-09-12"
	ChatModelO1Mini                           ChatModel = "o1-mini"
	ChatModelO1Mini2024_09_12                 ChatModel = "o1-mini-2024-09-12"
	ChatModelGPT4o                            ChatModel = "gpt-4o"
	ChatModelGPT4o2024_11_20                  ChatModel = "gpt-4o-2024-11-20"
	ChatModelGPT4o2024_08_06                  ChatModel = "gpt-4o-2024-08-06"
	ChatModelGPT4o2024_05_13                  ChatModel = "gpt-4o-2024-05-13"
	ChatModelGPT4oAudioPreview                ChatModel = "gpt-4o-audio-preview"
	ChatModelGPT4oAudioPreview2024_10_01      ChatModel = "gpt-4o-audio-preview-2024-10-01"
	ChatModelGPT4oAudioPreview2024_12_17      ChatModel = "gpt-4o-audio-preview-2024-12-17"
	ChatModelGPT4oMiniAudioPreview            ChatModel = "gpt-4o-mini-audio-preview"
	ChatModelGPT4oMiniAudioPreview2024_12_17  ChatModel = "gpt-4o-mini-audio-preview-2024-12-17"
	ChatModelGPT4oSearchPreview               ChatModel = "gpt-4o-search-preview"
	ChatModelGPT4oMiniSearchPreview           ChatModel = "gpt-4o-mini-search-preview"
	ChatModelGPT4oSearchPreview2025_03_11     ChatModel = "gpt-4o-search-preview-2025-03-11"
	ChatModelGPT4oMiniSearchPreview2025_03_11 ChatModel = "gpt-4o-mini-search-preview-2025-03-11"
	ChatModelChatgpt4oLatest                  ChatModel = "chatgpt-4o-latest"
	ChatModelCodexMiniLatest                  ChatModel = "codex-mini-latest"
	ChatModelGPT4oMini                        ChatModel = "gpt-4o-mini"
	ChatModelGPT4oMini2024_07_18              ChatModel = "gpt-4o-mini-2024-07-18"
	ChatModelGPT4Turbo                        ChatModel = "gpt-4-turbo"
	ChatModelGPT4Turbo2024_04_09              ChatModel = "gpt-4-turbo-2024-04-09"
	ChatModelGPT4_0125Preview                 ChatModel = "gpt-4-0125-preview"
	ChatModelGPT4TurboPreview                 ChatModel = "gpt-4-turbo-preview"
	ChatModelGPT4_1106Preview                 ChatModel = "gpt-4-1106-preview"
	ChatModelGPT4VisionPreview                ChatModel = "gpt-4-vision-preview"
	ChatModelGPT4                             ChatModel = "gpt-4"
	ChatModelGPT4_0314                        ChatModel = "gpt-4-0314"
	ChatModelGPT4_0613                        ChatModel = "gpt-4-0613"
	ChatModelGPT4_32k                         ChatModel = "gpt-4-32k"
	ChatModelGPT4_32k0314                     ChatModel = "gpt-4-32k-0314"
	ChatModelGPT4_32k0613                     ChatModel = "gpt-4-32k-0613"
	ChatModelGPT3_5Turbo                      ChatModel = "gpt-3.5-turbo"
	ChatModelGPT3_5Turbo16k                   ChatModel = "gpt-3.5-turbo-16k"
	ChatModelGPT3_5Turbo0301                  ChatModel = "gpt-3.5-turbo-0301"
	ChatModelGPT3_5Turbo0613                  ChatModel = "gpt-3.5-turbo-0613"
	ChatModelGPT3_5Turbo1106                  ChatModel = "gpt-3.5-turbo-1106"
	ChatModelGPT3_5Turbo0125                  ChatModel = "gpt-3.5-turbo-0125"
	ChatModelGPT3_5Turbo16k0613               ChatModel = "gpt-3.5-turbo-16k-0613"
)
