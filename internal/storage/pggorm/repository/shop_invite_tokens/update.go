package shopinvitetokens

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
	"time"
)

func (r *Repository) MarkUsed(ctx context.Context, token string, usedAt time.Time) error {
	res := r.DB.WithContext(ctx).
		Model(&model.ShopInviteToken{}).
		Where("token = ?", token).
		Where("used_at IS NULL").
		Update("used_at", usedAt)

	if res.Error != nil {
		return fmt.Errorf("shop_invite_tokens mark used: %w", res.Error)
	}
	return nil
}

// InvalidateActiveByShop — удобно, чтобы держать максимум 1 активный инвайт на магазин.
func (r *Repository) InvalidateActiveByShop(ctx context.Context, shopID int, usedAt time.Time) error {
	res := r.DB.WithContext(ctx).
		Model(&model.ShopInviteToken{}).
		Where("shop_id = ?", shopID).
		Where("used_at IS NULL").
		Update("used_at", usedAt)

	if res.Error != nil {
		return fmt.Errorf("shop_invite_tokens invalidate by shop: %w", res.Error)
	}
	return nil
}
