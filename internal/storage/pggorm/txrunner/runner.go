package txrunner

import (
	"context"
	"salehelp/internal/services/subscription"
	subRepoPkg "salehelp/internal/storage/pggorm/repository/subscription"
	addRepoPkg "salehelp/internal/storage/pggorm/repository/subscription_addon"

	"gorm.io/gorm"
)

type Runner struct {
	db      *gorm.DB
	subRepo *subRepoPkg.Repository
	addRepo *addRepoPkg.Repository
}

func New(db *gorm.DB, subRepo *subRepoPkg.Repository, addRepo *addRepoPkg.Repository) *Runner {
	return &Runner{db: db, subRepo: subRepo, addRepo: addRepo}
}

type deps struct {
	sub subscription.SubscriptionsRepo
	add subscription.AddonsRepo
}

func (d deps) Subscriptions() subscription.SubscriptionsRepo { return d.sub }
func (d deps) Addons() subscription.AddonsRepo               { return d.add }

func (r *Runner) WithinTx(ctx context.Context, fn func(ctx context.Context, deps subscription.TxDeps) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		d := deps{
			sub: r.subRepo.WithDB(tx),
			add: r.addRepo.WithDB(tx),
		}
		return fn(ctx, d)
	})
}
