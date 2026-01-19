package notificationendpoint

import (
	"context"
	"errors"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"strconv"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.NotificationEndpoint, error) {
	var m model.ShopNotificationEndpoint

	if err := r.DB.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("notification_endpoint get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []int) (map[int]*domain.NotificationEndpoint, error) {
	var rows []model.ShopNotificationEndpoint

	if err := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("notification_endpoint get by ids: %w", err)
	}

	out := make(map[int]*domain.NotificationEndpoint, len(rows))
	for i := range rows {
		d := rows[i].ToDomain()
		out[d.ID] = d
	}
	return out, nil
}

func (r *Repository) ListByShopID(ctx context.Context, shopID int, onlyEnabled bool) ([]*domain.NotificationEndpoint, error) {
	var rows []model.ShopNotificationEndpoint

	q := r.DB.WithContext(ctx).Where("shop_id = ?", shopID)
	if onlyEnabled {
		q = q.Where("enabled = TRUE")
	}

	if err := q.Order("id").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("notification_endpoint list by shop: %w", err)
	}

	out := make([]*domain.NotificationEndpoint, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}

func (r *Repository) Exists(ctx context.Context, shopID int, channel, address string) (bool, error) {
	var count int64

	err := r.DB.WithContext(ctx).
		Model(&model.ShopNotificationEndpoint{}).
		Where("shop_id = ?", shopID).
		Where("channel = ?", channel).
		Where("address = ?", address).
		Where("enabled = true").
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("shop_notification_endpoints exists: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) IsOwner(ctx context.Context, shopID int, channel, address string) (bool, error) {
	var m model.ShopNotificationEndpoint

	err := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Where("channel = ?", channel).
		Where("address = ?", address).
		Where("enabled = true").
		Where("label = ?", "owner").
		First(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("shop_notification_endpoints is owner: %w", err)
	}

	return true, nil
}

func (r *Repository) ListByChat(ctx context.Context, channel, address string) ([]domain.NotificationEndpoint, error) {
	var rows []model.ShopNotificationEndpoint

	if err := r.DB.WithContext(ctx).
		Where("channel = ?", channel).
		Where("address = ?", address).
		Where("enabled = true").
		Order("created_at ASC").
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("shop_notification_endpoints list by chat: %w", err)
	}

	out := make([]domain.NotificationEndpoint, 0, len(rows))
	for _, r := range rows {
		out = append(out, *r.ToDomain())
	}
	return out, nil
}

func (r *Repository) GetByShopAndAddress(
	ctx context.Context,
	shopID int,
	channel string,
	address string,
) (*domain.NotificationEndpoint, error) {

	var m model.ShopNotificationEndpoint

	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND channel = ? AND address = ?", shopID, channel, address).
		First(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("endpoint get by shop+address: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListEnabledTelegramChatIDsByShopID(ctx context.Context, shopID int) ([]int64, error) {
	var addresses []string

	err := r.DB.WithContext(ctx).
		Model(&model.ShopNotificationEndpoint{}).
		Select("address").
		Where("shop_id = ? AND channel = ? AND enabled = TRUE", shopID, "telegram").
		Find(&addresses).Error
	if err != nil {
		return nil, fmt.Errorf("list enabled telegram endpoints: %w", err)
	}

	res := make([]int64, 0, len(addresses))
	for _, a := range addresses {
		v, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			// v0: пропускаем мусор, но можно логировать
			continue
		}
		res = append(res, v)
	}

	return res, nil
}
