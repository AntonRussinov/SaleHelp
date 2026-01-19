package telegram

import (
	"salehelp/internal/services/telegrambot"
)

type Handler struct {
	bot    *telegrambot.Service
	secret string
}

func NewHandler(bot *telegrambot.Service, secret string) *Handler {
	return &Handler{
		bot:    bot,
		secret: secret,
	}
}
