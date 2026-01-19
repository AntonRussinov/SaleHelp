package notifications

import (
	"fmt"
	"strconv"

	"salehelp/domain"
	"salehelp/pkg/notify"
)

func recipientFromEndpoint(ep *domain.NotificationEndpoint, shopID int) (notify.Recipient, error) {
	rcpt := notify.Recipient{ShopID: shopID}

	ch := notify.Channel(ep.Channel)
	switch ch {
	case notify.ChannelTelegram:
		chatID, err := strconv.ParseInt(ep.Address, 10, 64)
		if err != nil {
			return notify.Recipient{}, fmt.Errorf("invalid telegram chat_id: %q", ep.Address)
		}
		rcpt.TelegramChatID = &chatID
		return rcpt, nil

	case notify.ChannelEmail:
		email := ep.Address
		if email == "" {
			return notify.Recipient{}, fmt.Errorf("empty email address")
		}
		rcpt.Email = &email
		return rcpt, nil

	case notify.ChannelWhatsApp:
		phone := ep.Address
		if phone == "" {
			return notify.Recipient{}, fmt.Errorf("empty whatsapp phone")
		}
		rcpt.WhatsAppPhone = &phone
		return rcpt, nil

	default:
		return notify.Recipient{}, fmt.Errorf("unsupported channel: %s", ep.Channel)
	}
}
