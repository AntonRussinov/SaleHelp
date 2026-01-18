package subscription

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) CountByState(ctx context.Context, now time.Time) (int, int, int, error) {
	var active, grace, expired int64

	if err := r.DB.WithContext(ctx).
		Table("subscriptions").
		Where("status = ? AND valid_to > ?", "active", now).
		Count(&active).Error; err != nil {
		return 0, 0, 0, fmt.Errorf("count active subs: %w", err)
	}

	if err := r.DB.WithContext(ctx).
		Table("subscriptions").
		Where("status = ? AND valid_to <= ? AND (valid_to + interval '3 days') > ?", "active", now, now).
		Count(&grace).Error; err != nil {
		return 0, 0, 0, fmt.Errorf("count grace subs: %w", err)
	}

	// expired
	if err := r.DB.WithContext(ctx).
		Table("subscriptions").
		Where("status = ? OR (valid_to + interval '3 days') <= ?", "expired", now).
		Count(&expired).Error; err != nil {
		return 0, 0, 0, fmt.Errorf("count expired subs: %w", err)
	}

	return int(active), int(grace), int(expired), nil
}
