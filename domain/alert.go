package domain

import "time"

// Alert types (shared across API/service/storage).
const (
	AlertTypePriceDrop = "price_drop"
	AlertTypeStockLow  = "stock_low"
	AlertTypeStockZero = "stock_zero"
)

// Alert scopes.
const (
	AlertScopeAuto    = "auto"
	AlertScopeProduct = "product"
	AlertScopeSize    = "size"
)

// Alert statuses.
const (
	AlertStatusActive   = "active"
	AlertStatusInactive = "inactive"
)

type Alert struct {
	ID                 int
	ShopID             int
	ProductID          *int
	Type               string
	Severity           int
	Status             string
	Message            string
	Payload            []byte
	TriggeredAt        time.Time
	ResolvedAt         *time.Time
	Scope              string // "auto" | "product" | "size"
	ProductSizeID      *int   // nullable
	AIRecommendationID *int
}
