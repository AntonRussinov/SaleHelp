package shopabcsettings

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

// Upsert создаёт или обновляет настройки ABC для магазина.
func (r *Repository) Upsert(ctx context.Context, s *domain.ShopABCSettings) error {
	if s == nil {
		return nil
	}
	now := time.Now()
	m := &model.ShopABCSettings{
		ShopID:            s.ShopID,
		ABCWeeklyEnabled:  s.ABCWeeklyEnabled,
		ABCMonthlyEnabled: s.ABCMonthlyEnabled,
		ABCAlertsEnabled:  s.ABCAlertsEnabled,
		UpdatedAt:         now,
	}
	err := r.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "shop_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"abc_weekly_enabled", "abc_monthly_enabled", "abc_alerts_enabled", "updated_at"}),
	}).Create(m).Error
	if err != nil {
		return fmt.Errorf("shop_abc_settings upsert: %w", err)
	}
	return nil
}
