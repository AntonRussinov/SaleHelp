package shop_ai_settings

import (
	"context"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

// GetByShopID возвращает настройки AI для магазина. Если записей нет — nil, nil.
func (r *Repository) GetByShopID(ctx context.Context, shopID int) (*domain.ShopAISettings, error) {
	var m model.ShopAISettings
	err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("shop_ai_settings get: %w", err)
	}
	return m.ToDomain(), nil
}
