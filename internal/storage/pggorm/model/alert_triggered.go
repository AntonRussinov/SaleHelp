package model

import (
	"salehelp/domain"
	"time"
)

type AlertTriggered struct {
	ID            int       `gorm:"primaryKey;column:id"`
	AlertID       int       `gorm:"column:alert_id;not null"`
	ProductID     int       `gorm:"column:product_id;not null"`
	SnapshotID    int       `gorm:"column:snapshot_id"`
	ShopID        int       `gorm:"column:shop_id;not null"`
	Message       string    `gorm:"column:message;type:text"`
	ProductSizeID *int      // nullable
	AlertType     string    `gorm:"column:alert_type;type:text;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
}

func AlertTriggeredColumns() []string {
	return []string{
		"id",
		"alert_id",
		"product_id",
		"snapshot_id",
		"shop_id",
		"message",
		"product_size_id",
		"alert_type",
		"created_at",
	}
}

func (at *AlertTriggered) ToDomain() *domain.AlertTriggered {
	if at == nil {
		return nil
	}

	return &domain.AlertTriggered{
		ID:            at.ID,
		AlertID:       at.AlertID,
		ProductID:     at.ProductID,
		SnapshotID:    at.SnapshotID,
		ShopID:        at.ShopID,
		Message:       at.Message,
		ProductSizeID: at.ProductSizeID,
		AlertType:     at.AlertType,
		CreatedAt:     at.CreatedAt,
	}
}
func (*AlertTriggered) TableName() string {
	return "alerts_triggered"
}
