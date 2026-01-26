package alert

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func (r *Repository) Delete(ctx context.Context, id int, shopID int) error {
	result := r.DB.WithContext(ctx).
		Where("id = ? AND shop_id = ?", id, shopID).
		Delete(&struct{}{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("alert delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}