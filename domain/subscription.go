package domain

import "time"

const (
	TierStart       = "start"
	TierBusiness    = "business"
	TierPro         = "pro"
	TierTrial       = "trial"
	TierPackageOnly = "package_only"
)

const (
	SubStatusActive    = "active"
	SubStatusTrial     = "trial"
	SubStatusCancelled = "cancelled"
	SubStatusExpired   = "expired"
)

// TierConfig holds per-tier feature limits.
type TierConfig struct {
	AIQuotaMonthly int
	MaxShops       int
	MaxAIRules     int
	ToneOverride   bool // per-category tone (Pro only)
	ModelOverride  bool // per-category model selection (Pro only)
	ExtraShopPrice int  // monthly extra-shop price in kopecks
}

// TierConfigs defines limits for each subscription tier.
var TierConfigs = map[string]TierConfig{
	TierPackageOnly: {AIQuotaMonthly: 0, MaxShops: 1, MaxAIRules: 0},
	TierStart:       {AIQuotaMonthly: 50, MaxShops: 1, MaxAIRules: 0, ExtraShopPrice: 99000},
	TierBusiness:    {AIQuotaMonthly: 200, MaxShops: 2, MaxAIRules: 2, ExtraShopPrice: 149000},
	TierPro:         {AIQuotaMonthly: 800, MaxShops: 5, MaxAIRules: 5, ToneOverride: true, ModelOverride: true, ExtraShopPrice: 199000},
	TierTrial:       {AIQuotaMonthly: 200, MaxShops: 2, MaxAIRules: 2}, // same as business
}

type Subscription struct {
	ID           int
	UserID       int
	Tier         string
	Status       string
	PeriodMonths int
	StartedAt    time.Time
	ExpiresAt    time.Time
	TrialEndsAt  *time.Time
	AIQuota      int
	AIUsed       int
	AIResetsAt   time.Time
	MaxShops     int
	ExtraShops   int
	PriceMinor   int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TierCfg returns the TierConfig for this subscription's tier, falling back to package_only.
func (s *Subscription) TierCfg() TierConfig {
	if cfg, ok := TierConfigs[s.Tier]; ok {
		return cfg
	}
	return TierConfigs[TierPackageOnly]
}

// AIRemaining returns how many AI replies are left in the subscription quota this month.
func (s *Subscription) AIRemaining() int {
	if r := s.AIQuota - s.AIUsed; r > 0 {
		return r
	}
	return 0
}

// AIPackage is a purchased bundle of additional AI replies.
type AIPackage struct {
	ID          int
	UserID      int
	Size        int
	Remaining   int
	PriceMinor  int
	PurchasedAt time.Time
}
