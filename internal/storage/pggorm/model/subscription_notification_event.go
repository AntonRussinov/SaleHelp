package model

import (
	"salehelp/domain"
	"time"
)

type SubscriptionNotificationEvent struct {
	ID             int       `gorm:"primaryKey;column:id"`
	SubscriptionID int       `gorm:"column:subscription_id;not null"`
	Event          string    `gorm:"column:event;type:text;not null"`
	SentAt         time.Time `gorm:"column:sent_at;autoCreateTime"`
}

func (*SubscriptionNotificationEvent) TableName() string {
	return "subscription_notification_events"
}

func SubscriptionNotificationEventColumns() []string {
	return []string{
		"id",
		"subscription_id",
		"event",
		"sent_at",
	}
}

func (e *SubscriptionNotificationEvent) ToDomain() *domain.SubscriptionNotificationEvent {
	if e == nil {
		return nil
	}
	return &domain.SubscriptionNotificationEvent{
		ID:             e.ID,
		SubscriptionID: e.SubscriptionID,
		Event:          e.Event,
		SentAt:         e.SentAt,
	}
}
