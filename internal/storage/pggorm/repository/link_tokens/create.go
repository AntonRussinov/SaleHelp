package linktokens

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) Create(ctx context.Context, t *domain.LinkToken) error {
	m := model.LinkToken{
		Token:          t.Token,
		UserID:         t.UserID,
		TelegramChatID: t.TelegramChatID,
		ExpiresAt:      t.ExpiresAt,
		CreatedAt:      t.CreatedAt,
	}
	if err := r.DB.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("link_tokens create: %w", err)
	}
	return nil
}
