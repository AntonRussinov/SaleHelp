package notify

import "context"

type Channel string

const (
	ChannelTelegram Channel = "telegram"
	ChannelEmail    Channel = "email"
	ChannelWhatsApp Channel = "whatsapp"
)

type Message struct {
	Title         string
	Text          string
	InlineKeyboard map[string]any // optional; при наличии Telegram отправляет последний фрагмент с reply_markup
}

type Recipient struct {
	ShopID         int
	TelegramChatID *int64
	Email          *string
	WhatsAppPhone  *string
}

type Sender interface {
	Channel() Channel
	Send(ctx context.Context, r Recipient, m Message) error
}
