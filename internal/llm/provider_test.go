package llm

import (
	"context"
	"testing"
)

// MockProvider implements Provider for testing
type MockProvider struct {
	name       string
	available  bool
	models     []string
	response   string
	err        error
}

func (m *MockProvider) Complete(ctx context.Context, prompt string, opts Options) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Available(ctx context.Context) bool {
	return m.available
}

func (m *MockProvider) Models(ctx context.Context) ([]string, error) {
	return m.models, m.err
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want 0.7", opts.Temperature)
	}
	if opts.MaxTokens != 2048 {
		t.Errorf("MaxTokens = %v, want 2048", opts.MaxTokens)
	}
	if opts.SystemPrompt == "" {
		t.Error("SystemPrompt should not be empty")
	}
}

func TestNewProvider(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "ollama provider",
			cfg: Config{
				Type:  ProviderOllama,
				Model: "llama2",
			},
			wantErr: false,
		},
		{
			name: "openai without key",
			cfg: Config{
				Type:  ProviderOpenAI,
				Model: "gpt-4",
			},
			wantErr: true,
		},
		{
			name: "openai with key",
			cfg: Config{
				Type:   ProviderOpenAI,
				APIKey: "sk-test",
				Model:  "gpt-4",
			},
			wantErr: false,
		},
		{
			name: "anthropic without key",
			cfg: Config{
				Type:  ProviderAnthropic,
				Model: "claude-3",
			},
			wantErr: true,
		},
		{
			name: "anthropic with key",
			cfg: Config{
				Type:   ProviderAnthropic,
				APIKey: "sk-ant-test",
				Model:  "claude-3",
			},
			wantErr: false,
		},
		{
			name: "unknown provider",
			cfg: Config{
				Type: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProvider(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProviderTypes(t *testing.T) {
	if ProviderOllama != "ollama" {
		t.Errorf("ProviderOllama = %v, want 'ollama'", ProviderOllama)
	}
	if ProviderOpenAI != "openai" {
		t.Errorf("ProviderOpenAI = %v, want 'openai'", ProviderOpenAI)
	}
	if ProviderAnthropic != "anthropic" {
		t.Errorf("ProviderAnthropic = %v, want 'anthropic'", ProviderAnthropic)
	}
}

func TestMockProvider(t *testing.T) {
	mock := &MockProvider{
		name:      "test",
		available: true,
		models:    []string{"model1", "model2"},
		response:  "test response",
	}

	ctx := context.Background()

	if mock.Name() != "test" {
		t.Errorf("Name() = %v, want 'test'", mock.Name())
	}

	if !mock.Available(ctx) {
		t.Error("Available() should return true")
	}

	models, _ := mock.Models(ctx)
	if len(models) != 2 {
		t.Errorf("Models() returned %d models, want 2", len(models))
	}

	resp, _ := mock.Complete(ctx, "test", DefaultOptions())
	if resp != "test response" {
		t.Errorf("Complete() = %v, want 'test response'", resp)
	}
}
