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

// OllamaProvider implements Provider for Ollama
type OllamaProvider struct {
	baseURL string
	model   string
	client  *http.Client
}

// NewOllamaProvider creates a new Ollama provider
func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if model == "" {
		model = "llama2"
	}
	return &OllamaProvider{
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (o *OllamaProvider) Name() string {
	return "ollama"
}

func (o *OllamaProvider) Available(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", o.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (o *OllamaProvider) Models(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", o.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	models := make([]string, len(result.Models))
	for i, m := range result.Models {
		models[i] = m.Name
	}
	return models, nil
}

func (o *OllamaProvider) Complete(ctx context.Context, prompt string, opts Options) (string, error) {
	model := opts.Model
	if model == "" {
		model = o.model
	}

	reqBody := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": opts.Temperature,
			"num_predict": opts.MaxTokens,
		},
	}
	if opts.SystemPrompt != "" {
		reqBody["system"] = opts.SystemPrompt
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/api/generate", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}
