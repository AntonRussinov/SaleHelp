package subscription

import "time"

// GrantInput — admin endpoint: grant a subscription to the user who owns the given shop.
type GrantInput struct {
	ShopID       int
	Tier         string  `json:"tier"          example:"business"` // start|business|pro
	PeriodMonths int     `json:"period_months" example:"1"`        // 1|3|6|12
	PriceMinor   int     `json:"price_minor"   example:"299000"`   // monthly price in kopecks
	Comment      *string `json:"comment"       example:"manual grant"`
}

// GrantResult — response from GrantSubscription.
type GrantResult struct {
	SubscriptionID int       `json:"subscription_id" example:"7"`
	UserID         int       `json:"user_id"         example:"3"`
	Tier           string    `json:"tier"            example:"business"`
	ExpiresAt      time.Time `json:"expires_at"      example:"2026-06-11T00:00:00Z"`
	AIQuota        int       `json:"ai_quota"        example:"200"`
	MaxShops       int       `json:"max_shops"       example:"2"`
}

// StatusResult — subscription status for display.
type StatusResult struct {
	HasActiveSubscription bool       `json:"has_active_subscription" example:"true"`
	Tier                  *string    `json:"tier,omitempty"          example:"business"`
	Status                *string    `json:"status,omitempty"        example:"active"`
	ExpiresAt             *time.Time `json:"expires_at,omitempty"    example:"2026-06-11T00:00:00Z"`
	AIQuota               int        `json:"ai_quota"                example:"200"`
	AIUsed                int        `json:"ai_used"                 example:"12"`
	MaxShops              int        `json:"max_shops"               example:"2"`
	AllowedTelegramChats  int        `json:"allowed_telegram_chats"  example:"3"`
}
