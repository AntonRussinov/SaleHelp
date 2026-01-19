package subscription

import (
	"context"
	"fmt"
)

func (s *Service) GetSubscriptionStatus(ctx context.Context, shopID int) (*StatusResult, error) {
	if shopID <= 0 {
		return nil, fmt.Errorf("shopID is required")
	}

	now := s.now()

	active, err := s.subscriptionsRepo.GetActiveByShopID(ctx, shopID, now)
	if err != nil {
		return nil, err
	}

	allowed, err := s.AllowedTelegramChats(ctx, shopID)
	if err != nil {
		return nil, err
	}

	enabled, err := s.endpoints.CountEnabledByShopAndChannel(ctx, shopID, "telegram")
	if err != nil {
		return nil, err
	}

	res := &StatusResult{
		HasActiveSubscription: active != nil,
		AllowedTelegramChats:  allowed,
		EnabledTelegramChats:  enabled,
	}

	if active == nil {
		// подписки нет — вернём только entitlements/usage
		return res, nil
	}

	plan := active.Plan
	status := active.Status
	validFrom := active.ValidFrom
	validTo := active.ValidTo

	res.Plan = &plan
	res.Status = &status
	res.ValidFrom = &validFrom
	res.ValidTo = &validTo

	res.CurrencyCode = &active.CurrencyCode
	res.BaseAmountMinor = &active.BaseAmountMinor
	res.PeriodMonths = &active.PeriodMonths

	return res, nil
}
