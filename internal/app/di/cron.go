package di

import (
	"context"
	"errors"
	"fmt"
	"salehelp/internal/storage/cron"
)

var ErrInitCronRunner = errors.New("failed to initialize cron runner")

func (prov *Provider) InitCronRunner(ctx context.Context) error {
	cron, err := cron.NewCronRunner(ctx, prov.Database.GormDB, prov.log)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInitCronRunner, err)
	}

	prov.Cron = cron

	return nil
}
