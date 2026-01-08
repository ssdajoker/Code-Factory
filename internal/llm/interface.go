package llm

import (
	"context"
)

// LLMService defines the interface for all LLM providers
type LLMService interface {
	// Generate text completion
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)
	
	// Stream text completion
	GenerateStream(ctx context.Context, req GenerateRequest) (<-chan GenerateChunk, error)
	
	// List available models
	ListModels(ctx context.Context) ([]Model, error)
	
	// Get provider name
	Provider() string
}

// GenerateRequest represents a request to generate text
type GenerateRequest struct {
	Prompt      string
	System      string
	Temperature float64
	MaxTokens   int
	Stop        []string
	Context     map[string]string
}

// GenerateResponse represents the response from text generation
type GenerateResponse struct {
	Text         string
	TokensUsed   int
	FinishReason string
	Model        string
}

// GenerateChunk represents a chunk of streamed response
type GenerateChunk struct {
	Text  string
	Done  bool
	Error error
}

// Model represents an LLM model
type Model struct {
	Name   string
	Size   int64
	Format string
}
