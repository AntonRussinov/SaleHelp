package domain

import "time"

type ProductSizeSKU struct {
	ID            int
	ProductSizeID int
	SKU           string

	CreatedAt time.Time
}
