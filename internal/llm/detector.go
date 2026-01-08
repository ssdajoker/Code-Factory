package llm

import (
	"context"
	"time"
)

// DetectionResult holds the result of provider detection
type DetectionResult struct {
	Available    bool
	ProviderType ProviderType
	ProviderName string
	Models       []string
	Message      string
}

// Detector handles auto-detection of available LLM providers
type Detector struct {
	ollamaURL     string
	openAIKey     string
	anthropicKey  string
}

// NewDetector creates a new detector with the given credentials
func NewDetector(ollamaURL, openAIKey, anthropicKey string) *Detector {
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	return &Detector{
		ollamaURL:    ollamaURL,
		openAIKey:    openAIKey,
		anthropicKey: anthropicKey,
	}
}

// Detect checks all providers and returns the best available one
func (d *Detector) Detect(ctx context.Context) DetectionResult {
	// Check Ollama first (local, free)
	if result := d.checkOllama(ctx); result.Available {
		return result
	}

	// Check OpenAI
	if d.openAIKey != "" {
		return DetectionResult{
			Available:    true,
			ProviderType: ProviderOpenAI,
			ProviderName: "OpenAI",
			Models:       []string{"gpt-4", "gpt-4-turbo", "gpt-3.5-turbo"},
			Message:      "OpenAI API key configured",
		}
	}

	// Check Anthropic
	if d.anthropicKey != "" {
		return DetectionResult{
			Available:    true,
			ProviderType: ProviderAnthropic,
			ProviderName: "Anthropic",
			Models:       []string{"claude-3-opus", "claude-3-sonnet", "claude-3-haiku"},
			Message:      "Anthropic API key configured",
		}
	}

	return DetectionResult{
		Available: false,
		Message:   "No LLM provider available. Install Ollama or configure an API key.",
	}
}

func (d *Detector) checkOllama(ctx context.Context) DetectionResult {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	ollama := NewOllamaProvider(d.ollamaURL, "")
	if !ollama.Available(ctx) {
		return DetectionResult{Available: false}
	}

	models, err := ollama.Models(ctx)
	if err != nil || len(models) == 0 {
		return DetectionResult{
			Available:    true,
			ProviderType: ProviderOllama,
			ProviderName: "Ollama",
			Message:      "Ollama running but no models installed. Run: ollama pull llama2",
		}
	}

	return DetectionResult{
		Available:    true,
		ProviderType: ProviderOllama,
		ProviderName: "Ollama",
		Models:       models,
		Message:      "Ollama running with " + models[0],
	}
}

// GetBestProvider returns a configured provider based on detection
func (d *Detector) GetBestProvider(ctx context.Context) (Provider, error) {
	result := d.Detect(ctx)
	if !result.Available {
		return nil, ErrNoProvider
	}

	model := ""
	if len(result.Models) > 0 {
		model = result.Models[0]
	}

	var apiKey string
	switch result.ProviderType {
	case ProviderOpenAI:
		apiKey = d.openAIKey
	case ProviderAnthropic:
		apiKey = d.anthropicKey
	}

	return NewProvider(Config{
		Type:    result.ProviderType,
		APIKey:  apiKey,
		BaseURL: d.ollamaURL,
		Model:   model,
	})
}
