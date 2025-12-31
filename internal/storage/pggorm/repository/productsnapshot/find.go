package productsnapshot

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"github.com/jackc/pgtype"
)

func (r *Repository) GetLastByProductID(ctx context.Context, productID int) (*domain.ProductSnapshot, error) {
	var m model.ProductSnapshot

	err := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("snapshot_time DESC").
		Limit(1).
		Find(&m).Error

	if err != nil {
		return nil, fmt.Errorf("get last snapshot: %w", err)
	}

	if m.ID == 0 {
		return nil, nil
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListByProductID(ctx context.Context, productID int, limit, offset int) ([]*domain.ProductSnapshot, error) {
	var rows []model.ProductSnapshot

	q := r.DB.WithContext(ctx).Where("product_id = ?", productID).Order("snapshot_time DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}

	out := make([]*domain.ProductSnapshot, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}

// GetLatestByProductID возвращает последний snapshot по snapshot_time для 1 продукта.
func (r *Repository) GetLatestByProductID(ctx context.Context, productID int) (*domain.ProductSnapshot, error) {
	var m model.ProductSnapshot

	err := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("snapshot_time DESC").
		Limit(1).
		Find(&m).Error
	if err != nil {
		return nil, fmt.Errorf("product_snapshot get latest by product_id: %w", err)
	}

	// Find + Limit(1) вернёт пустую структуру без ErrRecordNotFound,
	// поэтому проверяем PK.
	if m.ID == 0 {
		return nil, nil
	}

	return m.ToDomain(), nil
}

// GetLatestByProductIDs возвращает последние snapshots для набора productIDs.
// Возвращает map[productID]*ProductSnapshot
func (r *Repository) GetLatestByProductIDs(ctx context.Context, productIDs []int) (map[int]*domain.ProductSnapshot, error) {
	out := make(map[int]*domain.ProductSnapshot)

	if len(productIDs) == 0 {
		return out, nil
	}

	// В Postgres самый простой и быстрый способ: DISTINCT ON (product_id)
	// Он вернёт по одной строке на product_id с максимальным snapshot_time.
	var rows []model.ProductSnapshot

	if err := r.DB.WithContext(ctx).
		Raw(`
			SELECT DISTINCT ON (product_id)
				*
			FROM product_snapshots
			WHERE product_id = ANY (?)
			ORDER BY product_id, snapshot_time DESC
		`, productIDs).
		Scan(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_snapshot get latest by product_ids: %w", err)
	}

	for i := range rows {
		m := rows[i]
		d := m.ToDomain()
		out[m.ProductID] = d
	}

	return out, nil
}

// GetLatestByProductID возвращает последний snapshot по snapshot_time для 1 продукта.
func (r *Repository) GetLatestByProductIDLimit(ctx context.Context, productID int, limit int) (*domain.ProductSnapshot, error) {
	var m model.ProductSnapshot

	err := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("snapshot_time DESC").
		Limit(limit).
		Find(&m).Error
	if err != nil {
		return nil, fmt.Errorf("product_snapshot get latest by product_id: %w", err)
	}

	// Find + Limit(1) вернёт пустую структуру без ErrRecordNotFound,
	// поэтому проверяем PK.
	if m.ID == 0 {
		return nil, nil
	}

	return m.ToDomain(), nil
}

// GetLatestByProductIDs возвращает последние snapshots для набора productIDs.
// Возвращает map[productID]*ProductSnapshot
func (r *Repository) GetLatestByProductIDsLimit(
	ctx context.Context,
	productIDs []int,
	limit int,
) (map[int][]*domain.ProductSnapshot, error) {

	out := make(map[int][]*domain.ProductSnapshot)

	if len(productIDs) == 0 || limit <= 0 {
		return out, nil
	}

	var rows []model.ProductSnapshot

	// ВАЖНО: добавь "id DESC" как tie-breaker, чтобы порядок был стабильный
	// (если snapshot_time одинаковый).
	q := `
		WITH ranked AS (
			SELECT
				ps.*,
				ROW_NUMBER() OVER (
					PARTITION BY ps.product_id
					ORDER BY ps.snapshot_time DESC, ps.id DESC
				) AS rn
			FROM product_snapshots ps
			WHERE ps.product_id = ANY (?)
		)
		SELECT *
		FROM ranked
		WHERE rn <= ?
		ORDER BY product_id, rn;
	`

	arr := pgtype.Int8Array{}
	_ = arr.Set(productIDs) // или конвертнуть в []int64

	if err := r.DB.WithContext(ctx).
		Raw(q, arr, limit).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("product_snapshot get latest(limit) by product_ids: %w", err)
	}

	for i := range rows {
		m := rows[i]
		out[m.ProductID] = append(out[m.ProductID], m.ToDomain()) // rn=1 идёт первым
	}

	return out, nil
}
