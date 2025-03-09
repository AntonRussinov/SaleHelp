package cron

import (
	"context"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type (
	logger interface {
		Info(msg string, args ...any)
		Error(msg string, args ...any)
	}
)

func NewCronRunner(ctx context.Context, db *gorm.DB, log logger) (*cron.Cron, error) {
	log.Info("Starting cron job for deleting old soft deleted records and updating old pending requests")

	crn := cron.New()

	// Start cron job every day at 00:00:00 (min, hour, day, month, day of week)
	// _, err := crn.AddFunc("*/5 * * * *", func() { GetOrders(ctx, db, log) })
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = crn.AddFunc("1 0 * * *", func() { CronUpdateRequest(db, log) })
	// if err != nil {
	// 	return nil, err
	// }

	return crn, nil
}
