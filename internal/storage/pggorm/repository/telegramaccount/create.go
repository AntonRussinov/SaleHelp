package telegramaccount

import (
	"context"
	"fmt"
	"salehelp/domain"
)

// Create — регистрация связи user ↔ telegram
func (r *Repository) Create(ctx context.Context, a *domain.TelegramAccount) error {
	m := toModel(a)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("telegram_account create: %w", err)
	}

	a.ID = m.ID
	a.CreatedAt = m.CreatedAt

	return nil
}
