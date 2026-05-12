package subscription

import (
	"context"
	"salehelp/domain"
	"time"
)

type AccessLevel string

const (
	AccessNone    AccessLevel = "none"    // no subscription ever
	AccessActive  AccessLevel = "active"  // expires_at > now
	AccessGrace   AccessLevel = "grace"   // expired within 3 days
	AccessBlocked AccessLevel = "blocked" // expired more than 3 days ago
)

// GetShopAccessLevel checks the subscription status for the user who owns the shop.
func (s *Service) GetShopAccessLevel(ctx context.Context, shopID int) (AccessLevel, *domain.Subscription, error) {
	userID, err := s.userIDByShopID(ctx, shopID)
	if err != nil || userID == 0 {
		return AccessNone, nil, err
	}
	return s.GetUserAccessLevel(ctx, userID)
}

// GetUserAccessLevel checks the subscription status for a user directly.
func (s *Service) GetUserAccessLevel(ctx context.Context, userID int) (AccessLevel, *domain.Subscription, error) {
	now := s.now()

	sub, err := s.subscriptionsRepo.GetLatestByUserID(ctx, userID)
	if err != nil {
		return AccessNone, nil, err
	}
	if sub == nil {
		return AccessNone, nil, nil
	}

	if sub.ExpiresAt.After(now) {
		return AccessActive, sub, nil
	}
	if now.Before(sub.ExpiresAt.Add(72 * time.Hour)) {
		return AccessGrace, sub, nil
	}
	return AccessBlocked, sub, nil
}
