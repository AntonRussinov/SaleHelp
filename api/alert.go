package api

import (
	"encoding/json"
	"fmt"
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

// CreateAlertRequest запрос на создание алерта.
// severity, scope и message вычисляются автоматически на бэкенде.
// Если productId и productSizeId не переданы — алерт работает по всем товарам магазина.
type CreateAlertRequest struct {
	ProductID     *int                   `json:"productId"     validate:"omitempty,min=1" example:"42"`
	ProductSizeID *int                   `json:"productSizeId" validate:"omitempty,min=1" example:"7"`
	Type          AlertType              `json:"type"          validate:"required,is_valid_enum" example:"price_drop"`
	Payload       map[string]interface{} `json:"payload"       validate:"required" swaggertype:"object"`
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
	ProductID *int         `query:"productId" validate:"omitempty,min=1"`
	Type      *AlertType   `query:"type" validate:"omitempty,is_valid_enum"`
	Status    *AlertStatus `query:"status" validate:"omitempty,is_valid_enum"`
	Limit     int          `query:"limit" validate:"min=1,max=100"`
	Offset    int          `query:"offset" validate:"min=0"`
}

// ListAlertsResponse ответ со списком алертов
type ListAlertsResponse struct {
	Alerts []AlertResponse `json:"alerts"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// PriceDropPayload конфигурация для price_drop алерта.
// Можно задать один или оба порога (percent, amount).
// mode: "any" — сработает при выполнении любого, "all" — оба сразу.
type PriceDropPayload struct {
	Percent float64 `json:"percent" validate:"min=0" example:"10"`
	Amount  int     `json:"amount"  validate:"min=0" example:"500"`
	Mode    string  `json:"mode"    validate:"required,oneof=any all" example:"any"`
}

// StockLowPayload конфигурация для stock_low алерта.
// threshold — минимальный остаток в штуках, при котором срабатывает алерт.
type StockLowPayload struct {
	Threshold int `json:"threshold" validate:"required,min=0" example:"10"`
}

// StockZeroPayload конфигурация для stock_zero алерта (payload не требуется)
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
		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return err
		}
		if payload.Percent == 0 && payload.Amount == 0 {
			return fmt.Errorf("price_drop: укажите хотя бы один порог (percent или amount)")
		}
		return nil
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

// AutoScope определяет scope по наличию productSizeId / productId.
func (r *CreateAlertRequest) AutoScope() string {
	if r.ProductSizeID != nil {
		return string(AlertScopeSize)
	}
	if r.ProductID != nil {
		return string(AlertScopeProduct)
	}
	return string(AlertScopeAuto)
}

// AutoSeverity вычисляет severity по типу и параметрам payload.
func (r *CreateAlertRequest) AutoSeverity() int {
	switch r.Type {
	case AlertTypeStockZero:
		return 5
	case AlertTypeStockLow:
		threshold := payloadInt(r.Payload, "threshold")
		switch {
		case threshold >= 50:
			return 4
		case threshold >= 10:
			return 3
		default:
			return 2
		}
	case AlertTypePriceDrop:
		percent := payloadFloat(r.Payload, "percent")
		amount := payloadInt(r.Payload, "amount")

		sevPercent := 0
		if percent >= 20 {
			sevPercent = 5
		} else if percent >= 10 {
			sevPercent = 4
		} else if percent >= 5 {
			sevPercent = 3
		} else if percent > 0 {
			sevPercent = 2
		}

		sevAmount := 0
		if amount >= 1000 {
			sevAmount = 4
		} else if amount >= 500 {
			sevAmount = 3
		} else if amount >= 100 {
			sevAmount = 2
		} else if amount > 0 {
			sevAmount = 1
		}

		if sevPercent > sevAmount {
			return sevPercent
		}
		if sevAmount > 0 {
			return sevAmount
		}
	}
	return 1
}

// AutoMessage генерирует человекочитаемое описание алерта на русском языке.
func (r *CreateAlertRequest) AutoMessage() string {
	switch r.Type {
	case AlertTypeStockZero:
		return "Товар закончился"
	case AlertTypeStockLow:
		threshold := payloadInt(r.Payload, "threshold")
		return fmt.Sprintf("Мало остатков (< %d шт.)", threshold)
	case AlertTypePriceDrop:
		percent := payloadFloat(r.Payload, "percent")
		amount := payloadInt(r.Payload, "amount")
		mode, _ := r.Payload["mode"].(string)

		if percent > 0 && amount > 0 {
			conjunction := "или"
			if mode == "all" {
				conjunction = "и"
			}
			return fmt.Sprintf("Снижение цены > %.0f%% %s %d руб.", percent, conjunction, amount)
		}
		if percent > 0 {
			return fmt.Sprintf("Снижение цены > %.0f%%", percent)
		}
		return fmt.Sprintf("Снижение цены > %d руб.", amount)
	}
	return ""
}

func payloadFloat(p map[string]interface{}, key string) float64 {
	v, _ := p[key].(float64)
	return v
}

func payloadInt(p map[string]interface{}, key string) int {
	switch v := p[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return 0
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
