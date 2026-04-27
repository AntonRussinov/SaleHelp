package api

import "time"

// CreateShopRequest тело POST /shops.
type CreateShopRequest struct {
	Name             string `json:"name"             validate:"required,max=255"           example:"Мой магазин WB"`
	APIToken         string `json:"apiToken"         validate:"required"                   example:"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJleGFtcGxlIn0.ZXhhbXBsZQ"`
	Marketplace      string `json:"marketplace"      validate:"required,oneof=wildberries" example:"wildberries"`
	WBReportTimezone string `json:"wbReportTimezone"                                       example:"Europe/Moscow"`
}

// UpdateShopRequest тело PATCH /shop/:id.
type UpdateShopRequest struct {
	WBReportTimezone *string `json:"wbReportTimezone" example:"Asia/Almaty"`
}

// SubscriptionSummary краткая информация о подписке магазина.
type SubscriptionSummary struct {
	HasActive bool       `json:"hasActive"        example:"true"`
	Plan      *string    `json:"plan,omitempty"   example:"base"`
	ValidTo   *time.Time `json:"validTo,omitempty" example:"2025-01-15T00:00:00Z"`
}

// ShopResponse один магазин в ответе.
type ShopResponse struct {
	ID                 int                  `json:"id"                          example:"1"`
	UserID             int                  `json:"userId"                      example:"5"`
	Marketplace        string               `json:"marketplace"                 example:"wildberries"`
	Name               string               `json:"name"                        example:"Мой магазин WB"`
	IsActive           bool                 `json:"isActive"                    example:"true"`
	TokenType          string               `json:"tokenType"                   example:"basic"`
	WBReportTimezone   string               `json:"wbReportTimezone,omitempty"  example:"Europe/Moscow"`
	LastSyncAt         *time.Time           `json:"lastSyncAt,omitempty"`
	LastSyncStatusCode int                  `json:"lastSyncStatusCode,omitempty" example:"200"`
	CreatedAt          time.Time            `json:"createdAt"                   example:"2024-01-15T10:00:00Z"`
	Subscription       *SubscriptionSummary `json:"subscription"`
}
