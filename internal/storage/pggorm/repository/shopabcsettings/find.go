package shopabcsettings

import (
	"context"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

// GetByShopID возвращает настройки ABC для магазина. Если записи нет — nil (использовать дефолты: всё включено).
func (r *Repository) GetByShopID(ctx context.Context, shopID int) (*domain.ShopABCSettings, error) {
	var m model.ShopABCSettings
	err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("shop_abc_settings get: %w", err)
	}
	return m.ToDomain(), nil
}
