package alert

import (
	"context"
	"fmt"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) Delete(ctx context.Context, id int, shopID int) error {
	result := r.DB.WithContext(ctx).
		Delete(&model.Alert{}, "id = ? AND shop_id = ?", id, shopID)

	if result.Error != nil {
		return fmt.Errorf("alert delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}