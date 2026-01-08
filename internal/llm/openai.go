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

// OpenAIProvider implements Provider for OpenAI
type OpenAIProvider struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	if model == "" {
		model = "gpt-4"
	}
	return &OpenAIProvider{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.openai.com/v1",
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (o *OpenAIProvider) Name() string {
	return "openai"
}

func (o *OpenAIProvider) Available(ctx context.Context) bool {
	return o.apiKey != ""
}

func (o *OpenAIProvider) Models(ctx context.Context) ([]string, error) {
	return []string{"gpt-4", "gpt-4-turbo", "gpt-3.5-turbo"}, nil
}

func (o *OpenAIProvider) Complete(ctx context.Context, prompt string, opts Options) (string, error) {
	model := opts.Model
	if model == "" {
		model = o.model
	}

	messages := []map[string]string{}
	if opts.SystemPrompt != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": opts.SystemPrompt,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": prompt,
	})

	reqBody := map[string]interface{}{
		"model":       model,
		"messages":    messages,
		"temperature": opts.Temperature,
		"max_tokens":  opts.MaxTokens,
	}
	if len(opts.Stop) > 0 {
		reqBody["stop"] = opts.Stop
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return result.Choices[0].Message.Content, nil
}
