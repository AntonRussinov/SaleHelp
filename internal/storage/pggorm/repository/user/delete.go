package user

import (
	"context"
	"fmt"
)

// DeleteByID permanently deletes a user and all cascaded data.
func (r *Repository) DeleteByID(ctx context.Context, id int) error {
	res := r.DB.WithContext(ctx).Exec("DELETE FROM users WHERE id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("user delete by id: %w", res.Error)
	}
	return nil
}
