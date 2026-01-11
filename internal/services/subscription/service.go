package subscription

import (
	"context"
	"salehelp/domain"
	"time"
)

const (
	PlanBasic              = "basic"
	StatusActive           = "active"
	AddonTelegramExtraChat = "telegram_extra_chat"
)

type SubscriptionsRepo interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetActiveByShopID(ctx context.Context, shopID int, now time.Time) (*domain.Subscription, error)
	GetLatestByShopID(ctx context.Context, shopID int) (*domain.Subscription, error)
	ExpireActiveByShopID(ctx context.Context, shopID int) (int64, error)
}

type AddonsRepo interface {
	Upsert(ctx context.Context, a *domain.SubscriptionAddon) error
	ListBySubscriptionID(ctx context.Context, subscriptionID int) ([]domain.SubscriptionAddon, error)
}

// TransactionRunner — абстракция транзакции (реализация живёт в storage).
// Внутри fn выдаём репозитории, уже привязанные к tx.
type TransactionRunner interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, deps TxDeps) error) error
}

type EndpointsRepo interface {
	CountEnabledByShopAndChannel(ctx context.Context, shopID int, channel string) (int, error)
}

type TxDeps interface {
	Subscriptions() SubscriptionsRepo
	Addons() AddonsRepo
}

type Service struct {
	tx                TransactionRunner
	subscriptionsRepo SubscriptionsRepo
	addonsRepo        AddonsRepo
	now               func() time.Time
	endpoints         EndpointsRepo
}

func New(tx TransactionRunner,
	SubscriptionsRepo SubscriptionsRepo,
	AddonsRepo AddonsRepo,
	EndpointsRepo EndpointsRepo,
) *Service {
	return &Service{
		tx:                tx,
		subscriptionsRepo: SubscriptionsRepo,
		addonsRepo:        AddonsRepo,
		now:               time.Now,
		endpoints:         EndpointsRepo,
	}
}

func (s *Service) SetNow(fn func() time.Time) { s.now = fn }

// type GrantInput struct {
// 	ShopID int

// 	PeriodMonths int // 1|3|6

// 	BaseAmountMinor int64
// 	CurrencyCode    string

// 	ExtraTelegramChats       int
// 	TelegramAddonAmountMinor int64

// 	PaidAt  *time.Time
// 	Source  string
// 	Comment *string
// }

// type GrantResult struct {
// 	SubscriptionID       int
// 	ValidFrom            time.Time
// 	ValidTo              time.Time
// 	AllowedTelegramChats int
// }
