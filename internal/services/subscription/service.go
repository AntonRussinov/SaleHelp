package subscription

import (
	"context"
	"salehelp/domain"
	"time"
)

type SubscriptionsRepo interface {
	Create(ctx context.Context, s *domain.Subscription) error
	GetActiveByUserID(ctx context.Context, userID int) (*domain.Subscription, error)
	GetLatestByUserID(ctx context.Context, userID int) (*domain.Subscription, error)
	ExpireActiveByUserID(ctx context.Context, userID int) error
	ListExpiringBetween(ctx context.Context, from, to time.Time) ([]*domain.Subscription, error)
	IncrementAIUsed(ctx context.Context, userID int) (bool, error)
	ResetMonthlyQuota(ctx context.Context, now time.Time) error
}

type AIPackagesRepo interface {
	ConsumeOne(ctx context.Context, userID int) (packageID int, ok bool, err error)
	TotalRemaining(ctx context.Context, userID int) (int, error)
}

// ShopsRepo allows resolving shopID → userID.
type ShopsRepo interface {
	GetByID(ctx context.Context, id int) (*domain.Shop, error)
}

type Service struct {
	subscriptionsRepo SubscriptionsRepo
	aiPackagesRepo    AIPackagesRepo
	shopsRepo         ShopsRepo
	now               func() time.Time
}

func New(
	subscriptionsRepo SubscriptionsRepo,
	aiPackagesRepo AIPackagesRepo,
	shopsRepo ShopsRepo,
) *Service {
	return &Service{
		subscriptionsRepo: subscriptionsRepo,
		aiPackagesRepo:    aiPackagesRepo,
		shopsRepo:         shopsRepo,
		now:               time.Now,
	}
}

func (s *Service) SetNow(fn func() time.Time) { s.now = fn }

// userIDByShopID resolves shopID → userID via shops repo.
func (s *Service) userIDByShopID(ctx context.Context, shopID int) (int, error) {
	shop, err := s.shopsRepo.GetByID(ctx, shopID)
	if err != nil {
		return 0, err
	}
	if shop == nil {
		return 0, nil
	}
	return shop.UserID, nil
}
