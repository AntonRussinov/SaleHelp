package model

import (
	"salehelp/domain"
	"time"
)

type AIRecommendation struct {
	ID        int       `gorm:"primaryKey;column:id"`
	ShopID    int       `gorm:"column:shop_id;not null"`
	ProductID *int      `gorm:"column:product_id"` // NULLable
	AlertID   *int      `gorm:"column:alert_id"`   // NULLable
	Type      string    `gorm:"column:type;type:text;not null"`
	Model     string    `gorm:"column:model;type:text"`
	Prompt    string    `gorm:"column:prompt;type:text"`
	Content   string    `gorm:"column:content;type:text;not null"`
	Meta      []byte    `gorm:"column:meta;type:jsonb"` // raw JSON
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func AIRecommendationColumns() []string {
	return []string{
		"id",
		"shop_id",
		"product_id",
		"alert_id",
		"type",
		"model",
		"prompt",
		"content",
		"meta",
		"created_at",
	}
}

func (r *AIRecommendation) ToDomain() *domain.AIRecommendation {
	if r == nil {
		return nil
	}

	return &domain.AIRecommendation{
		ID:        r.ID,
		ShopID:    r.ShopID,
		ProductID: r.ProductID,
		AlertID:   r.AlertID,
		Type:      r.Type,
		Model:     r.Model,
		Prompt:    r.Prompt,
		Content:   r.Content,
		Meta:      r.Meta,
		CreatedAt: r.CreatedAt,
	}
}

func (*AIRecommendation) TableName() string {
	return "ai_recommendations"
}
