package model

import "time"

type OrderSyncCursor struct {
	ShopID         int       `gorm:"primaryKey;column:shop_id"`
	LastChangeDate time.Time `gorm:"column:last_change_date;not null"`
}

func (*OrderSyncCursor) TableName() string {
	return "order_sync_cursor"
}
