package abcsnapshot

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
)

// GetLatestByShopAndPeriod возвращает последний по created_at снэпшот для магазина и периода (все строки).
func (r *Repository) GetLatestByShopAndPeriod(ctx context.Context, shopID, periodDays int) (*domain.ABCSnapshot, error) {
	var createdAt time.Time
	err := r.DB.WithContext(ctx).
		Model(&model.ABCSnapshotRow{}).
		Select("created_at").
		Where("shop_id = ? AND period_days = ?", shopID, periodDays).
		Order("created_at DESC").
		Limit(1).
		Pluck("created_at", &createdAt).Error
	if err != nil {
		return nil, fmt.Errorf("abc_snapshots get latest created_at: %w", err)
	}
	if createdAt.IsZero() {
		return nil, nil
	}
	var allRows []model.ABCSnapshotRow
	err = r.DB.WithContext(ctx).
		Where("shop_id = ? AND period_days = ? AND created_at = ?", shopID, periodDays, createdAt).
		Order("revenue DESC").
		Find(&allRows).Error
	if err != nil {
		return nil, fmt.Errorf("abc_snapshots get latest rows: %w", err)
	}
	if len(allRows) == 0 {
		return nil, nil
	}
	out := &domain.ABCSnapshot{
		PeriodDays:  periodDays,
		PeriodStart: allRows[0].PeriodStart,
		PeriodEnd:   allRows[0].PeriodEnd,
		CreatedAt:   createdAt,
		Rows:        make([]domain.ABCSnapshotRow, 0, len(allRows)),
	}
	for _, m := range allRows {
		out.Rows = append(out.Rows, *m.ToDomain())
		out.TotalRevenue += m.Revenue
	}
	return out, nil
}

// ListTwoLatestByShopAndPeriod возвращает два последних снэпшота (по дате created_at) для сравнения.
// Первый элемент — самый новый, второй — предыдущий.
func (r *Repository) ListTwoLatestByShopAndPeriod(ctx context.Context, shopID, periodDays int) ([]*domain.ABCSnapshot, error) {
	type createdGroup struct {
		CreatedAt time.Time `gorm:"column:created_at"`
	}
	var groups []createdGroup
	err := r.DB.WithContext(ctx).
		Model(&model.ABCSnapshotRow{}).
		Select("DISTINCT created_at").
		Where("shop_id = ? AND period_days = ?", shopID, periodDays).
		Order("created_at DESC").
		Limit(2).
		Scan(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("abc_snapshots list two latest groups: %w", err)
	}
	if len(groups) == 0 {
		return nil, nil
	}
	result := make([]*domain.ABCSnapshot, 0, len(groups))
	for _, g := range groups {
		var rows []model.ABCSnapshotRow
		err = r.DB.WithContext(ctx).
			Where("shop_id = ? AND period_days = ? AND created_at = ?", shopID, periodDays, g.CreatedAt).
			Order("revenue DESC").
			Find(&rows).Error
		if err != nil {
			return nil, fmt.Errorf("abc_snapshots get rows for created_at %v: %w", g.CreatedAt, err)
		}
		if len(rows) == 0 {
			continue
		}
		snap := &domain.ABCSnapshot{
			PeriodDays:  periodDays,
			PeriodStart: rows[0].PeriodStart,
			PeriodEnd:   rows[0].PeriodEnd,
			CreatedAt:   g.CreatedAt,
			Rows:        make([]domain.ABCSnapshotRow, 0, len(rows)),
		}
		for _, m := range rows {
			snap.Rows = append(snap.Rows, *m.ToDomain())
			snap.TotalRevenue += m.Revenue
		}
		result = append(result, snap)
	}
	return result, nil
}

// ListHistory возвращает последние N снэпшотов (по created_at) для магазина и периода.
func (r *Repository) ListHistory(ctx context.Context, shopID, periodDays, limit int) ([]*domain.ABCSnapshot, error) {
	if limit <= 0 {
		limit = 12
	}
	type createdGroup struct {
		CreatedAt time.Time `gorm:"column:created_at"`
	}
	var groups []createdGroup
	err := r.DB.WithContext(ctx).
		Model(&model.ABCSnapshotRow{}).
		Select("DISTINCT created_at").
		Where("shop_id = ? AND period_days = ?", shopID, periodDays).
		Order("created_at DESC").
		Limit(limit).
		Scan(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("abc_snapshots list history groups: %w", err)
	}
	result := make([]*domain.ABCSnapshot, 0, len(groups))
	for _, g := range groups {
		var rows []model.ABCSnapshotRow
		err = r.DB.WithContext(ctx).
			Where("shop_id = ? AND period_days = ? AND created_at = ?", shopID, periodDays, g.CreatedAt).
			Order("revenue DESC").
			Find(&rows).Error
		if err != nil {
			return nil, fmt.Errorf("abc_snapshots get rows for history: %w", err)
		}
		if len(rows) == 0 {
			continue
		}
		snap := &domain.ABCSnapshot{
			PeriodDays:  periodDays,
			PeriodStart: rows[0].PeriodStart,
			PeriodEnd:   rows[0].PeriodEnd,
			CreatedAt:   g.CreatedAt,
			Rows:        make([]domain.ABCSnapshotRow, 0, len(rows)),
		}
		for _, m := range rows {
			snap.Rows = append(snap.Rows, *m.ToDomain())
			snap.TotalRevenue += m.Revenue
		}
		result = append(result, snap)
	}
	return result, nil
}
