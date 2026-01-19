package productsizesnapshot

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.ProductSizeSnapshot, error) {
	var m model.ProductSizeSnapshot

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("product_size_snapshot get by id: %w", err)
	}

	return m.ToDomain(), nil
}

func (r *Repository) ListBySnapshotID(ctx context.Context, snapshotID int) ([]*domain.ProductSizeSnapshot, error) {
	var rows []model.ProductSizeSnapshot

	if err := r.DB.WithContext(ctx).
		Where("product_snapshot_id = ?", snapshotID).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size_snapshot list: %w", err)
	}

	out := make([]*domain.ProductSizeSnapshot, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.ToDomain())
	}
	return out, nil
}

func (r *Repository) ListBySnapshotIDs(ctx context.Context, snapshotIDs []int) ([]*domain.ProductSizeSnapshot, error) {
	if len(snapshotIDs) == 0 {
		return nil, nil
	}

	var rows []model.ProductSizeSnapshot

	if err := r.DB.WithContext(ctx).
		Where("product_snapshot_id IN ?", snapshotIDs).
		Find(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size_snapshot list by snapshot_ids: %w", err)
	}

	out := make([]*domain.ProductSizeSnapshot, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}

func (r *Repository) GetLatestBySizeIDsLimitMap(
	ctx context.Context,
	sizeIDs []int,
	limit int,
) (map[int][]*domain.ProductSizeSnapshot, error) {

	out := make(map[int][]*domain.ProductSizeSnapshot)

	if len(sizeIDs) == 0 || limit <= 0 {
		return out, nil
	}

	var rows []model.ProductSizeSnapshot

	q := `
        WITH ranked AS (
            SELECT
                pss.*,
                ROW_NUMBER() OVER (
                    PARTITION BY pss.product_size_id
                    ORDER BY ps.snapshot_time DESC, ps.id DESC, pss.id DESC
                ) AS rn
            FROM product_size_snapshots pss
            JOIN product_snapshots ps ON ps.id = pss.product_snapshot_id
            WHERE pss.product_size_id = ANY (?)
        )
        SELECT *
        FROM ranked
        WHERE rn <= ?
        ORDER BY product_size_id, rn;
    `

	arr := pgtype.Int8Array{}
	_ = arr.Set(sizeIDs)

	// ВАЖНО: как и раньше, чтобы ANY(?) корректно работал
	if err := r.DB.WithContext(ctx).
		Raw(q, arr, limit).
		Scan(&rows).Error; err != nil {

		return nil, fmt.Errorf("product_size_snapshot get latest(limit) by size_ids: %w", err)
	}

	for i := range rows {
		d := rows[i].ToDomain()
		out[d.ProductSizeID] = append(out[d.ProductSizeID], d) // rn=1 идёт первым
	}

	return out, nil
}
