package model

import "time"

// ShopSyncLease — техническая таблица для leasing синков магазина.
// Это НЕ доменная сущность, поэтому живёт только в storage/model.
type ShopSyncLease struct {
	ShopID     int        `gorm:"column:shop_id;primaryKey"`
	LeasedBy   *string    `gorm:"column:leased_by"`
	LeaseUntil *time.Time `gorm:"column:lease_until"`
}

func (ShopSyncLease) TableName() string { return "shop_sync_leases" }
