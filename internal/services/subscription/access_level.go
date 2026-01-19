package subscription

import (
	"context"
	"salehelp/domain"
	"time"
)

type AccessLevel string

const (
	AccessNone    AccessLevel = "none"    // никогда не было подписки
	AccessActive  AccessLevel = "active"  // valid_to > now
	AccessGrace   AccessLevel = "grace"   // now <= valid_to + 3d
	AccessBlocked AccessLevel = "blocked" // now > valid_to + 3d
)

func (s *Service) GetShopAccessLevel(ctx context.Context, shopID int) (AccessLevel, *domain.Subscription, error) {
	now := s.now()

	sub, err := s.subscriptionsRepo.GetLatestByShopID(ctx, shopID)
	if err != nil {
		return AccessNone, nil, err
	}
	if sub == nil {
		return AccessNone, nil, nil
	}

	if sub.ValidTo.After(now) {
		return AccessActive, sub, nil
	}

	graceUntil := sub.ValidTo.Add(72 * time.Hour) // 3 дня
	if !now.After(graceUntil) {                   // now <= graceUntil
		return AccessGrace, sub, nil
	}

	return AccessBlocked, sub, nil
}
