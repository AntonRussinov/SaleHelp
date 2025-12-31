package wbwarehouse

import (
	"context"
	"fmt"
	"salehelp/domain"

	"gorm.io/gorm/clause"
)

func (r *Repository) Update(ctx context.Context, d *domain.WBWarehouse) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("wb_warehouse update: %w", err)
	}

	d.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *Repository) Upsert(ctx context.Context, d *domain.WBWarehouse) error {
	m := toModel(d)

	if err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "wb_id"}}, // UNIQUE (wb_id)
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"address",
				"city",
				"longitude",
				"latitude",
				"cargo_type",
				"delivery_type",
				"federal_district",
				"selected",
				"updated_at",
			}),
		}).
		Create(m).Error; err != nil {

		return fmt.Errorf("wb_warehouse upsert: %w", err)
	}

	// На случай, если вернётся ID/таймштампы
	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt

	return nil
}
