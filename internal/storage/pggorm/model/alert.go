package model

import (
	"salehelp/domain"
	"time"
)

type Alert struct {
	ID                 int        `gorm:"primaryKey;column:id"`
	ShopID             int        `gorm:"column:shop_id;not null"`
	ProductID          *int       `gorm:"column:product_id"`
	Type               string     `gorm:"column:type;type:text;not null"`
	Severity           int        `gorm:"column:severity;type:smallint;not null;default:1"`
	Status             string     `gorm:"column:status;type:text;not null;default:'new'"`
	Message            string     `gorm:"column:message;type:text"`
	Payload            []byte     `gorm:"column:payload"` // JSONB
	TriggeredAt        time.Time  `gorm:"column:triggered_at;autoCreateTime"`
	ResolvedAt         *time.Time `gorm:"column:resolved_at"`
	AIRecommendationID *int       `gorm:"column:ai_recommendation_id"`
}

func AlertColumns() []string {
	return []string{
		"id",
		"shop_id",
		"product_id",
		"type",
		"severity",
		"status",
		"message",
		"payload",
		"triggered_at",
		"resolved_at",
		"ai_recommendation_id",
	}
}

func (a *Alert) ToDomain() *domain.Alert {
	if a == nil {
		return nil
	}

	return &domain.Alert{
		ID:                 a.ID,
		ShopID:             a.ShopID,
		ProductID:          a.ProductID,
		Type:               a.Type,
		Severity:           a.Severity,
		Status:             a.Status,
		Message:            a.Message,
		Payload:            a.Payload,
		TriggeredAt:        a.TriggeredAt,
		ResolvedAt:         a.ResolvedAt,
		AIRecommendationID: a.AIRecommendationID,
	}
}

func (*Alert) TableName() string {
	return "alerts"
}
