package shop

import (
	"context"
	"fmt"
)

func (r *Repository) CountAll(ctx context.Context) (int, error) {
	var cnt int64
	if err := r.DB.WithContext(ctx).Table("shops").Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("count all shops: %w", err)
	}
	return int(cnt), nil
}

func (r *Repository) CountActive(ctx context.Context) (int, error) {
	var cnt int64
	if err := r.DB.WithContext(ctx).
		Table("shops").
		Where("is_active = TRUE").
		Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("count active shops: %w", err)
	}
	return int(cnt), nil
}

func (r *Repository) CountByLastSyncStatusCode(ctx context.Context, code int) (int, error) {
	var cnt int64
	if err := r.DB.WithContext(ctx).
		Table("shops").
		Where("last_sync_status_code = ?", code).
		Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("count shops by last_sync_status_code: %w", err)
	}
	return int(cnt), nil
}
