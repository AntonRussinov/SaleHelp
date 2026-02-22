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

// BotCommand — команда для меню бота (Bot Command Menu, иконка ☰).
type BotCommand struct {
	Command     string `json:"command"`     // 1-32 символа, латиница/цифры/подчёркивание
	Description string `json:"description"` // 3-256 символов, описание на любом языке
}

// SetMyCommands устанавливает список команд бота (меню в чате).
// Команды отображаются в меню рядом с полем ввода при нажатии на ☰.
func (c *Client) SetMyCommands(ctx context.Context, commands []BotCommand) error {
	if len(commands) == 0 {
		return nil
	}
	payload := map[string]any{"commands": commands}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram marshal setMyCommands: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/setMyCommands", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("telegram new setMyCommands request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("telegram setMyCommands request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var terr telegramErrorResponse
	_ = json.Unmarshal(body, &terr)
	if !terr.Ok && terr.Description != "" {
		return fmt.Errorf("telegram setMyCommands api error (%d): %s", terr.ErrorCode, terr.Description)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if terr.Description != "" {
			return fmt.Errorf("telegram setMyCommands http %s: %s", resp.Status, terr.Description)
		}
		return fmt.Errorf("telegram setMyCommands http status: %s", resp.Status)
	}
	return nil
}

// AnswerCallbackQuery снимает "часики" с inline-кнопки после нажатия.
func (c *Client) AnswerCallbackQuery(ctx context.Context, callbackQueryID string) error {
	payload := map[string]any{"callback_query_id": callbackQueryID}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram marshal answerCallbackQuery: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/answerCallbackQuery", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("telegram new answerCallbackQuery request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("telegram answerCallbackQuery request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var terr telegramErrorResponse
	_ = json.Unmarshal(body, &terr)
	if !terr.Ok && terr.Description != "" {
		return fmt.Errorf("telegram answerCallbackQuery api error (%d): %s", terr.ErrorCode, terr.Description)
	}
	return nil
}

func (c *Client) SendRaw(ctx context.Context, payload map[string]any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("telegram marshal raw request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/sendMessage",
		bytes.NewReader(b),
	)
	if err != nil {
		return fmt.Errorf("telegram new raw request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("telegram http raw request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var terr telegramErrorResponse
	_ = json.Unmarshal(body, &terr)

	if terr.Ok == false && terr.Description != "" {
		return fmt.Errorf("telegram api error (%d): %s", terr.ErrorCode, terr.Description)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram http status: %s", resp.Status)
	}

	return nil
}
