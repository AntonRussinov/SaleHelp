package cronjob

import (
	"context"
	"fmt"
	"salehelp/domain"
)

// Update — обновить статус/ошибку/время завершения
func (r *Repository) Update(ctx context.Context, c *domain.CronJob) error {
	m := toModel(c)

	if err := r.DB.WithContext(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("cron_job update: %w", err)
	}

	c.FinishedAt = m.FinishedAt
	c.StartedAt = m.StartedAt

	return nil
}
