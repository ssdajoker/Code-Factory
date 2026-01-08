package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AnthropicProvider implements Provider for Claude
type AnthropicProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	if model == "" {
		model = "claude-3-sonnet-20240229"
	}
	return &AnthropicProvider{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.anthropic.com/v1",
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (a *AnthropicProvider) Name() string {
	return "anthropic"
}

func (a *AnthropicProvider) Available(ctx context.Context) bool {
	return a.apiKey != ""
}

func (a *AnthropicProvider) Models(ctx context.Context) ([]string, error) {
	return []string{"claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307"}, nil
}

func (a *AnthropicProvider) Complete(ctx context.Context, prompt string, opts Options) (string, error) {
	model := opts.Model
	if model == "" {
		model = a.model
	}

	reqBody := map[string]interface{}{
		"model":      model,
		"max_tokens": opts.MaxTokens,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	if opts.SystemPrompt != "" {
		reqBody["system"] = opts.SystemPrompt
	}
	if opts.Temperature > 0 {
		reqBody["temperature"] = opts.Temperature
	}
	if len(opts.Stop) > 0 {
		reqBody["stop_sequences"] = opts.Stop
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/messages", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anthropic error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no response from Anthropic")
	}

	return result.Content[0].Text, nil
}
