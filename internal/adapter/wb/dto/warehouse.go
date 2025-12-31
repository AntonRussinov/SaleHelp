package dto

// Склады WB (глобальные)
type WBWarehouse struct {
	Address         string  `json:"address"`
	Name            string  `json:"name"`
	City            string  `json:"city"`
	ID              int64   `json:"id"`
	Longitude       float64 `json:"longitude"`
	Latitude        float64 `json:"latitude"`
	CargoType       int     `json:"cargoType"`
	DeliveryType    int     `json:"deliveryType"`
	FederalDistrict string  `json:"federalDistrict"`
	Selected        bool    `json:"selected"`
}

// Склады продавца (FBS и т.п.)
type SellerWarehouse struct {
	Name         string `json:"name"`
	OfficeID     int64  `json:"officeId"`
	ID           int64  `json:"id"`
	CargoType    int    `json:"cargoType"`
	DeliveryType int    `json:"deliveryType"`
	IsDeleting   bool   `json:"isDeleting"`
	IsProcessing bool   `json:"isProcessing"`
}
