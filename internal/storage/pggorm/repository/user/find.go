package user

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	var m model.User

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("user get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m model.User

	if err := r.DB.WithContext(ctx).
		First(&m, "email = ?", email).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("user get by email: %w", err)
	}

	return m.ToDomain(), nil
}
