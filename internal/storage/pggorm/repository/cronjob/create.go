package cronjob

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (r *Repository) Create(ctx context.Context, c *domain.CronJob) error {
	m := toModel(c)

	if err := r.DB.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("cron_job create: %w", err)
	}

	c.ID = m.ID
	c.StartedAt = m.StartedAt
	c.FinishedAt = m.FinishedAt

	return nil
}
