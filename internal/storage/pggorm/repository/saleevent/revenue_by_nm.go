package saleevent

import (
	"context"
	"fmt"
	"time"

	"salehelp/internal/storage/pggorm/model"
)

// RevenueByNmIDRow — выручка и количество выкупов по nm_id за период (только event_type = 'buyout').
type RevenueByNmIDRow struct {
	NmID       int64
	Revenue    float64
	SalesCount int
	Subject    string
}

// RevenueAndCountByNmIDInPeriod возвращает по каждому nm_id сумму выручки (for_pay) и количество выкупов
// за период [from, to], только события с event_type = 'buyout'.
func (r *Repository) RevenueAndCountByNmIDInPeriod(
	ctx context.Context,
	shopID int,
	from, to time.Time,
) ([]RevenueByNmIDRow, error) {
	type row struct {
		NmID       int64   `gorm:"column:nm_id"`
		Revenue    float64 `gorm:"column:revenue"`
		SalesCount int     `gorm:"column:sales_count"`
		Subject    string  `gorm:"column:subject"`
	}
	var rows []row
	err := r.DB.WithContext(ctx).
		Model(&model.SaleEvent{}).
		Select("nm_id, SUM(for_pay) AS revenue, COUNT(*)::int AS sales_count, MAX(subject) AS subject").
		Where("shop_id = ? AND event_type = ? AND date >= ? AND date < ?", shopID, "buyout", from, to).
		Group("nm_id").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("sale_events revenue by nm_id: %w", err)
	}
	out := make([]RevenueByNmIDRow, 0, len(rows))
	for _, rr := range rows {
		out = append(out, RevenueByNmIDRow{NmID: rr.NmID, Revenue: rr.Revenue, SalesCount: rr.SalesCount, Subject: rr.Subject})
	}
	return out, nil
}
