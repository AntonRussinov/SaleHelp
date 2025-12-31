package model

import (
	"salehelp/domain"
	"time"
)

type AlertTriggeredDelivery struct {
	ID          int `gorm:"primaryKey;column:id"`
	TriggeredID int `gorm:"column:triggered_id;not null"`
	EndpointID  int `gorm:"column:endpoint_id;not null"`

	Status    string     `gorm:"column:status;not null;default:new"`
	SentAt    *time.Time `gorm:"column:sent_at"`
	Attempts  int        `gorm:"column:attempts;not null;default:0"`
	LastError *string    `gorm:"column:last_error"`

	LeaseUntil *time.Time `gorm:"column:lease_until"`
	LeasedBy   *string    `gorm:"column:leased_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func AlertTriggeredDeliveryColumns() []string {
	return []string{
		"id",
		"triggered_id",
		"endpoint_id",
		"status",
		"sent_at",
		"attempts",
		"last_error",
		"lease_until",
		"leased_by",
		"created_at",
		"updated_at",
	}
}

func (m *AlertTriggeredDelivery) ToDomain() *domain.AlertTriggeredDelivery {
	if m == nil {
		return nil
	}

	return &domain.AlertTriggeredDelivery{
		ID:          m.ID,
		TriggeredID: m.TriggeredID,
		EndpointID:  m.EndpointID,
		Status:      m.Status,
		SentAt:      m.SentAt,
		Attempts:    m.Attempts,
		LastError:   m.LastError,
		LeaseUntil:  m.LeaseUntil,
		LeasedBy:    m.LeasedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (*AlertTriggeredDelivery) TableName() string {
	return "alerts_triggered_deliveries"
}
