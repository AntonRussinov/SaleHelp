package linktokens

import (
	"context"
	"errors"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"time"

	"gorm.io/gorm"
)

// GetValidByToken возвращает токен, если он существует и не истёк. Токен не удаляется — потребление отдельно.
func (r *Repository) GetValidByToken(ctx context.Context, token string, now time.Time) (*domain.LinkToken, error) {
	var m model.LinkToken
	err := r.DB.WithContext(ctx).
		Where("token = ?", token).
		Where("expires_at > ?", now).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("link_tokens get by token: %w", err)
	}
	return m.ToDomain(), nil
}
