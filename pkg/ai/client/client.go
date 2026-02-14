package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"salehelp/pkg/ai"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

const maxOutputTokens = 600 // ~500 символов + запас

// Reasoning-модели (gpt-5, gpt-5-mini, gpt-5.2 и т.д.) через Chat Completions возвращают пустой Content.
// Для них используем Responses API, где ответ в output[].content[].text.
var reasoningModelPrefixes = []string{"gpt-5", "o1-", "o3-", "o4-"}

func isReasoningModel(model string) bool {
	for _, p := range reasoningModelPrefixes {
		if strings.HasPrefix(model, p) {
			return true
		}
	}
	return false
}

// Client — обёртка над OpenAI LLM для генерации ответов на отзывы.
type Client struct {
	llm       llms.Model
	models    ai.ModelConfig
	apiKey    string
	httpClient *http.Client
}

// New создаёт клиент. apiKey — ключ из конфига; если пустой, используется env OPENAI_API_KEY.
func New(models ai.ModelConfig, apiKey string, opts ...openai.Option) (*Client, error) {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if apiKey != "" {
		opts = append(opts, openai.WithToken(apiKey))
	}
	llm, err := openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("openai new: %w", err)
	}
	if models.ModelForPositive == "" {
		models = ai.DefaultModelConfig()
	}
	return &Client{
		llm:        llm,
		models:     models,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

// GenerateFeedbackReply генерирует ответ на отзыв. Модель выбирается по рейтингу (5 → mini, 1–4 → other).
// Для reasoning-моделей (gpt-5-mini, gpt-5.2 и т.д.) вызывается Responses API; для остальных — Chat Completions через langchaingo.
func (c *Client) GenerateFeedbackReply(ctx context.Context, prompt string, rating int) (*ai.FeedbackReplyResult, error) {
	model := c.models.ModelForOther
	if rating == 5 {
		model = c.models.ModelForPositive
	}
	if isReasoningModel(model) {
		return c.generateViaResponsesAPI(ctx, model, prompt)
	}
	return c.generateViaChatCompletions(ctx, model, prompt)
}

func (c *Client) generateViaResponsesAPI(ctx context.Context, model, prompt string) (*ai.FeedbackReplyResult, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("openai api key required for responses api")
	}
	reqBody := map[string]any{
		"model":             model,
		"input":             []map[string]string{{"role": "user", "content": prompt}},
		"max_output_tokens": maxOutputTokens,
		"reasoning":         map[string]string{"effort": "low"},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("responses api marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("responses api request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("responses api do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errBody struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return nil, fmt.Errorf("responses api status %d: %s", resp.StatusCode, errBody.Error.Message)
	}
	var out responsesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("responses api decode: %w", err)
	}
	text := extractOutputTextFromResponses(out)
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("model returned empty content (model=%s)", model)
	}
	tokensUsed := 0
	if out.Usage != nil {
		tokensUsed = out.Usage.TotalTokens
	}
	return &ai.FeedbackReplyResult{
		Text:       text,
		TokensUsed: tokensUsed,
		ModelUsed:  model,
	}, nil
}

type responsesAPIResponse struct {
	Status string                   `json:"status"`
	Output []responsesAPIOutputItem  `json:"output"`
	Usage  *responsesAPIUsage        `json:"usage"`
}

type responsesAPIOutputItem struct {
	Type    string                     `json:"type"`
	Content []responsesAPIContentPart  `json:"content"`
}

type responsesAPIContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type responsesAPIUsage struct {
	TotalTokens int `json:"total_tokens"`
}

func extractOutputTextFromResponses(r responsesAPIResponse) string {
	for i := len(r.Output) - 1; i >= 0; i-- {
		item := r.Output[i]
		if item.Type != "message" {
			continue
		}
		for _, part := range item.Content {
			if part.Type == "output_text" && part.Text != "" {
				return part.Text
			}
		}
	}
	return ""
}

func (c *Client) generateViaChatCompletions(ctx context.Context, model, prompt string) (*ai.FeedbackReplyResult, error) {
	opts := []llms.CallOption{
		llms.WithModel(model),
		llms.WithMaxTokens(maxOutputTokens),
	}
	resp, err := c.llm.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}, opts...)
	if err != nil {
		return nil, fmt.Errorf("generate content: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty response from model")
	}
	text := resp.Choices[0].Content
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("model returned empty content (model=%s)", model)
	}
	tokensUsed := extractTokensUsed(resp)
	return &ai.FeedbackReplyResult{
		Text:       text,
		TokensUsed: tokensUsed,
		ModelUsed:  model,
	}, nil
}

func extractTokensUsed(resp *llms.ContentResponse) int {
	if len(resp.Choices) == 0 {
		return 0
	}
	info := resp.Choices[0].GenerationInfo
	if info == nil {
		return 0
	}
	if v, ok := info["total_tokens"]; ok {
		switch n := v.(type) {
		case int:
			return n
		case float64:
			return int(n)
		}
	}
	return 0
}
