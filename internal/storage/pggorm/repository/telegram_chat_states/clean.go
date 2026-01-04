package telegramchatstates

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
	"time"
)

func (r *Repository) Clear(ctx context.Context, chatID int64) error {
	if err := r.DB.WithContext(ctx).
		Delete(&model.TelegramChatState{}, "chat_id = ?", chatID).Error; err != nil {
		return fmt.Errorf("telegram_chat_states clear: %w", err)
	}
	return nil
}

func (r *Repository) CleanupOlderThan(ctx context.Context, olderThan time.Time) (int64, error) {
	res := r.DB.WithContext(ctx).
		Where("updated_at < ?", olderThan).
		Delete(&model.TelegramChatState{})

	if res.Error != nil {
		return 0, fmt.Errorf("telegram_chat_states cleanup: %w", res.Error)
	}
	return res.RowsAffected, nil
}
