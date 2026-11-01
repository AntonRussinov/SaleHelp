package user

import (
	"context"
	"fmt"

	"salehelp/internal/storage/pggorm/model"
)

// SetTelegramChatID привязывает telegram_chat_id к пользователю (для объединения аккаунтов).
func (r *Repository) SetTelegramChatID(ctx context.Context, userID int, chatID int64) error {
	if userID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	res := r.DB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("telegram_chat_id", chatID)
	if res.Error != nil {
		return fmt.Errorf("user set telegram_chat_id: %w", res.Error)
	}
	return nil
}

// ClearTelegramChatID отвязывает telegram от пользователя (после слияния в другого пользователя).
func (r *Repository) ClearTelegramChatID(ctx context.Context, userID int) error {
	if userID <= 0 {
		return fmt.Errorf("user_id is required")
	}
	res := r.DB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("telegram_chat_id", nil)
	if res.Error != nil {
		return fmt.Errorf("user clear telegram_chat_id: %w", res.Error)
	}
	return nil
}
