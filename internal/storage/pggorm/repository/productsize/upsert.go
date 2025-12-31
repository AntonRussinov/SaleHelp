package productsize

import (
	"context"
	"fmt"
	"salehelp/domain"
	"time"

	"gorm.io/gorm/clause"
)

// Upsert по UNIQUE(product_id, external_id)
func (r *Repository) Upsert(ctx context.Context, d *domain.ProductSize) error {
	m := toModel(d)
	now := time.Now()

	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now

	if err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "product_id"},
				{Name: "external_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"tech_size_name": m.TechSizeName,
				"updated_at":     now,
			}),
		}).
		Clauses(clause.Returning{}).
		Create(m).Error; err != nil {

		return fmt.Errorf("product_size upsert: %w", err)
	}

	d.ID = m.ID
	d.CreatedAt = m.CreatedAt
	d.UpdatedAt = m.UpdatedAt
	return nil
}
