package di

import (
	"context"
	"errors"
	"salehelp/internal/cron"
)

var ErrInitCronRunner = errors.New("failed to initialize cron runner")

func (prov *Provider) InitCronRunner(ctx context.Context) error {

	cron.GetOrders(ctx, prov.Database.GormDB, prov.log)

	// cron, err := cron.NewCronRunner(ctx, prov.Database.GormDB, prov.log)
	// if err != nil {
	// 	return fmt.Errorf("%w: %w", ErrInitCronRunner, err)
	// }

	// prov.Cron = cron

	return nil
}
