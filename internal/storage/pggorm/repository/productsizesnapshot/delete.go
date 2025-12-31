package productsizesnapshot

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"
)

func (r *Repository) DeleteBySnapshotID(ctx context.Context, snapshotID int) error {
	if err := r.DB.WithContext(ctx).
		Where("product_snapshot_id = ?", snapshotID).
		Delete(&model.ProductSizeSnapshot{}).Error; err != nil {

		return fmt.Errorf("product_size_snapshot delete by snapshot: %w", err)
	}

	return nil
}
