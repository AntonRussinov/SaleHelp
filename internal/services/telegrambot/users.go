package telegrambot

import (
	"context"
	"fmt"

	"salehelp/domain"
)

func (s *Service) ensureUserByChatID(ctx context.Context, chatID int64) (*domain.User, error) {
	u, err := s.usersRepo.EnsureByTelegramChatID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("ensure user by telegram_chat_id: %w", err)
	}
	return u, nil
}
