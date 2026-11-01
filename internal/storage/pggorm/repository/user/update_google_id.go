package user

import (
	"context"
	"fmt"

	"salehelp/internal/storage/pggorm/model"
)

// SetGoogleID привязывает google_id к пользователю (например при первом входе по Google для user, созданного ранее по email/Telegram).
func (r *Repository) SetGoogleID(ctx context.Context, userID int, googleID string) error {
	if userID <= 0 || googleID == "" {
		return fmt.Errorf("user_id and google_id are required")
	}
	res := r.DB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("google_id", googleID)
	if res.Error != nil {
		return fmt.Errorf("user set google_id: %w", res.Error)
	}
	return nil
}
