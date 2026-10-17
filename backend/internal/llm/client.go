package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ai-bot-expenses/backend/internal/config"
)

// Client handles interaction with OpenAI-compatible LLM endpoints.
type Client struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClient creates a new LLM client from application configuration.
func NewClient(cfg *config.Config) *Client {
	baseURL := strings.TrimRight(cfg.LLMBaseURL, "/")
	if !strings.HasSuffix(baseURL, "/chat/completions") {
		baseURL = baseURL + "/chat/completions"
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  cfg.LLMAPIKey,
		model:   cfg.LLMModel,
		httpClient: &http.Client{
			Timeout: 35 * time.Second,
		},
	}
}

// ParseIntent sends the user's message and system prompt to the LLM and returns the parsed structured response.
func (c *Client) ParseIntent(ctx context.Context, userMessage string) (*Response, error) {
	systemPrompt := BuildSystemPrompt(time.Now())

	reqBody := ChatRequest{
		Model: c.model,
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		Temperature: 0.1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal LLM request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request to LLM: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("LLM API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read LLM response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("LLM API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode LLM response: %w", err)
	}

	if chatResp.Error != nil && chatResp.Error.Message != "" {
		return nil, fmt.Errorf("LLM API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("LLM returned empty choices")
	}

	rawContent := strings.TrimSpace(chatResp.Choices[0].Message.Content)

	// Clean code fence formatting if the LLM wrapped JSON in ```json ... ```
	cleanedContent := cleanJSONContent(rawContent)

	var result Response
	if err := json.Unmarshal([]byte(cleanedContent), &result); err != nil {
		return nil, fmt.Errorf("failed to parse LLM structured JSON output: %w (raw response: %s)", err, rawContent)
	}

	return &result, nil
}

// cleanJSONContent strips markdown code fences if present.
func cleanJSONContent(content string) string {
	s := strings.TrimSpace(content)
	if strings.HasPrefix(s, "```") {
		// Remove first line (e.g. ```json or ```)
		if idx := strings.Index(s, "\n"); idx != -1 {
			s = s[idx+1:]
		}
		// Remove trailing ```
		s = strings.TrimSuffix(strings.TrimSpace(s), "```")
	}
	return strings.TrimSpace(s)
}
