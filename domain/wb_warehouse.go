package domain

import "time"

type WBWarehouse struct {
	ID              int
	WBID            int64
	Name            string
	Address         *string
	City            *string
	Longitude       *float64
	Latitude        *float64
	CargoType       *int
	DeliveryType    *int
	FederalDistrict *string
	Selected        *bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
