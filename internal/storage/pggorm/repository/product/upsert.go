package product

import (
	"context"
	"fmt"

	"salehelp/domain"
)

func (r *Repository) Upsert(ctx context.Context, p *domain.Product) error {
	// 1. Проверяем, есть ли продукт
	existing, err := r.GetByShopAndExternalID(ctx, p.ShopID, p.ExternalID)
	if err != nil {
		return fmt.Errorf("product upsert get: %w", err)
	}

	// 2. Если нет — создаём новый
	if existing == nil {
		m := toModel(p) // конвертация в модель БД
		if err := r.DB.WithContext(ctx).Create(&m).Error; err != nil {
			return fmt.Errorf("product create: %w", err)
		}
		// обновляем ID в домене
		p.ID = m.ID
		return nil
	}

	// 3. Если есть — обновляем поля
	existing.Title = p.Title
	existing.Brand = p.Brand
	existing.Category = p.Category
	existing.Description = p.Description
	existing.URL = p.URL
	existing.ImageURL = p.ImageURL

	m := toModel(existing)

	if err := r.DB.WithContext(ctx).Save(&m).Error; err != nil {
		return fmt.Errorf("product update: %w", err)
	}

	p.ID = m.ID
	return nil
}
