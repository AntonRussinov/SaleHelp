package notifications

import (
	"context"
	"salehelp/domain"
	"time"
)

type Service struct {
	dispatcher *Dispatcher
	interval   time.Duration
}

type DeliveriesRepo interface {
	LeaseUnsent(ctx context.Context, workerID string, limit int, leaseSeconds int) ([]*domain.AlertTriggeredDelivery, error)
	MarkSent(ctx context.Context, id int, sentAt time.Time) error
	MarkFailed(ctx context.Context, id int, errText string) error
	MarkSentMany(ctx context.Context, ids []int, sentAt time.Time) error
	MarkFailedMany(ctx context.Context, ids []int, errText string) error
}

type FanoutRepo interface {
	CreateDeliveriesForTriggered(ctx context.Context, triggeredID int, shopID int) error
	// CreateDeliveriesForTriggeredCount(ctx context.Context, triggeredID int, shopID int) (created int64, err error)
}

type TriggeredRepo interface {
	// нужен, чтобы взять данные события и сформировать сообщение
	GetByID(ctx context.Context, id int) (*domain.AlertTriggered, error)
	GetByIDs(ctx context.Context, ids []int) (map[int]*domain.AlertTriggered, error)
}

type ProductsRepo interface {
	GetByIDs(ctx context.Context, ids []int) (map[int]*domain.Product, error)
}

type ProductSizesRepo interface {
	GetByIDs(ctx context.Context, ids []int) (map[int]*domain.ProductSize, error)
	ListByProductIDs(ctx context.Context, productIDs []int) ([]*domain.ProductSize, error)
}

type EndpointsRepo interface {
	GetByID(ctx context.Context, id int) (*domain.NotificationEndpoint, error)
}

func NewService(dispatcher *Dispatcher, interval time.Duration) *Service {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	return &Service{
		dispatcher: dispatcher,
		interval:   interval,
	}
}

func (s *Service) Run(ctx context.Context) error {
	t := time.NewTicker(s.interval)
	defer t.Stop()

	for {
		// запускаем сразу, а потом по тикам
		if err := s.dispatcher.RunOnce(ctx); err != nil {
			// можно залогировать и продолжать, но лучше вернуть ошибку наверх и решить на уровне main
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}
	}
}
