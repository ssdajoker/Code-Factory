package llm

import "context"

// Provider interface for LLM integrations
type Provider interface {
	Name() string
	IsAvailable() bool
	Complete(ctx context.Context, prompt string) (string, error)
}

// TODO: Implement Ollama, OpenAI, Anthropic, Azure, Gemini providers
