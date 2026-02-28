package user

import (
	"context"
	"errors"
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

func (r *Repository) GetByTelegramChatID(ctx context.Context, chatID int64) (*domain.User, error) {
	var m model.User

	err := r.DB.WithContext(ctx).
		Model(&model.User{}).
		Where("telegram_chat_id = ?", chatID).
		First(&m).Error

	if err == nil {
		return m.ToDomain(), nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, fmt.Errorf("user get by telegram_chat_id %d: %w", chatID, err)
}

func (r *Repository) GetByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	if googleID == "" {
		return nil, nil
	}
	var m model.User
	err := r.DB.WithContext(ctx).
		Where("google_id = ?", googleID).
		First(&m).Error
	if err == nil {
		return m.ToDomain(), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, fmt.Errorf("user get by google_id: %w", err)
}
