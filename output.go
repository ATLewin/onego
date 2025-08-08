package onego

type APIResponse struct {
	Code   uint16             `json:"code"`
	Output LlmUnifiedResponse `json:"output"`
}

type LlmUnifiedResponse struct {
	Role         *string   `json:"role,omitempty"`
	Content      string    `json:"content"`          // main field
	Usage        *LlmUsage `json:"usage,omitempty"`
	FinishReason *string   `json:"finish_reason,omitempty"`
}

type LlmUsage struct {
	InputTokens  *uint32 `json:"input_tokens,omitempty"`
	OutputTokens *uint32 `json:"output_tokens,omitempty"`
	TotalTokens  *uint32 `json:"total_tokens,omitempty"`
}
