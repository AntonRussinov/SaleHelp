package cronjob

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.CronJob, error) {
	var m model.CronJob

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("cron_job get by id: %w", err)
	}

	return m.ToDomain(), nil
}

// GetLastByName — последний запуск джобы с таким именем
func (r *Repository) GetLastByName(ctx context.Context, name string) (*domain.CronJob, error) {
	var m model.CronJob

	err := r.DB.WithContext(ctx).
		Where("name = ?", name).
		Order("started_at DESC").
		Limit(1).
		Find(&m).Error
	if err != nil {
		return nil, fmt.Errorf("cron_job get last by name: %w", err)
	}

	if m.ID == 0 {
		return nil, nil
	}

	return m.ToDomain(), nil
}

// ListByName — история запусков конкретной джобы
func (r *Repository) ListByName(ctx context.Context, name string, limit, offset int) ([]*domain.CronJob, error) {
	var rows []model.CronJob

	q := r.DB.WithContext(ctx).
		Where("name = ?", name).
		Order("started_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("cron_job list by name: %w", err)
	}

	out := make([]*domain.CronJob, 0, len(rows))
	for i := range rows {
		out = append(out, rows[i].ToDomain())
	}

	return out, nil
}
