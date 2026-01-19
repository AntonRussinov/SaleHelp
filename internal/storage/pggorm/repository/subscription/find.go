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

func (r *Repository) ListByValidToRange(ctx context.Context, from, to time.Time) ([]*domain.Subscription, error) {
	var ms []model.Subscription

	err := r.DB.WithContext(ctx).
		Where("status = ? AND valid_to >= ? AND valid_to < ?", "active", from, to).
		Order("valid_to ASC").
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("list subscriptions by valid_to range: %w", err)
	}

	res := make([]*domain.Subscription, 0, len(ms))
	for i := range ms {
		res = append(res, ms[i].ToDomain())
	}
	return res, nil
}

// GetLatestByShopID — последняя подписка по valid_to DESC, либо nil.
func (r *Repository) GetLatestByShopID(ctx context.Context, shopID int) (*domain.Subscription, error) {
	var m model.Subscription
	err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Order("valid_to DESC").
		Limit(1).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get latest subscription by shop_id: %w", err)
	}
	return m.ToDomain(), nil
}
