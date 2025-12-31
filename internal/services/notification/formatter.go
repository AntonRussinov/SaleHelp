package notifications

import (
	"salehelp/domain"
	"salehelp/pkg/notify"
)

func FormatMessage(at *domain.AlertTriggered) notify.Message {
	// MVP: просто прокидываем готовый текст
	return notify.Message{
		Title: "SaleHelp alert",
		Text:  at.Message,
	}
}
