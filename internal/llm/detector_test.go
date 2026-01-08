package llm

import (
	"context"
	"testing"
	"time"
)

func TestNewDetector(t *testing.T) {
	tests := []struct {
		name         string
		ollamaURL    string
		wantOllamaURL string
	}{
		{"default URL", "", "http://localhost:11434"},
		{"custom URL", "http://custom:8080", "http://custom:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDetector(tt.ollamaURL, "", "")
			if d.ollamaURL != tt.wantOllamaURL {
				t.Errorf("ollamaURL = %v, want %v", d.ollamaURL, tt.wantOllamaURL)
			}
		})
	}
}

func TestDetectWithAPIKeys(t *testing.T) {
	tests := []struct {
		name         string
		openAIKey    string
		anthropicKey string
		wantProvider ProviderType
		wantAvail    bool
	}{
		{
			name:         "OpenAI key configured",
			openAIKey:    "sk-test",
			wantProvider: ProviderOpenAI,
			wantAvail:    true,
		},
		{
			name:         "Anthropic key configured",
			anthropicKey: "sk-ant-test",
			wantProvider: ProviderAnthropic,
			wantAvail:    true,
		},
		{
			name:         "OpenAI takes priority",
			openAIKey:    "sk-test",
			anthropicKey: "sk-ant-test",
			wantProvider: ProviderOpenAI,
			wantAvail:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use invalid Ollama URL to skip Ollama check
			d := NewDetector("http://invalid:99999", tt.openAIKey, tt.anthropicKey)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			result := d.Detect(ctx)
			if result.Available != tt.wantAvail {
				t.Errorf("Available = %v, want %v", result.Available, tt.wantAvail)
			}
			if result.ProviderType != tt.wantProvider {
				t.Errorf("ProviderType = %v, want %v", result.ProviderType, tt.wantProvider)
			}
		})
	}
}

func TestDetectNoProviders(t *testing.T) {
	d := NewDetector("http://invalid:99999", "", "")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	result := d.Detect(ctx)
	if result.Available {
		t.Error("expected no provider available")
	}
	if result.Message == "" {
		t.Error("expected error message")
	}
}

func TestDetectionResult(t *testing.T) {
	result := DetectionResult{
		Available:    true,
		ProviderType: ProviderOllama,
		ProviderName: "Ollama",
		Models:       []string{"llama2", "codellama"},
		Message:      "Ollama running",
	}

	if !result.Available {
		t.Error("expected Available to be true")
	}
	if len(result.Models) != 2 {
		t.Errorf("expected 2 models, got %d", len(result.Models))
	}
}
