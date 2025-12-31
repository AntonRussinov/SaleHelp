package domain

import "time"

type SellerWarehouse struct {
	ID           int
	ShopID       int
	WBOfficeID   int64
	WBID         int64
	Name         string
	CargoType    *int
	DeliveryType *int
	IsDeleting   bool
	IsProcessing bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
