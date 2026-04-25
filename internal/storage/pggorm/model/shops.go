package model

import (
	"salehelp/domain"
	"time"
)

type Shop struct {
	ID                  int        `gorm:"primaryKey;column:id"`
	UserID              int        `gorm:"column:user_id;not null"`
	Marketplace         string     `gorm:"column:marketplace;type:text;not null"` // "WB", "OZON", "Kaspi", ...
	ExternalID          string     `gorm:"column:external_id;type:text"`
	Name                string     `gorm:"column:name;type:text;not null"`
	APIToken            string     `gorm:"column:api_token;type:text;not null"`
	TokenType           string     `gorm:"column:token_type;type:text;not null;default:'basic'"`
	IsActive            bool       `gorm:"column:is_active;not null;default:true"`
	WBReportTimezone    string     `gorm:"column:wb_report_timezone;type:text"`
	LastSyncAt          *time.Time `gorm:"column:last_sync_at"`
	FirstSyncNotifiedAt *time.Time `gorm:"column:first_sync_notified_at"`
	LastSyncStatusCode  int        `gorm:"column:last_sync_status_code"`
	LastSyncError       string     `gorm:"column:last_sync_error;type:text"`
	LastSyncFailedAt    *time.Time `gorm:"column:last_sync_failed_at"`
	LastSyncOkAt        *time.Time `gorm:"column:last_sync_ok_at"`
	LastHeavySyncAt     *time.Time `gorm:"column:last_heavy_sync_at"`
	CardsCursorUpdatedAt string    `gorm:"column:cards_cursor_updated_at;type:text"`
	CardsCursorNmID      int       `gorm:"column:cards_cursor_nmid"`
	LastFullCardsSyncAt  *time.Time `gorm:"column:last_full_cards_sync_at"`
	CreatedAt           time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

type SyncStatusUpdate struct {
	LastSyncAt         *time.Time
	LastSyncStatusCode *int
	LastSyncError      *string
	LastSyncFailedAt   *time.Time
	LastSyncOkAt       *time.Time
}

// HeavySyncUpdate содержит поля для обновления после тяжёлого синка (карточки, цены, остатки).
type HeavySyncUpdate struct {
	LastHeavySyncAt      *time.Time
	CardsCursorUpdatedAt *string
	CardsCursorNmID      *int
	LastFullCardsSyncAt  *time.Time
}

type MoneyByCurrency struct {
	CurrencyCode string `gorm:"column:currency_code"`
	AmountMinor  int64  `gorm:"column:amount_minor"`
	Count        int64  `gorm:"column:cnt"`
}

func ShopColumns() []string {
	return []string{
		"id",
		"user_id",
		"marketplace",
		"external_id",
		"name",
		"api_token",
		"token_type",
		"is_active",
		"wb_report_timezone",
		"last_sync_at",
		"first_sync_notified_at",
		"last_sync_status_code",
		"last_sync_error",
		"last_sync_failed_at",
		"last_sync_ok_at",
		"last_heavy_sync_at",
		"cards_cursor_updated_at",
		"cards_cursor_nmid",
		"last_full_cards_sync_at",
		"created_at",
		"updated_at",
	}
}

func (s *Shop) ToDomain() *domain.Shop {
	if s == nil {
		return nil
	}

	return &domain.Shop{
		ID:                   s.ID,
		UserID:               s.UserID,
		Marketplace:          s.Marketplace,
		ExternalID:           s.ExternalID,
		Name:                 s.Name,
		APIToken:             s.APIToken,
		TokenType:            s.TokenType,
		IsActive:             s.IsActive,
		WBReportTimezone:     s.WBReportTimezone,
		LastSyncAt:           s.LastSyncAt,
		FirstSyncNotifiedAt:  s.FirstSyncNotifiedAt,
		LastSyncStatusCode:   s.LastSyncStatusCode,
		LastSyncError:        s.LastSyncError,
		LastSyncFailedAt:     s.LastSyncFailedAt,
		LastSyncOkAt:         s.LastSyncOkAt,
		LastHeavySyncAt:      s.LastHeavySyncAt,
		CardsCursorUpdatedAt: s.CardsCursorUpdatedAt,
		CardsCursorNmID:      s.CardsCursorNmID,
		LastFullCardsSyncAt:  s.LastFullCardsSyncAt,
		CreatedAt:            s.CreatedAt,
		UpdatedAt:            s.UpdatedAt,
	}
}

func (*Shop) TableName() string {
	return "shops"
}
