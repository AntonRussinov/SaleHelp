package subscription

import "time"

// GrantInput — вход для Service.GrantSubscription (не домен, а handler-level DTO).
type GrantInput struct {
	ShopID int

	PeriodMonths int    `json:"period_months" example:"1"` // 1|3|6
	CurrencyCode string `json:"currency_code" example:"RUB"` // "RUB"|"KZT"|...

	// База за период (в minor units: копейки/тиыны)
	BaseAmountMinor int64 `json:"base_amount_minor" example:"49900"`

	// Сколько доп. телеграм чатов купить (помимо 1 базового)
	ExtraTelegramChats int `json:"extra_telegram_chats" example:"0"`

	// Цена аддона за весь период (minor units).
	// Например: 2 доп чата на 1 месяц * 500₽ = 100000 (копеек)
	TelegramAddonAmountMinor int64 `json:"telegram_addon_amount_minor" example:"0"`

	PaidAt  *time.Time `json:"paid_at"  example:"2024-01-15T10:00:00Z"` // опционально
	Source  string     `json:"source"   example:"manual"`               // manual|...
	Comment *string    `json:"comment"  example:"Оплата по счёту №123"` // опционально
}

type GrantResult struct {
	SubscriptionID       int       `json:"subscription_id"        example:"7"`
	ValidFrom            time.Time `json:"valid_from"             example:"2024-01-15T00:00:00Z"`
	ValidTo              time.Time `json:"valid_to"               example:"2024-02-15T00:00:00Z"`
	AllowedTelegramChats int       `json:"allowed_telegram_chats" example:"1"`
}

// StatusResult — то, что удобно смотреть в Postman.
type StatusResult struct {
	HasActiveSubscription bool       `json:"has_active_subscription" example:"true"`
	Plan                  *string    `json:"plan,omitempty"          example:"base"`
	Status                *string    `json:"status,omitempty"        example:"active"`
	ValidFrom             *time.Time `json:"valid_from,omitempty"    example:"2024-01-15T00:00:00Z"`
	ValidTo               *time.Time `json:"valid_to,omitempty"      example:"2024-02-15T00:00:00Z"`

	CurrencyCode    *string `json:"currency_code,omitempty"    example:"RUB"`
	BaseAmountMinor *int64  `json:"base_amount_minor,omitempty" example:"49900"`
	PeriodMonths    *int    `json:"period_months,omitempty"    example:"1"`

	AllowedTelegramChats int `json:"allowed_telegram_chats" example:"1"`
	EnabledTelegramChats int `json:"enabled_telegram_chats" example:"1"`
}
