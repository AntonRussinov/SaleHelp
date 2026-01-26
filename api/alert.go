package api

import (
	"encoding/json"
	"salehelp/domain"
	"time"
)

// AlertType представляет тип алерта
type AlertType string

const (
	AlertTypePriceDrop AlertType = AlertType(domain.AlertTypePriceDrop)
	AlertTypeStockLow  AlertType = AlertType(domain.AlertTypeStockLow)
	AlertTypeStockZero AlertType = AlertType(domain.AlertTypeStockZero)
)

// AlertScope представляет область действия алерта
type AlertScope string

const (
	AlertScopeAuto    AlertScope = AlertScope(domain.AlertScopeAuto)
	AlertScopeProduct AlertScope = AlertScope(domain.AlertScopeProduct)
	AlertScopeSize    AlertScope = AlertScope(domain.AlertScopeSize)
)

// AlertStatus представляет статус алерта
type AlertStatus string

const (
	AlertStatusActive   AlertStatus = AlertStatus(domain.AlertStatusActive)
	AlertStatusInactive AlertStatus = AlertStatus(domain.AlertStatusInactive)
)

// CreateAlertRequest запрос на создание алерта
type CreateAlertRequest struct {
	ProductID     *int                   `json:"productId" validate:"omitempty,min=1"`
	ProductSizeID *int                   `json:"productSizeId" validate:"omitempty,min=1"` 
	Type          AlertType              `json:"type" validate:"required,is_valid_enum"`
	Severity      int                    `json:"severity" validate:"required,min=1,max=5"`
	Status        AlertStatus            `json:"status" validate:"required,is_valid_enum"`
	Message       string                 `json:"message" validate:"required,max=500"`
	Scope         AlertScope             `json:"scope" validate:"required,is_valid_enum"`
	Payload       map[string]interface{} `json:"payload" validate:"required"`
}

// UpdateAlertRequest запрос на обновление алерта
type UpdateAlertRequest struct {
	ProductID     *int                   `json:"productId" validate:"omitempty,min=1"`
	ProductSizeID *int                   `json:"productSizeId" validate:"omitempty,min=1"`
	Severity      *int                   `json:"severity" validate:"omitempty,min=1,max=5"`
	Status        *AlertStatus           `json:"status" validate:"omitempty,is_valid_enum"`
	Message       *string                `json:"message" validate:"omitempty,max=500"`
	Scope         *AlertScope            `json:"scope" validate:"omitempty,is_valid_enum"`
	Payload       map[string]interface{} `json:"payload" validate:"omitempty"`
}

// AlertResponse ответ с данными алерта
type AlertResponse struct {
	ID                 int                    `json:"id"`
	ShopID             int                    `json:"shopId"`
	ProductID          *int                   `json:"productId"`
	ProductSizeID      *int                   `json:"productSizeId"`
	Type               AlertType              `json:"type"`
	Severity           int                    `json:"severity"`
	Status             AlertStatus            `json:"status"`
	Message            string                 `json:"message"`
	Payload            map[string]interface{} `json:"payload"`
	Scope              AlertScope             `json:"scope"`
	TriggeredAt        time.Time              `json:"triggeredAt"`
	ResolvedAt         *time.Time             `json:"resolvedAt"`
	AIRecommendationID *int                   `json:"aiRecommendationId"`
}

// ListAlertsRequest запрос на получение списка алертов
type ListAlertsRequest struct {
	ProductID *int        `query:"productId" validate:"omitempty,min=1"`
	Type      *AlertType  `query:"type" validate:"omitempty,is_valid_enum"`
	Status    *AlertStatus `query:"status" validate:"omitempty,is_valid_enum"`
	Limit     int         `query:"limit" validate:"min=1,max=100"`
	Offset    int         `query:"offset" validate:"min=0"`
}

// ListAlertsResponse ответ со списком алертов
type ListAlertsResponse struct {
	Alerts []AlertResponse `json:"alerts"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// PriceDropPayload конфигурация для price_drop алерта
type PriceDropPayload struct {
	Threshold float64 `json:"threshold" validate:"required,min=0"`
	Mode      string  `json:"mode" validate:"required,oneof=percent amount"`
	Condition string  `json:"condition" validate:"required,oneof=any all"`
}

// StockLowPayload конфигурация для stock_low алерта
type StockLowPayload struct {
	Threshold int `json:"threshold" validate:"required,min=0"`
}

// StockZeroPayload конфигурация для stock_zero алерта (пустая структура)
type StockZeroPayload struct{}

// ValidatePayload валидирует payload в зависимости от типа алерта
func (r *CreateAlertRequest) ValidatePayload() error {
	payloadBytes, err := json.Marshal(r.Payload)
	if err != nil {
		return err
	}

	switch r.Type {
	case AlertTypePriceDrop:
		var payload PriceDropPayload
		return json.Unmarshal(payloadBytes, &payload)
	case AlertTypeStockLow:
		var payload StockLowPayload
		return json.Unmarshal(payloadBytes, &payload)
	case AlertTypeStockZero:
		var payload StockZeroPayload
		return json.Unmarshal(payloadBytes, &payload)
	default:
		return nil
	}
}

// IsValid методы для валидации enum'ов
func (a AlertType) IsValid() bool {
	switch a {
	case AlertTypePriceDrop, AlertTypeStockLow, AlertTypeStockZero:
		return true
	default:
		return false
	}
}

func (a AlertScope) IsValid() bool {
	switch a {
	case AlertScopeAuto, AlertScopeProduct, AlertScopeSize:
		return true
	default:
		return false
	}
}

func (a AlertStatus) IsValid() bool {
	switch a {
	case AlertStatusActive, AlertStatusInactive:
		return true
	default:
		return false
	}
}