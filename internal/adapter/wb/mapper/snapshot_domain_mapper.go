package mapper

import (
	"time"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

func ptr[T any](v T) *T { return &v }

// SnapshotWithSizes связывает снапшот товара и его размерные снапшоты.
// Это удобно, чтобы потом присвоить ProductSnapshotID после сохранения в БД.
type SnapshotWithSizes struct {
	Snapshot domain.ProductSnapshot
}

// ToSnapshots строит снапшоты по данным WB:
//   - prodByNm: nmID -> Product из БД
//   - prices: ответ /list/goods/filter
//   - stockByNm: nmID -> общий остаток по всем складам
func ToSnapshots(
	prodByNm map[int64]domain.Product,
	prices []dto.PriceGood,
	stockByNm map[int64]int,
) []SnapshotWithSizes {

	if prices == nil {
		return nil
	}

	now := time.Now()
	result := make([]SnapshotWithSizes, 0, len(prices))

	for _, goods := range prices {
		p, ok := prodByNm[goods.NmID]
		if !ok {
			continue
		}

		// агрегаты по товару: можешь выбрать любую стратегию
		// (первый размер / минимальный / средний). Пока оставим первый.
		var basePrice, finalPrice, clubPrice float64
		if len(goods.Sizes) > 0 {
			basePrice = goods.Sizes[0].Price
			finalPrice = goods.Sizes[0].DiscountedPrice
			clubPrice = goods.Sizes[0].ClubDiscountedPrice
		}

		var totalStockPtr *int
		if stockByNm != nil {
			totalStock := stockByNm[goods.NmID]
			totalStockPtr = ptr(totalStock)
		}

		snap := domain.ProductSnapshot{
			ProductID:    p.ID,
			SnapshotTime: now,

			CurrencyCode: ptr(goods.CurrencyISOCode4217),

			BasePrice:  ptr(basePrice),
			FinalPrice: ptr(finalPrice),
			ClubPrice:  ptr(clubPrice),

			Discount:     ptr(goods.Discount),
			ClubDiscount: ptr(goods.ClubDiscount),

			// Если stockByNm == nil, значит stocks недоступен (ошибка WB) и сток неизвестен.
			// Нельзя превращать "нет данных" в "0" — это ломает алерты.
			StockTotal: totalStockPtr,

			IsBadTurnover:     ptr(goods.IsBadTurnover),
			EditableSizePrice: ptr(goods.EditableSizePrice),
		}

		result = append(result, SnapshotWithSizes{
			Snapshot: snap,
		})
	}

	return result
}
