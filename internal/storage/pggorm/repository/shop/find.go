package shop

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.Shop, error) {
	var m model.Shop

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("shop get by id: %w", err)
	}

	return m.ToDomain(), nil
}

// GetByUserAndMarketplace — один магазин на пользователя и маркетплейс
func (r *Repository) GetByUserAndMarketplace(
	ctx context.Context,
	userID int,
	marketplace string,
) (*domain.Shop, error) {
	var m model.Shop

	if err := r.DB.WithContext(ctx).
		First(&m, "user_id = ? AND marketplace = ?", userID, marketplace).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("shop get by user_id+marketplace: %w", err)
	}

	return m.ToDomain(), nil
}

// ListByUser — все магазины пользователя (для UI/аналитики)
func (r *Repository) ListByUser(
	ctx context.Context,
	userID int,
	limit, offset int,
) ([]*domain.Shop, error) {
	var models []model.Shop

	q := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id ASC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("shop list by user: %w", err)
	}

	result := make([]*domain.Shop, 0, len(models))
	for i := range models {
		result = append(result, models[i].ToDomain())
	}

	return result, nil
}
