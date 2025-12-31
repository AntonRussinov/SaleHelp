package sellerwarehouse

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByShopAndWBOfficeID(ctx context.Context, shopID int, officeID int64) (*domain.SellerWarehouse, error) {
	var m model.SellerWarehouse

	err := r.DB.WithContext(ctx).
		First(&m, "shop_id = ? AND wb_office_id = ?", shopID, officeID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("seller_warehouse get: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByShop(ctx context.Context, shopID int) ([]*domain.SellerWarehouse, error) {
	var rows []model.SellerWarehouse

	if err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("seller_warehouse list: %w", err)
	}

	result := make([]*domain.SellerWarehouse, 0, len(rows))
	for _, r := range rows {
		result = append(result, r.ToDomain())
	}

	return result, nil
}
