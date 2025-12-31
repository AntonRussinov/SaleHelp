package notificationendpoint

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

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
