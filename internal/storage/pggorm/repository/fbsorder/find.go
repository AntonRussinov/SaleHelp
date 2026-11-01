package fbsorder

import (
	"context"
	"fmt"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)



// CountNewByShopGroupedByNmID возвращает количество активных FBS-заказов (supplier_status='new')
// по каждому nmID для магазина. Используется для раннего обнаружения нулевого остатка.
func (r *Repository) CountNewByShopGroupedByNmID(ctx context.Context, shopID int) (map[int64]int, error) {
	type row struct {
		NmID  int64 `gorm:"column:nm_id"`
		Count int   `gorm:"column:cnt"`
	}
	var rows []row
	err := r.DB.WithContext(ctx).
		Table("fbs_orders").
		Select("nm_id, COUNT(*) AS cnt").
		Where("shop_id = ? AND supplier_status = 'new'", shopID).
		Group("nm_id").
		Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("fbs_orders count new by nm_id: %w", err)
	}
	out := make(map[int64]int, len(rows))
	for _, row := range rows {
		out[row.NmID] = row.Count
	}
	return out, nil
}

// ListUnnotifiedByShop возвращает заказы магазина, по которым ещё не отправлено уведомление (notified_at IS NULL).
func (r *Repository) ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.FBSOrder, error) {
	if limit <= 0 {
		limit = 100
	}
	var rows []model.FBSOrder
	err := r.DB.WithContext(ctx).
		Where("shop_id = ? AND notified_at IS NULL", shopID).
		Order("created_at ASC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("fbs_orders list unnotified: %w", err)
	}
	out := make([]*domain.FBSOrder, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}
	return out, nil
}
