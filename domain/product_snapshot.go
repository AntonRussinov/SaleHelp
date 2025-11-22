package domain

import "time"

type ProductSnapshot struct {
	ID             int
	ProductID      int
	SnapshotTime   time.Time
	Price          float64
	OldPrice       float64
	Discount       float64
	StockTotal     int
	Rating         float64
	ReviewsCount   int
	Sales24h       int
	Sales7d        int
	PositionSearch int
	Extra          []byte
	CreatedAt      time.Time
}
