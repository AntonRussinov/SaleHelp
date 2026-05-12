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

// GetActiveByUserID returns the current active/trial subscription for a user, or nil.
func (r *Repository) GetActiveByUserID(ctx context.Context, userID int) (*domain.Subscription, error) {
	var m model.Subscription
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND status IN ('active','trial') AND expires_at > ?", userID, time.Now()).
		Order("expires_at DESC").
		Limit(1).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get active subscription by user_id: %w", err)
	}
	return m.ToDomain(), nil
}

// GetLatestByUserID returns the most recent subscription for a user (any status), or nil.
func (r *Repository) GetLatestByUserID(ctx context.Context, userID int) (*domain.Subscription, error) {
	var m model.Subscription
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("expires_at DESC").
		Limit(1).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get latest subscription by user_id: %w", err)
	}
	return m.ToDomain(), nil
}

// ListExpiringBetween returns active/trial subscriptions with expires_at in [from, to).
func (r *Repository) ListExpiringBetween(ctx context.Context, from, to time.Time) ([]*domain.Subscription, error) {
	var ms []model.Subscription
	err := r.DB.WithContext(ctx).
		Where("status IN ('active','trial') AND expires_at >= ? AND expires_at < ?", from, to).
		Order("expires_at ASC").
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("list expiring subscriptions: %w", err)
	}
	res := make([]*domain.Subscription, 0, len(ms))
	for i := range ms {
		res = append(res, ms[i].ToDomain())
	}
	return res, nil
}
