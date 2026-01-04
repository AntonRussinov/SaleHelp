package shopinvitetokens

import (
	"context"
	"errors"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetValidByToken(ctx context.Context, token string, now time.Time) (*domain.ShopInviteToken, error) {
	var m model.ShopInviteToken

	err := r.DB.WithContext(ctx).
		Where("token = ?", token).
		Where("used_at IS NULL").
		Where("expires_at > ?", now).
		First(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("shop_invite_tokens get valid by token: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) GetActiveByShopID(
	ctx context.Context,
	shopID int,
) (*domain.ShopInviteToken, error) {

	var m model.ShopInviteToken

	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND expires_at > now()", shopID).
		Order("expires_at DESC").
		First(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("invite get active by shop: %w", err)
	}

	return m.ToDomain(), nil
}
