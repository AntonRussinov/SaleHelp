package subscription

import (
	"context"
	"fmt"
	"salehelp/domain"
)

func (s *Service) GrantSubscription(ctx context.Context, in GrantInput) (*GrantResult, error) {
	if in.ShopID <= 0 {
		return nil, fmt.Errorf("shop_id is required")
	}
	if in.PeriodMonths != 1 && in.PeriodMonths != 3 && in.PeriodMonths != 6 {
		return nil, fmt.Errorf("invalid period_months: %d", in.PeriodMonths)
	}
	if len(in.CurrencyCode) != 3 {
		return nil, fmt.Errorf("currency_code must be 3 letters, got %q", in.CurrencyCode)
	}
	if in.BaseAmountMinor < 0 || in.TelegramAddonAmountMinor < 0 {
		return nil, fmt.Errorf("amount_minor must be >= 0")
	}
	if in.ExtraTelegramChats < 0 {
		return nil, fmt.Errorf("extra telegram chats must be >= 0")
	}
	if in.Source == "" {
		in.Source = "manual"
	}

	now := s.now()

	var out *GrantResult

	err := s.tx.WithinTx(ctx, func(ctx context.Context, deps TxDeps) error {
		subRepo := deps.Subscriptions()
		addRepo := deps.Addons()

		active, err := subRepo.GetActiveByShopID(ctx, in.ShopID, now)
		if err != nil {
			return err
		}

		validFrom := now
		if active != nil && active.ValidTo.After(validFrom) {
			validFrom = active.ValidTo
		}
		validTo := validFrom.AddDate(0, in.PeriodMonths, 0)

		if _, err := subRepo.ExpireActiveByShopID(ctx, in.ShopID); err != nil {
			return err
		}

		sub := domain.Subscription{
			ShopID:          in.ShopID,
			Plan:            PlanBasic,
			Status:          StatusActive,
			ValidFrom:       validFrom,
			ValidTo:         validTo,
			PeriodMonths:    in.PeriodMonths,
			BaseAmountMinor: in.BaseAmountMinor,
			CurrencyCode:    in.CurrencyCode,
			PaidAt:          in.PaidAt,
			Source:          in.Source,
			Comment:         in.Comment,
		}

		if err := subRepo.Create(ctx, &sub); err != nil {
			return err
		}

		if in.ExtraTelegramChats > 0 {
			add := domain.SubscriptionAddon{
				SubscriptionID: sub.ID,
				Code:           AddonTelegramExtraChat,
				Quantity:       in.ExtraTelegramChats,
				AmountMinor:    in.TelegramAddonAmountMinor,
				CurrencyCode:   in.CurrencyCode,
				Meta:           []byte(`{}`),
			}
			if err := addRepo.Upsert(ctx, &add); err != nil {
				return err
			}
		}

		out = &GrantResult{
			SubscriptionID:       sub.ID,
			ValidFrom:            validFrom,
			ValidTo:              validTo,
			AllowedTelegramChats: 1 + in.ExtraTelegramChats,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return out, nil
}
