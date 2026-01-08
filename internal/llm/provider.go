package llm

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNoProvider     = errors.New("no LLM provider available")
	ErrProviderFailed = errors.New("provider request failed")
)

// Provider is the main interface for LLM providers
type Provider interface {
	// Complete generates a completion for the given prompt
	Complete(ctx context.Context, prompt string, opts Options) (string, error)
	// Name returns the provider name
	Name() string
	// Available checks if the provider is available
	Available(ctx context.Context) bool
	// Models returns available models
	Models(ctx context.Context) ([]string, error)
}

// Options configures the completion request
type Options struct {
	Temperature  float64
	MaxTokens    int
	SystemPrompt string
	Model        string
	Stop         []string
}

// DefaultOptions returns sensible defaults
func DefaultOptions() Options {
	return Options{
		Temperature:  0.7,
		MaxTokens:    2048,
		SystemPrompt: "You are a helpful assistant for software development.",
	}
}

// ProviderType represents the type of LLM provider
type ProviderType string

const (
	ProviderOllama    ProviderType = "ollama"
	ProviderOpenAI    ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
)

// Config holds provider configuration
type Config struct {
	Type       ProviderType
	APIKey     string
	BaseURL    string
	Model      string
}

// NewProvider creates a provider based on config
func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Type {
	case ProviderOllama:
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		return NewOllamaProvider(baseURL, cfg.Model), nil
	case ProviderOpenAI:
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("OpenAI API key required")
		}
		return NewOpenAIProvider(cfg.APIKey, cfg.Model), nil
	case ProviderAnthropic:
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("Anthropic API key required")
		}
		return NewAnthropicProvider(cfg.APIKey, cfg.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider type: %s", cfg.Type)
	}
}
