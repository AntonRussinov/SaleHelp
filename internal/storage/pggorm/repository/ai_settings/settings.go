package ai_settings

import (
	"context"
	"errors"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByShopID(ctx context.Context, shopID int) (*domain.AISettings, error) {
	var m model.AISettings
	err := r.DB.WithContext(ctx).Where("shop_id = ?", shopID).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ai_settings get by shop: %w", err)
	}
	return m.ToDomain(), nil
}

func (r *Repository) Upsert(ctx context.Context, d *domain.AISettings) error {
	d.UpdatedAt = time.Now()
	m := model.AISettingsFromDomain(d)
	if m.CreatedAt.IsZero() {
		m.CreatedAt = d.UpdatedAt
	}
	result := r.DB.WithContext(ctx).
		Where(model.AISettings{ShopID: d.ShopID}).
		Assign(m).
		FirstOrCreate(m)
	if result.Error != nil {
		return fmt.Errorf("ai_settings upsert: %w", result.Error)
	}
	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *Repository) ListActiveShopIDs(ctx context.Context, limit int) ([]int, error) {
	if limit <= 0 {
		limit = 500
	}
	var ids []int
	err := r.DB.WithContext(ctx).
		Model(&model.AISettings{}).
		Where("is_active = true").
		Limit(limit).
		Pluck("shop_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("ai_settings list active: %w", err)
	}
	return ids, nil
}
