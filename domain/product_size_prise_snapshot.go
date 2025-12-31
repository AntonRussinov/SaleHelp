package domain

type ProductSizePriceSnapshot struct {
	ID                int
	ProductSnapshotID int

	ProductSizeID int

	BasePrice    *float64
	FinalPrice   *float64
	ClubPrice    *float64
	Discount     *float64
	ClubDiscount *float64

	Extra []byte
}
