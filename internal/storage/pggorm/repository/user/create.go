package user

import (
	"context"
	"errors"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Create(ctx context.Context, u *domain.User) error {
	m := toModel(u)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("user create: %w", err)
	}

	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt

	return nil
}

func (r *Repository) EnsureByTelegramChatID(ctx context.Context, chatID int64) (*domain.User, error) {
	// 1) fast path
	u, err := r.GetByTelegramChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return u, nil
	}

	// 2) try create
	m := &model.User{
		Email:          nil, // ⚠️ если email nullable в БД — ЛУЧШЕ сделать *string и хранить nil
		PasswordHash:   nil,
		TelegramChatID: &chatID,
	}

	// ВАЖНО: если ты переведёшь email в nullable, переделай model.User.Email на *string.
	// Тогда тут будет Email:nil, PasswordHash:nil.

	err = r.DB.WithContext(ctx).Create(m).Error
	if err == nil {
		return m.ToDomain(), nil
	}

	// 3) если это дубль по telegram_chat_id — просто дочитаем
	if isUniqueViolation(err) {
		u2, err2 := r.GetByTelegramChatID(ctx, chatID)
		if err2 != nil {
			return nil, err2
		}
		if u2 != nil {
			return u2, nil
		}
	}

	return nil, fmt.Errorf("user ensure by telegram_chat_id %d: %w", chatID, err)
}

func isUniqueViolation(err error) bool {
	// pgx
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
