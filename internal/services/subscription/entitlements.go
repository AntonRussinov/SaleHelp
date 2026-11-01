package subscription

import (
	"context"
	"salehelp/domain"
)

type Entitlements struct {
	Tier                 string
	AIQuotaMonthly       int
	AIQuotaRemaining     int
	AIPackagesRemaining  int
	MaxShops             int
	ExtraShops           int
	MaxAIRules           int
	ToneOverride         bool
	ModelOverride        bool
	AllowedTelegramChats int
}

// GetShopEntitlements returns entitlements for the user who owns the given shop.
func (s *Service) GetShopEntitlements(ctx context.Context, shopID int) (*Entitlements, error) {
	userID, err := s.userIDByShopID(ctx, shopID)
	if err != nil || userID == 0 {
		return defaultEntitlements(), nil
	}
	return s.GetUserEntitlements(ctx, userID)
}

// GetUserEntitlements returns entitlements for a user.
func (s *Service) GetUserEntitlements(ctx context.Context, userID int) (*Entitlements, error) {
	sub, err := s.subscriptionsRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		// No active subscription — check AI packages only
		pkgRemaining, _ := s.aiPackagesRepo.TotalRemaining(ctx, userID)
		return &Entitlements{
			Tier:                domain.TierPackageOnly,
			AllowedTelegramChats: 1,
			AIPackagesRemaining:  pkgRemaining,
		}, nil
	}

	cfg := sub.TierCfg()
	pkgRemaining, _ := s.aiPackagesRepo.TotalRemaining(ctx, userID)

	return &Entitlements{
		Tier:                 sub.Tier,
		AIQuotaMonthly:       sub.AIQuota,
		AIQuotaRemaining:     sub.AIRemaining(),
		AIPackagesRemaining:  pkgRemaining,
		MaxShops:             sub.MaxShops,
		ExtraShops:           sub.ExtraShops,
		MaxAIRules:           cfg.MaxAIRules,
		ToneOverride:         cfg.ToneOverride,
		ModelOverride:        cfg.ModelOverride,
		AllowedTelegramChats: 3,
	}, nil
}

// AllowedTelegramChats returns max telegram chats for the shop's owner.
func (s *Service) AllowedTelegramChats(ctx context.Context, shopID int) (int, error) {
	ent, err := s.GetShopEntitlements(ctx, shopID)
	if err != nil {
		return 0, err
	}
	return ent.AllowedTelegramChats, nil
}

func defaultEntitlements() *Entitlements {
	return &Entitlements{
		Tier:                 domain.TierPackageOnly,
		AllowedTelegramChats: 1,
	}
}
