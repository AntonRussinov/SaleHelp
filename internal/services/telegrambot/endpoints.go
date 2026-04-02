package telegrambot

import (
	"context"
	"fmt"
	"salehelp/domain"
)

// ensureOwnerEndpointsForUser создаёт endpoint (label=owner) для всех магазинов пользователя,
// если такого endpoint ещё нет. Вызывается после привязки Telegram-чата к аккаунту.
func (s *Service) ensureOwnerEndpointsForUser(ctx context.Context, userID int, chatID int64) {
	shops, err := s.shopsRepo.ListByUser(ctx, userID, 100, 0)
	if err != nil {
		s.log.Warn("ensureOwnerEndpoints: list shops failed", "user_id", userID, "err", err)
		return
	}

	address := fmt.Sprintf("%d", chatID)
	label := "owner"

	for _, shop := range shops {
		exists, err := s.endpointsRepo.Exists(ctx, shop.ID, "telegram", address)
		if err != nil {
			s.log.Warn("ensureOwnerEndpoints: exists check failed", "shop_id", shop.ID, "err", err)
			continue
		}
		if exists {
			continue
		}

		ep := &domain.NotificationEndpoint{
			ShopID:  shop.ID,
			Channel: "telegram",
			Address: address,
			Label:   &label,
			Enabled: true,
		}
		if err := s.endpointsRepo.Create(ctx, ep); err != nil {
			s.log.Warn("ensureOwnerEndpoints: create failed", "shop_id", shop.ID, "err", err)
		}
	}
}
