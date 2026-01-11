// internal/service/subscription/entitlements.go
package subscription

import "context"

func (s *Service) AllowedTelegramChats(ctx context.Context, shopID int) (int, error) {
	ent, err := s.GetCurrentEntitlements(ctx, shopID)
	if err != nil {
		return 0, err
	}
	return ent.AllowedTelegramChats, nil
}

type Entitlements struct {
	AllowedTelegramChats int
}

func (s *Service) GetCurrentEntitlements(ctx context.Context, shopID int) (*Entitlements, error) {
	now := s.now()

	active, err := s.subscriptionsRepo.GetActiveByShopID(ctx, shopID, now)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return &Entitlements{AllowedTelegramChats: 0}, nil
	}

	addons, err := s.addonsRepo.ListBySubscriptionID(ctx, active.ID)
	if err != nil {
		return nil, err
	}

	allowed := 1 // базовый чат
	for _, a := range addons {
		if a.Code == AddonTelegramExtraChat {
			allowed += a.Quantity
		}
	}

	return &Entitlements{AllowedTelegramChats: allowed}, nil
}
