package alerttriggered

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm/clause"
)

func (r *Repository) Create(ctx context.Context, d *domain.AlertTriggered) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("alert_triggered create: %w", err)
	}

	d.ID = m.ID
	d.CreatedAt = m.CreatedAt

	return nil
}
func (r *Repository) CreateIfNotExists(
	ctx context.Context,
	at *domain.AlertTriggered,
) (id int, inserted bool, err error) {

	m := model.AlertTriggered{
		AlertID:    at.AlertID,
		ShopID:     at.ShopID,
		ProductID:  at.ProductID,
		SnapshotID: at.SnapshotID,
		Message:    at.Message,
	}

	tx := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "alert_id"},
				{Name: "product_id"},
				{Name: "snapshot_id"},
			},
			DoNothing: true,
		}).
		Create(&m)

	if tx.Error != nil {
		return 0, false, fmt.Errorf("alert_triggered create if not exists: %w", tx.Error)
	}

	// 1 -> вставили, 0 -> конфликт
	if tx.RowsAffected == 0 {
		return 0, false, nil
	}

	// inserted
	inserted = true

	// m.ID заполнен только при реальной вставке
	at.ID = int(m.ID) // подстрой если ID другой тип
	return at.ID, true, nil
}
