package model

import (
	"salehelp/domain"
	"time"
)

type PriceRestoreBatch struct {
	ID               int        `gorm:"primaryKey;column:id"`
	ShopID           int        `gorm:"column:shop_id"`
	RestoreStatus    string     `gorm:"column:restore_status;not null;default:'pending'"`
	RestoreStartedAt *time.Time `gorm:"column:restore_started_at"`
	CreatedAt        time.Time  `gorm:"column:created_at;not null;autoCreateTime"`
	RestoredAt       *time.Time `gorm:"column:restored_at"`
}

func (b *PriceRestoreBatch) ToDomain() *domain.PriceRestoreBatch {
	if b == nil {
		return nil
	}
	return &domain.PriceRestoreBatch{
		ID:               b.ID,
		ShopID:           b.ShopID,
		RestoreStatus:    b.RestoreStatus,
		RestoreStartedAt: b.RestoreStartedAt,
		CreatedAt:        b.CreatedAt,
		RestoredAt:       b.RestoredAt,
	}
}

func (*PriceRestoreBatch) TableName() string {
	return "price_restore_batches"
}
