package domain

import "time"

// OrderRefusal — запись об отказе (заказ не забран): заказ из WB Orders с isCancel=true.
// Связь с Sales по srid. Lifecycle: заказ создан (Orders isCancel=false) → 2a выкуп (Sales S...) / 2b отказ (Orders isCancel=true) / 3 возврат (Sales R...).
type OrderRefusal struct {
	ID              int
	ShopID          int
	Srid            string // ID заказа в WB, общий с Sales
	NmID            int64
	Barcode         string
	SupplierArticle string
	Date            time.Time // дата заказа
	CancelDate      time.Time // дата отмены/незабора
	LastChangeDate  time.Time
	CreatedAt       time.Time
	NotifiedAt      *time.Time
}
