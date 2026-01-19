package domain

import "time"

type ProductSnapshot struct {
	ID           int
	ProductID    int
	SnapshotTime time.Time

	// цены и валюта
	CurrencyCode *string
	BasePrice    *float64
	FinalPrice   *float64
	ClubPrice    *float64
	Discount     *float64
	ClubDiscount *float64
	OldPrice     *float64

	// остатки и продажи
	StockTotal   *int
	StockInWay   *int
	Rating       *float64
	ReviewsCount *int
	Sales24h     *int
	Sales7d      *int

	PositionSearch *int

	IsBadTurnover     *bool
	EditableSizePrice *bool

	Extra     []byte
	CreatedAt time.Time
}
