package onego

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ==== Model Enum Equivalent ====

type Model string

const (
	// ==== OpenAI ====
	ModelGpt4_1            Model = "GPT-4.1"
	ModelGpt4_1Mini        Model = "GPT-4.1-Mini"
	ModelGpt4_1Nano        Model = "GPT-4.1-Nano"
	ModelGptO3             Model = "GPT-o3"
	ModelGptO4Mini         Model = "GPT-o4-mini"
	ModelGptO3Pro          Model = "GPT-o3-pro"
	ModelGpt4o             Model = "GPT-4o"
	ModelGpt4oMini         Model = "GPT-4o-mini"
	ModelGptO1             Model = "GPT-o1"
	ModelGptO3DeepResearch Model = "GPT-o3-DeepResearch"
	ModelGptO3Mini         Model = "GPT-o3-Mini"
	ModelGptO1Mini         Model = "GPT-o1-Mini"

	// ==== Anthropic ====
	ModelClaudeOpus4     Model = "Opus-4"
	ModelClaudeSonnet4   Model = "Sonnet-4"
	ModelClaudeHaiku3_5  Model = "Haiku-3.5"
	ModelClaudeOpus3     Model = "Opus-3"
	ModelClaudeSonnet3_7 Model = "Sonnet-3.7"
	ModelClaudeHaiku3    Model = "Haiku-3"

	// ==== DeepSeek ====
	ModelDeepSeekR1 Model = "DeepSeek-Reasoner"
	ModelDeepSeekV3 Model = "DeepSeek-Chat"

	// ==== Gemini (Google) ====
	ModelGemini25FlashPreview Model = "gemini-2.5-Flash-preview"
	ModelGemini25ProPreview   Model = "gemini-2.5-Pro-preview"
	ModelGemini20Flash        Model = "gemini-2.0-Flash"
	ModelGemini20FlashLite    Model = "gemini-2.0-Flash-lite"
	ModelGemini15Flash        Model = "gemini-1.5-Flash"
	ModelGemini15Flash8B      Model = "gemini-1.5-Flash-8B"
	ModelGemini15Pro          Model = "gemini-1.5-Pro"

	// ==== Mistral ====
	ModelMistralMedium3  Model = "Mistral-Medium-3"
	ModelMagistralMedium Model = "Magistral-Medium"
	ModelCodestral       Model = "Codestral"
	ModelDevstralMedium  Model = "Devstral-Medium"
	ModelMistralLarge    Model = "Mistral-Large"
	ModelPixtralLarge    Model = "Pixtral-Large"
	ModelMinistral8B2410 Model = "Ministral-8B-24.10"
	ModelMinistral3B2410 Model = "Ministral-3B-24.10"
	ModelMistralSmall3_2 Model = "Mistral-Small-3.2"
	ModelMagistralSmall  Model = "Magistral-Small"
	ModelDevstralSmall   Model = "Devstral-Small"
	ModelPixtral12B      Model = "Pixtral-12B"
	ModelMistralNemo     Model = "Mistral-NeMo"
)

