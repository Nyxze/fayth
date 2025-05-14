package internal

type ChatCompletionResponse struct {
	ID                string       `json:"id,omitempty"`
	Object            string       `json:"object,omitempty"`
	Created           int64        `json:"created,omitempty"`
	Model             string       `json:"model,omitempty"`
	Choices           []ChatChoice `json:"choices,omitempty"`
	Usage             ChatUsage    `json:"usage,omitempty"`
	SystemFingerprint string       `json:"system_fingerprint,omitempty"`
	ServiceTier       string       `json:"service_tier,omitempty"`
}

type ChatChoice struct {
	Index        int         `json:"index,omitempty"`
	Message      ChatMessage `json:"message,omitempty"`
	Logprobs     *Logprobs   `json:"logprobs,omitempty"`
	FinishReason string      `json:"finish_reason,omitempty"`
}

type ChatMessage struct {
	Role        string        `json:"role,omitempty"`
	Content     string        `json:"content,omitempty"`
	Refusal     interface{}   `json:"refusal,omitempty"`
	Annotations []interface{} `json:"annotations,omitempty"`
}

type ChatUsage struct {
	PromptTokens            int                     `json:"prompt_tokens,omitempty"`
	CompletionTokens        int                     `json:"completion_tokens,omitempty"`
	TotalTokens             int                     `json:"total_tokens,omitempty"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens,omitempty"`
	AudioTokens  int `json:"audio_tokens,omitempty"`
}

type CompletionTokensDetails struct {
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`
	AudioTokens              int `json:"audio_tokens,omitempty"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"`
}

type Logprobs struct {
}
