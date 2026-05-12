package subscription

import (
	"context"
	"fmt"
	"time"
)

// CountByState returns counts of active (incl. trial), grace-period, and expired subscriptions.
func (r *Repository) CountByState(ctx context.Context, now time.Time) (active, grace, expired int, err error) {
	var aCount, gCount, eCount int64

	graceCutoff := now.Add(-72 * time.Hour)

	if err = r.DB.WithContext(ctx).Table("subscriptions").
		Where("status IN ('active','trial') AND expires_at > ?", now).
		Count(&aCount).Error; err != nil {
		return 0, 0, 0, fmt.Errorf("count active subs: %w", err)
	}
	if err = r.DB.WithContext(ctx).Table("subscriptions").
		Where("expires_at <= ? AND expires_at > ?", now, graceCutoff).
		Count(&gCount).Error; err != nil {
		return 0, 0, 0, fmt.Errorf("count grace subs: %w", err)
	}
	if err = r.DB.WithContext(ctx).Table("subscriptions").
		Where("status = 'expired' OR expires_at <= ?", graceCutoff).
		Count(&eCount).Error; err != nil {
		return 0, 0, 0, fmt.Errorf("count expired subs: %w", err)
	}
	return int(aCount), int(gCount), int(eCount), nil
}
