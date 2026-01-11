package subscription

import "time"

// GrantInput — вход для Service.GrantSubscription (не домен, а handler-level DTO).
type GrantInput struct {
	ShopID int

	PeriodMonths int    `json:"period_months"` // 1|3|6
	CurrencyCode string `json:"currency_code"` // "RUB"|"KZT"|...

	// База за период (в minor units: копейки/тиыны)
	BaseAmountMinor int64 `json:"base_amount_minor"`

	// Сколько доп. телеграм чатов купить (помимо 1 базового)
	ExtraTelegramChats int `json:"extra_telegram_chats"`

	// Цена аддона за весь период (minor units).
	// Например: 2 доп чата на 1 месяц * 500₽ = 100000 (копеек)
	TelegramAddonAmountMinor int64 `json:"telegram_addon_amount_minor"`

	PaidAt  *time.Time `json:"paid_at"` // опционально
	Source  string     `json:"source"`  // manual|...
	Comment *string    `json:"comment"` // опционально
}

type GrantResult struct {
	SubscriptionID int       `json:"subscription_id"`
	ValidFrom      time.Time `json:"valid_from"`
	ValidTo        time.Time `json:"valid_to"`

	AllowedTelegramChats int `json:"allowed_telegram_chats"`
}

// StatusResult — то, что удобно смотреть в Postman.
type StatusResult struct {
	HasActiveSubscription bool       `json:"has_active_subscription"`
	Plan                  *string    `json:"plan,omitempty"`
	Status                *string    `json:"status,omitempty"`
	ValidFrom             *time.Time `json:"valid_from,omitempty"`
	ValidTo               *time.Time `json:"valid_to,omitempty"`

	CurrencyCode    *string `json:"currency_code,omitempty"`
	BaseAmountMinor *int64  `json:"base_amount_minor,omitempty"`
	PeriodMonths    *int    `json:"period_months,omitempty"`

	AllowedTelegramChats int `json:"allowed_telegram_chats"`
	EnabledTelegramChats int `json:"enabled_telegram_chats"`
}
