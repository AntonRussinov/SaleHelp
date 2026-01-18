package shop

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"time"
)

func (r *Repository) Update(ctx context.Context, s *domain.Shop) error {
	m := toModel(s)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("shop update: %w", err)
	}

	s.UpdatedAt = m.UpdatedAt

	return nil
}

func (r *Repository) SetActive(
	ctx context.Context,
	shopID int,
	active bool,
) error {
	err := r.DB.WithContext(ctx).
		Model(&model.Shop{}).
		Where("id = ?", shopID).
		Update("is_active", active).Error

	if err != nil {
		return fmt.Errorf("shop set active=%v: %w", active, err)
	}
	return nil
}

func (r *Repository) UpdateLastSyncAt(ctx context.Context, shopID int, t time.Time) error {
	if shopID <= 0 {
		return fmt.Errorf("shopID is required")
	}

	if err := r.DB.WithContext(ctx).
		Table("shops").
		Where("id = ?", shopID).
		Update("last_sync_at", t).Error; err != nil {
		return fmt.Errorf("update shops.last_sync_at: %w", err)
	}

	return nil
}

func (r *Repository) UpdateFirstSyncNotifiedAt(ctx context.Context, shopID int, t time.Time) error {
	if err := r.DB.WithContext(ctx).
		Table("shops").
		Where("id = ?", shopID).
		Update("first_sync_notified_at", t).Error; err != nil {
		return fmt.Errorf("update shops.first_sync_notified_at: %w", err)
	}
	return nil
}

func (r *Repository) UpdateSyncStatus(ctx context.Context, shopID int, upd model.SyncStatusUpdate) error {
	if shopID <= 0 {
		return fmt.Errorf("shopID is required")
	}

	updates := map[string]any{}

	if upd.LastSyncAt != nil {
		updates["last_sync_at"] = *upd.LastSyncAt
	}
	if upd.LastSyncStatusCode != nil {
		updates["last_sync_status_code"] = *upd.LastSyncStatusCode
	}
	if upd.LastSyncError != nil {
		updates["last_sync_error"] = *upd.LastSyncError
	}
	if upd.LastSyncFailedAt != nil {
		updates["last_sync_failed_at"] = *upd.LastSyncFailedAt
	}
	if upd.LastSyncOkAt != nil {
		updates["last_sync_ok_at"] = *upd.LastSyncOkAt
	}

	// если нечего обновлять — это ок
	if len(updates) == 0 {
		return nil
	}

	if err := r.DB.WithContext(ctx).
		Table("shops").
		Where("id = ?", shopID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("update shops sync status: %w", err)
	}

	return nil
}
