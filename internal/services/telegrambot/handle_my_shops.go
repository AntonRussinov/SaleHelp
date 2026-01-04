package telegrambot

import (
	"context"
	"fmt"
)

func (s *Service) handleMyShops(ctx context.Context, chatID int64) error {
	// 1) user
	u, err := s.ensureUserByChatID(ctx, chatID)
	if err != nil {
		return err
	}

	// 2) shops
	shops, err := s.shopsRepo.ListByUser(ctx, u.ID, 0, 0)
	if err != nil {
		return fmt.Errorf("list shops by user: %w", err)
	}

	// 3) message
	var text string
	if len(shops) == 0 {
		text = "📦 У вас пока нет подключённых магазинов.\n\nИспользуйте /addshop, чтобы добавить магазин."
	} else {
		text = "📦 <b>Ваши магазины:</b>\n\n"
		for i, sh := range shops {
			status := "⛔ выключен"
			if sh.IsActive {
				status = "✅ активен"
			}

			text += fmt.Sprintf(
				"%d) <b>%s</b>\n   id: %d\n   статус: %s\n\n",
				i+1,
				sh.Name,
				sh.ID,
				status,
			)
		}
	}

	// 4) send
	return s.tg.SendMessage(
		ctx,
		chatID,
		text,
	)
}
