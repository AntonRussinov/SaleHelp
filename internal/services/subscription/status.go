package subscription

import (
	"context"
	"fmt"
)

func (s *Service) GetSubscriptionStatus(ctx context.Context, shopID int) (*StatusResult, error) {
	if shopID <= 0 {
		return nil, fmt.Errorf("shopID is required")
	}

	userID, err := s.userIDByShopID(ctx, shopID)
	if err != nil {
		return nil, err
	}
	if userID == 0 {
		return &StatusResult{HasActiveSubscription: false}, nil
	}

	sub, err := s.subscriptionsRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := &StatusResult{
		HasActiveSubscription: sub != nil,
		AllowedTelegramChats:  1,
	}

	if sub != nil {
		tier := sub.Tier
		status := sub.Status
		res.Tier = &tier
		res.Status = &status
		res.ExpiresAt = &sub.ExpiresAt
		res.AIQuota = sub.AIQuota
		res.AIUsed = sub.AIUsed
		res.MaxShops = sub.MaxShops
		res.AllowedTelegramChats = 3
	}

	return res, nil
}
