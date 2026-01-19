package shopinvitetokens

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) Create(ctx context.Context, t *domain.ShopInviteToken) error {
	m := model.ShopInviteToken{
		ShopID:          t.ShopID,
		Token:           t.Token,
		Role:            t.Role,
		ExpiresAt:       t.ExpiresAt,
		UsedAt:          t.UsedAt,
		CreatedByChatID: t.CreatedByChatID,
	}

	if err := r.DB.WithContext(ctx).Create(&m).Error; err != nil {
		return fmt.Errorf("shop_invite_tokens create: %w", err)
	}

	t.ID = m.ID
	t.CreatedAt = m.CreatedAt
	return nil
}
