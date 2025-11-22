package telegramaccount

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

// GetByUserID — найти по user_id (UNIQUE user_id)
func (r *Repository) GetByUserID(ctx context.Context, userID int) (*domain.TelegramAccount, error) {
	var m model.TelegramAccount

	if err := r.DB.WithContext(ctx).
		First(&m, "user_id = ?", userID).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("telegram_account get by user_id: %w", err)
	}

	return m.ToDomain(), nil
}

// GetByTelegramID — найти по telegram_id (UNIQUE telegram_id)
func (r *Repository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.TelegramAccount, error) {
	var m model.TelegramAccount

	if err := r.DB.WithContext(ctx).
		First(&m, "telegram_id = ?", telegramID).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("telegram_account get by telegram_id: %w", err)
	}

	return m.ToDomain(), nil
}
