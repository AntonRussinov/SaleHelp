package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		baseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type sendMessageRequest struct {
	ChatID                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode,omitempty"`               // "HTML" или "MarkdownV2"
	DisableWebPagePreview bool   `json:"disable_web_page_preview,omitempty"` // true = без превью ссылок
}

type telegramErrorResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	ErrorCode   int    `json:"error_code"`
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string, parseMode string, disablePreview bool) error {
	if chatID == 0 {
		return fmt.Errorf("telegram chat_id is 0")
	}
	if text == "" {
		return fmt.Errorf("telegram text is empty")
	}

	reqBody, err := json.Marshal(sendMessageRequest{
		ChatID:                chatID,
		Text:                  text,
		ParseMode:             parseMode,
		DisableWebPagePreview: disablePreview,
	})
	if err != nil {
		return fmt.Errorf("telegram marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/sendMessage", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("telegram new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("telegram http request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Telegram ошибки могут прийти и с 200, и с 4xx/5xx
	var terr telegramErrorResponse
	_ = json.Unmarshal(body, &terr)
	if terr.Ok == false && terr.Description != "" {
		return fmt.Errorf("telegram api error (%d): %s", terr.ErrorCode, terr.Description)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if terr.Description != "" {
			return fmt.Errorf("telegram http %s: %s", resp.Status, terr.Description)
		}
		return fmt.Errorf("telegram http status: %s", resp.Status)
	}

	return nil
}
