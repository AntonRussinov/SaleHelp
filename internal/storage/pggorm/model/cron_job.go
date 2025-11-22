package model

import (
	"salehelp/domain"
	"time"
)

type CronJob struct {
	ID           int        `gorm:"primaryKey;column:id"`
	Name         string     `gorm:"column:name;type:varchar(100);not null"`
	Status       string     `gorm:"column:status;type:varchar(20);not null"` // STARTED, SUCCESS, ERROR
	ErrorMessage *string    `gorm:"column:error_message;type:text"`
	StartedAt    time.Time  `gorm:"column:started_at;not null"`
	FinishedAt   *time.Time `gorm:"column:finished_at"`
}

func CronJobColumns() []string {
	return []string{
		"id",
		"name",
		"status",
		"error_message",
		"started_at",
		"finished_at",
	}
}

func (cj *CronJob) ToDomain() *domain.CronJob {
	if cj == nil {
		return nil
	}

	return &domain.CronJob{
		ID:           cj.ID,
		Name:         cj.Name,
		Status:       cj.Status,
		ErrorMessage: cj.ErrorMessage,
		StartedAt:    cj.StartedAt,
		FinishedAt:   cj.FinishedAt,
	}
}

func (*CronJob) TableName() string {
	return "cron_jobs"
}
