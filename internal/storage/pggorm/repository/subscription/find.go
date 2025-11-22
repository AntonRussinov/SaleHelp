package subscription

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.Subscription, error) {
	var m model.Subscription

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("subscription get by id: %w", err)
	}

	return m.ToDomain(), nil
}

// ListByUser возвращает все подписки пользователя
func (r *Repository) ListByUser(ctx context.Context, userID int) ([]*domain.Subscription, error) {
	var rows []model.Subscription

	if err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id ASC").
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("subscription list by user: %w", err)
	}

	res := make([]*domain.Subscription, 0, len(rows))
	for i := range rows {
		res = append(res, rows[i].ToDomain())
	}

	return res, nil
}
