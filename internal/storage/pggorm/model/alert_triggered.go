package model

import (
	"salehelp/domain"
	"time"
)

type AlertTriggered struct {
	ID         int       `gorm:"primaryKey;column:id"`
	AlertID    int       `gorm:"column:alert_id;not null"`
	ProductID  int       `gorm:"column:product_id;not null"`
	SnapshotID int       `gorm:"column:snapshot_id"`
	Message    string    `gorm:"column:message;type:text"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func AlertTriggeredColumns() []string {
	return []string{
		"id",
		"alert_id",
		"product_id",
		"snapshot_id",
		"message",
		"created_at",
	}
}

func (at *AlertTriggered) ToDomain() *domain.AlertTriggered {
	if at == nil {
		return nil
	}

	return &domain.AlertTriggered{
		ID:         at.ID,
		AlertID:    at.AlertID,
		ProductID:  at.ProductID,
		SnapshotID: at.SnapshotID,
		Message:    at.Message,
		CreatedAt:  at.CreatedAt,
	}
}

func (*AlertTriggered) TableName() string {
	return "alerts_triggered"
}
