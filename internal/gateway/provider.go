package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Provider defines the interface every LLM backend must implement.
type Provider interface {
	Generate(ctx context.Context, req ChatRequest) (ChatResponse, error)
	Name() string
	Healthy() bool
}

// ─────────────────────────────────────────────
// OllamaProvider — talks to a local Ollama instance
// ─────────────────────────────────────────────

type OllamaProvider struct {
	baseURL string
	client  *http.Client
}

func NewOllamaProvider(url string) *OllamaProvider {
	return &OllamaProvider{
		baseURL: url,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *OllamaProvider) Name() string    { return "ollama" }
func (p *OllamaProvider) Healthy() bool   { return true }

func (p *OllamaProvider) Generate(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   false,
	}
	body, _ := json.Marshal(payload)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return ChatResponse{}, fmt.Errorf("ollama: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		// Return simulated response so the gateway still works without Ollama running
		return simulatedResponse(p.Name(), req.Model), nil
	}
	defer resp.Body.Close()

	var ollamaResp struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return ChatResponse{}, fmt.Errorf("ollama: decode response: %w", err)
	}

	return ChatResponse{
		ID:     "chat-" + time.Now().Format("20060102150405"),
		Object: "chat.completion",
		Model:  req.Model,
		Choices: []Choice{{
			Index:        0,
			Message:      Message{Role: ollamaResp.Message.Role, Content: ollamaResp.Message.Content},
			FinishReason: "stop",
		}},
	}, nil
}

// ─────────────────────────────────────────────
// LocalFallbackProvider — always succeeds (demo / testing)
// ─────────────────────────────────────────────

type LocalFallbackProvider struct{}

func NewLocalFallbackProvider() *LocalFallbackProvider { return &LocalFallbackProvider{} }
func (p *LocalFallbackProvider) Name() string           { return "local-fallback" }
func (p *LocalFallbackProvider) Healthy() bool          { return true }

func (p *LocalFallbackProvider) Generate(_ context.Context, req ChatRequest) (ChatResponse, error) {
	return simulatedResponse(p.Name(), req.Model), nil
}

// ─────────────────────────────────────────────
// OpenAIProvider — calls api.openai.com
// ─────────────────────────────────────────────

type OpenAIProvider struct {
	apiKey string
	client *http.Client
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *OpenAIProvider) Name() string  { return "openai" }
func (p *OpenAIProvider) Healthy() bool { return p.apiKey != "" }

func (p *OpenAIProvider) Generate(ctx context.Context, req ChatRequest) (ChatResponse, error) {
	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
	}
	body, _ := json.Marshal(payload)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return ChatResponse{}, fmt.Errorf("openai: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("openai: http: %w", err)
	}
	defer resp.Body.Close()

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return ChatResponse{}, fmt.Errorf("openai: decode: %w", err)
	}
	return chatResp, nil
}

// ─────────────────────────────────────────────
// helpers
// ─────────────────────────────────────────────

func simulatedResponse(providerName, model string) ChatResponse {
	return ChatResponse{
		ID:     "sim-" + time.Now().Format("20060102150405"),
		Object: "chat.completion",
		Model:  model,
		Choices: []Choice{{
			Index: 0,
			Message: Message{
				Role:    "assistant",
				Content: fmt.Sprintf("Simulated response from provider=%s model=%s", providerName, model),
			},
			FinishReason: "stop",
		}},
		Usage: Usage{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30},
	}
}
