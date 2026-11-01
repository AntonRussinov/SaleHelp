package domain

import "time"

// ABCCategory — категория ABC-анализа (A/B/C).
type ABCCategory string

const (
	ABCCategoryA ABCCategory = "A"
	ABCCategoryB ABCCategory = "B"
	ABCCategoryC ABCCategory = "C"
)

// ABCSnapshotRow — одна строка снэпшота ABC (один товар в одном срезе).
type ABCSnapshotRow struct {
	ID         int
	ShopID     int
	ProductID  *int
	NmID       int64
	Category   ABCCategory
	Revenue    float64
	SalesCount int
	PeriodDays int
	PeriodStart time.Time
	PeriodEnd   time.Time
	CreatedAt   time.Time
}

// ABCSnapshot — снэпшот ABC-анализа за период (набор строк по товарам).
type ABCSnapshot struct {
	Rows       []ABCSnapshotRow
	PeriodDays int
	PeriodStart time.Time
	PeriodEnd   time.Time
	CreatedAt   time.Time
	TotalRevenue float64
}

// ABCChange — изменение категории товара между двумя снэпшотами.
type ABCChange struct {
	NmID       int64
	ProductID  *int
	Subject    string
	FromCat    ABCCategory
	ToCat      ABCCategory
	Revenue    float64
	ChangePct  float64
}

// ShopABCSettings — настройки ABC для магазина.
type ShopABCSettings struct {
	ShopID             int
	ABCWeeklyEnabled   bool
	ABCMonthlyEnabled  bool
	ABCAlertsEnabled   bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
