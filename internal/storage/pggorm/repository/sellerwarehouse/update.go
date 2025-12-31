package sellerwarehouse

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Update(ctx context.Context, d *domain.SellerWarehouse) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("seller_warehouse update: %w", err)
	}

	d.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *Repository) Upsert(ctx context.Context, d *domain.SellerWarehouse) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "shop_id"},
				{Name: "wb_office_id"},
			}, // UNIQUE (shop_id, wb_office_id)
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"cargo_type",
				"delivery_type",
				"is_deleting",
				"is_processing",
				"updated_at",
			}),
		}).
		Create(m).Error; err != nil {

		return fmt.Errorf("seller_warehouse upsert: %w", err)
	}

	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt

	return nil
}

func (r *Repository) ReplaceForShop(ctx context.Context, shopID int64, warehouses []domain.SellerWarehouse) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Удаляем все старые склады этого магазина
		if err := tx.
			Where("shop_id = ?", shopID).
			Delete(&model.SellerWarehouse{}).Error; err != nil {
			return err
		}

		// 2. Если новые не пришли – всё, на этом заканчиваем
		if len(warehouses) == 0 {
			return nil
		}

		// 3. Вставляем новые
		models := make([]model.SellerWarehouse, 0, len(warehouses))
		for _, d := range warehouses {
			models = append(models, *toModel(&d))
		}

		return tx.Create(&models).Error
	})
}
