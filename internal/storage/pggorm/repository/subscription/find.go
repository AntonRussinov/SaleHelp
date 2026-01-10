package subscription

import (
	"context"
	"errors"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"time"

	"gorm.io/gorm"
)

// GetActiveByShopID возвращает текущую активную подписку (valid_to > now) либо nil.
func (r *Repository) GetActiveByShopID(ctx context.Context, shopID int, now time.Time) (*domain.Subscription, error) {
	var m model.Subscription

	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND status = ? AND valid_to > ?", shopID, "active", now).
		Order("valid_to DESC").
		Limit(1).
		First(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get active subscription by shop_id: %w", err)
	}

	return m.ToDomain(), nil
}
