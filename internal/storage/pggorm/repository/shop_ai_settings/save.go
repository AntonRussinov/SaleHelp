package shop_ai_settings

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

// Save создаёт или обновляет настройки AI для магазина.
func (r *Repository) Save(ctx context.Context, d *domain.ShopAISettings) error {
	now := time.Now()
	m := &model.ShopAISettings{
		ShopID:      d.ShopID,
		ToneOfVoice: d.ToneOfVoice,
		BrandName:   d.BrandName,
		CustomRules: d.CustomRules,
		IsActive:    d.IsActive,
		UpdatedAt:   now,
	}
	err := r.DB.WithContext(ctx).
		Save(m).Error
	if err != nil {
		return fmt.Errorf("shop_ai_settings save: %w", err)
	}
	d.UpdatedAt = now
	return nil
}
