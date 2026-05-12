package subscription

import (
	"context"
	"fmt"
	"salehelp/domain"
	"time"
)

var validTiers = map[string]bool{
	domain.TierStart: true, domain.TierBusiness: true, domain.TierPro: true,
}

func (s *Service) GrantSubscription(ctx context.Context, in GrantInput) (*GrantResult, error) {
	if in.ShopID <= 0 {
		return nil, fmt.Errorf("shop_id is required")
	}
	if !validTiers[in.Tier] {
		return nil, fmt.Errorf("invalid tier %q: must be start|business|pro", in.Tier)
	}
	if in.PeriodMonths < 1 {
		return nil, fmt.Errorf("period_months must be >= 1")
	}

	userID, err := s.userIDByShopID(ctx, in.ShopID)
	if err != nil {
		return nil, fmt.Errorf("resolve user for shop %d: %w", in.ShopID, err)
	}
	if userID == 0 {
		return nil, fmt.Errorf("shop %d not found", in.ShopID)
	}

	cfg, ok := domain.TierConfigs[in.Tier]
	if !ok {
		return nil, fmt.Errorf("unknown tier %q", in.Tier)
	}

	now := s.now()
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	if err := s.subscriptionsRepo.ExpireActiveByUserID(ctx, userID); err != nil {
		return nil, fmt.Errorf("expire old subscription: %w", err)
	}

	sub := &domain.Subscription{
		UserID:       userID,
		Tier:         in.Tier,
		Status:       domain.SubStatusActive,
		PeriodMonths: in.PeriodMonths,
		StartedAt:    now,
		ExpiresAt:    now.AddDate(0, in.PeriodMonths, 0),
		AIQuota:      cfg.AIQuotaMonthly,
		AIUsed:       0,
		AIResetsAt:   nextMonth,
		MaxShops:     cfg.MaxShops,
		ExtraShops:   0,
		PriceMinor:   in.PriceMinor,
	}

	if err := s.subscriptionsRepo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	return &GrantResult{
		SubscriptionID: sub.ID,
		UserID:         userID,
		Tier:           sub.Tier,
		ExpiresAt:      sub.ExpiresAt,
		AIQuota:        sub.AIQuota,
		MaxShops:       sub.MaxShops,
	}, nil
}
