package domain

import "time"

type CronJob struct {
	ID           int
	Name         string
	Status       string
	ErrorMessage *string
	StartedAt    time.Time
	FinishedAt   *time.Time
}
