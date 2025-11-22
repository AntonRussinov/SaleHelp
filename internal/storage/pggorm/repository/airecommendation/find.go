package airecommendation

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.AIRecommendation, error) {
	var m model.AIRecommendation

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("ai_recommendation get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByShop(ctx context.Context, shopID int) ([]*domain.AIRecommendation, error) {
	var rows []model.AIRecommendation

	if err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Order("id DESC").
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("ai_recommendation list by shop: %w", err)
	}

	out := make([]*domain.AIRecommendation, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}

func (r *Repository) ListByProduct(ctx context.Context, productID int) ([]*domain.AIRecommendation, error) {
	var rows []model.AIRecommendation

	if err := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("id DESC").
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("ai_recommendation list by product: %w", err)
	}

	out := make([]*domain.AIRecommendation, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}

func (r *Repository) ListByAlert(ctx context.Context, alertID int) ([]*domain.AIRecommendation, error) {
	var rows []model.AIRecommendation

	if err := r.DB.WithContext(ctx).
		Where("alert_id = ?", alertID).
		Order("id DESC").
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("ai_recommendation list by alert: %w", err)
	}

	out := make([]*domain.AIRecommendation, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}
