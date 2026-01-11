package telegrambot

import (
	"context"
	"fmt"
)

func (s *Service) canAddTelegramEndpoint(ctx context.Context, shopID int) (bool, string) {
	allowed, err := s.subEntitlements.AllowedTelegramChats(ctx, shopID)
	if err != nil {
		// безопасное сообщение пользователю
		return false, "Не удалось проверить подписку. Попробуй позже."
	}

	// Политика v0: без активной подписки нельзя добавлять чаты
	if allowed <= 0 {
		return false, "У магазина нет активной подписки. Алерты и уведомления доступны только при активной подписке."
	}

	current, err := s.endpointsRepo.CountEnabledByShopAndChannel(ctx, shopID, "telegram")
	if err != nil {
		return false, "Не удалось проверить текущие чаты. Попробуй позже."
	}

	if current >= allowed {
		return false, fmt.Sprintf(
			"Достигнут лимит Telegram-чатов по подписке: %d из %d.\n\nПродли подписку или подключи доп. чаты, и попробуй снова.",
			current, allowed,
		)
	}

	return true, ""
}
