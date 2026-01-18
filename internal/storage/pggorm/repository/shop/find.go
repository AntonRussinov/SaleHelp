package shop

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"strconv"
	"time"

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

func (r *Repository) ListAll(ctx context.Context) ([]*domain.Shop, error) {
	var m []model.Shop

	if err := r.DB.WithContext(ctx).Find(&m).Error; err != nil {
		return nil, err
	}

	result := make([]*domain.Shop, 0, len(m))
	for i := range m {
		result = append(result, m[i].ToDomain())
	}

	return result, nil
}

func (r *Repository) ListActiveIDs(ctx context.Context, limit, offset int) ([]int, error) {
	var ids []int

	q := r.DB.WithContext(ctx).
		Model(&model.Shop{}).
		Select("id").
		Where("is_active = true").
		Order("id ASC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Pluck("id", &ids).Error; err != nil {
		return nil, fmt.Errorf("shop list active ids: %w", err)
	}

	return ids, nil
}

func (r *Repository) ListByTelegramChat(
	ctx context.Context,
	chatID int64,
) ([]*domain.Shop, error) {
	var models []model.Shop

	err := r.DB.WithContext(ctx).
		Joins(`
			JOIN shop_notification_endpoints e
			  ON e.shop_id = shops.id
		`).
		Where("e.channel = ?", "telegram").
		Where("e.address = ?", strconv.FormatInt(chatID, 10)).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	res := make([]*domain.Shop, 0, len(models))
	for i := range models {
		res = append(res, models[i].ToDomain())
	}
	return res, nil
}

func (r *Repository) GetLastSyncAtMap(ctx context.Context, shopIDs []int) (map[int]*time.Time, error) {
	out := make(map[int]*time.Time, len(shopIDs))
	if len(shopIDs) == 0 {
		return out, nil
	}

	type row struct {
		ID         int        `gorm:"column:id"`
		LastSyncAt *time.Time `gorm:"column:last_sync_at"`
	}

	var rows []row
	if err := r.DB.WithContext(ctx).
		Table("shops").
		Select("id, last_sync_at").
		Where("id IN ?", shopIDs).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("get shops.last_sync_at map: %w", err)
	}

	for _, r := range rows {
		out[r.ID] = r.LastSyncAt
	}

	return out, nil
}
