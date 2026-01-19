package telegramchatstates

import (
	"context"
	"errors"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByChatID(ctx context.Context, chatID int64) (*domain.TelegramChatState, error) {
	var m model.TelegramChatState

	err := r.DB.WithContext(ctx).
		Where("chat_id = ?", chatID).
		First(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("telegram_chat_states get: %w", err)
	}

	// конвертим payload в RawMessage
	return &domain.TelegramChatState{
		ChatID:    m.ChatID,
		State:     m.State,
		Payload:   domain.JSONRaw([]byte(m.Payload)),
		UpdatedAt: m.UpdatedAt,
	}, nil
}