var ModelPrices = map[Model]uint32{
	ModelGpt4_1:            1040,
	ModelGpt4_1Mini:        208,
	ModelGpt4_1Nano:        52,
	ModelGptO3:             1040,
	ModelGptO4Mini:         572,
	ModelGptO3Pro:          10400,
	ModelGpt4o:             1300,
	ModelGpt4oMini:         78,
	ModelGptO1:             7800,
	ModelGptO3DeepResearch: 5200,
	ModelGptO3Mini:         572,
	ModelGptO1Mini:         572,

	ModelClaudeOpus4:     9360,
	ModelClaudeSonnet4:   1872,
	ModelClaudeHaiku3_5:  499,
	ModelClaudeOpus3:     9360,
	ModelClaudeSonnet3_7: 1872,
	ModelClaudeHaiku3:    182,

	ModelDeepSeekR1: 142,
	ModelDeepSeekV3: 242,

	ModelGemini25FlashPreview: 380,
	ModelGemini25ProPreview:   1820,
	ModelGemini20Flash:        52,
	ModelGemini20FlashLite:    39,
	ModelGemini15Flash:        78,
	ModelGemini15Flash8B:      39,
	ModelGemini15Pro:          1300,

	ModelMistralMedium3:  2496,
	ModelMagistralMedium: 7280,
	ModelCodestral:       1248,
	ModelDevstralMedium:  2496,
	ModelMistralLarge:    8320,
	ModelPixtralLarge:    8320,
	ModelMinistral8B2410: 208,
	ModelMinistral3B2410: 83,
	ModelMistralSmall3_2: 416,
	ModelMagistralSmall:  2080,
	ModelDevstralSmall:   416,
	ModelPixtral12B:      312,
	ModelMistralNemo:     312,
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Function struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GenerationConfig struct {
	Temperature     float64  `json:"temperature"`
	TopP            float64  `json:"top_p"`
	TopK            uint32   `json:"top_k"`
	CandidateCount  uint32   `json:"candidate_count"`
	MaxOutputTokens uint32   `json:"max_output_tokens"`
	StopSequences   []string `json:"stop_sequences"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type APIInput struct {
	Endpoint      string    `json:"endpoint"`
	Model         Model     `json:"model"`
	Temperature   *float64  `json:"temperature,omitempty"`
	Stream        *bool     `json:"stream,omitempty"`
	Messages      []Message `json:"messages"`
	MaxTokens     uint32    `json:"max_tokens"`
	TopP          float64   `json:"top_p"`
	StopSequences *[]string `json:"stop_sequences,omitempty"`
	Tools         *[]Tool   `json:"tools,omitempty"`

	Contents         *[]Content        `json:"contents,omitempty"`
	SafetySettings   *[]SafetySetting  `json:"safety_settings,omitempty"`
	GenerationConfig *GenerationConfig `json:"generation_config,omitempty"`

	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`

	N              *uint32         `json:"n,omitempty"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
	Seed           *uint32         `json:"seed,omitempty"`
	ToolChoice     *string         `json:"tool_choice,omitempty"`
	User           *string         `json:"user,omitempty"`

	Logprobs    *bool   `json:"logprobs,omitempty"`
	TopLogprobs *uint32 `json:"top_logprobs,omitempty"`

	System *string `json:"system,omitempty"`
	TopK   *uint32 `json:"top_k,omitempty"`
}

func NewAPIInput(endpoint string, model Model, messages []Message, maxTokens uint32) *APIInput {
	return &APIInput{
		Endpoint:         endpoint,
		Model:            model,
		Messages:         messages,
		MaxTokens:        maxTokens,
		Temperature:      floatPtr(1.0),
		Stream:           boolPtr(false),
		TopP:             1.0,
		StopSequences:    nil,
		Tools:            nil,
		Contents:         nil,
		SafetySettings:   nil,
		GenerationConfig: nil,
		FrequencyPenalty: nil,
		PresencePenalty:  nil,
		N:                nil,
		ResponseFormat:   nil,
		Seed:             nil,
		ToolChoice:       nil,
		User:             nil,
		Logprobs:         nil,
		TopLogprobs:      nil,
		System:           nil,
		TopK:             nil,
	}
}

func (a *APIInput) SetTemperature(temp float64) {
	a.Temperature = &temp
}

func (a *APIInput) SetStopSequences(seqs []string) {
	a.StopSequences = &seqs
}

// Send method equivalent to Rust's async send()
func (a *APIInput) Send(apiKey string) (*APIResponse, error) {
	// Marshal APIInput to JSON
	body, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal APIInput: %w", err)
	}

	// HTTP client
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", "https://onellm.dev/api", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Decode response JSON into ApiResponse
	var output APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &output, nil
}

// Helper functions for pointers
func floatPtr(v float64) *float64 { return &v }
func boolPtr(v bool) *bool        { return &v }
