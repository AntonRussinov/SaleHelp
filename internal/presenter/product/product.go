package product

// import (
// 	"salehelp/domain"
// 	"salehelp/internal/presenter/product/view"
// )

type Presenter struct{}

func NewPresenter() *Presenter {
	return &Presenter{}
}

// func (p *Presenter) Product(prod *domain.Products) *view.ProductView {
// 	if prod == nil {
// 		return nil
// 	}

// 	offers := make([]view.OfferView, len(prod.Data))
// 	for i, off := range prod.Data {
// 		offers[i] = *p.Offer(&off)
// 	}

// 	return &view.ProductView{
// 		Data:  offers,
// 		Total: prod.Total,
// 	}
// }

// // Offer преобразует доменную модель предложения в OfferView.
// func (p *Presenter) Offer(off *domain.Offer) *view.OfferView {
// 	if off == nil {
// 		return nil
// 	}

// 	return &view.OfferView{
// 		MasterSku:               off.MasterSku,
// 		Images:                  off.Images,
// 		ShopLink:                off.ShopLink,
// 		MasterTitle:             off.MasterTitle,
// 		OfferId:                 off.OfferId,
// 		FileId:                  off.FileId,
// 		MerchantUid:             off.MerchantUid,
// 		Sku:                     off.Sku,
// 		Model:                   off.Model,
// 		Brand:                   off.Brand,
// 		Price:                   off.Price,
// 		Oldprice:                off.Oldprice,
// 		Availabilities:          p.Availabilities(off.Availabilities),
// 		CityPrices:              p.CityPrices(off.CityPrices),
// 		DeliveryOptions:         off.DeliveryOptions,
// 		Version:                 off.Version,
// 		DistinctVersion:         off.DistinctVersion,
// 		OfferSource:             off.OfferSource,
// 		ReindexReason:           off.ReindexReason,
// 		MatchesLastBatch:        off.MatchesLastBatch,
// 		SkipDistinctCheck:       off.SkipDistinctCheck,
// 		Priority:                off.Priority,
// 		RangePrice:              p.RangePrice(off.RangePrice),
// 		Updates:                 p.Updates(off.Updates),
// 		ArchivedAt:              off.ArchivedAt,
// 		Modification:            off.Modification,
// 		ColourCode:              off.ColourCode,
// 		Attributes:              off.Attributes,
// 		EditionDeadline:         off.EditionDeadline,
// 		VerticalCategory:        off.VerticalCategory,
// 		MasterCategory:          off.MasterCategory,
// 		Title:                   off.Title,
// 		Available:               off.Available,
// 		MinPrice:                off.MinPrice,
// 		MaxPrice:                off.MaxPrice,
// 		Points:                  off.Points,
// 		UpdatedAt:               off.UpdatedAt, // тип models.CustomTime
// 		AnyPickup:               off.AnyPickup,
// 		AnyKaspiDelivery:        off.AnyKaspiDelivery,
// 		AnyKaspiDeliveryLocal:   off.AnyKaspiDeliveryLocal,
// 		AnyKaspiDeliveryExpress: off.AnyKaspiDeliveryExpress,
// 		AnyMerchantDelivery:     off.AnyMerchantDelivery,
// 		BrandCode:               off.BrandCode,
// 		SelectedCity:            off.SelectedCity,
// 		AllCityPrices:           p.AllCityPrices(off.AllCityPrices),
// 		Stocks:                  p.Stocks(off.Stocks),
// 		OfferError:              off.OfferError,
// 	}
// }

// // Availabilities преобразует срез доменных availabilities в срез AvailabilityView.
// func (p *Presenter) Availabilities(avails []domain.Availability) []view.AvailabilityView {
// 	res := make([]view.AvailabilityView, len(avails))
// 	for i, a := range avails {
// 		res[i] = view.AvailabilityView{
// 			Available:    a.Available,
// 			StoreId:      a.StoreId,
// 			StockCount:   a.StockCount,
// 			PreOrder:     a.PreOrder,
// 			StockEnabled: a.StockEnabled,
// 		}
// 	}
// 	return res
// }

// // CityPrices преобразует срез доменных cityPrices в срез CityPriceView.
// func (p *Presenter) CityPrices(prices []domain.CityPrice) []view.CityPriceView {
// 	res := make([]view.CityPriceView, len(prices))
// 	for i, cp := range prices {
// 		res[i] = view.CityPriceView{
// 			Value:    cp.Value,
// 			CityId:   cp.CityId,
// 			Oldprice: cp.Oldprice,
// 		}
// 	}
// 	return res
// }

// // RangePrice преобразует доменную модель диапазона цен в RangePriceView.
// func (p *Presenter) RangePrice(rp domain.RangePrice) view.RangePriceView {
// 	return view.RangePriceView{
// 		MIN: rp.MIN,
// 		MAX: rp.MAX,
// 	}
// }

// // Updates преобразует срез доменных обновлений в срез UpdateView.
// func (p *Presenter) Updates(updates []domain.Update) []view.UpdateView {
// 	res := make([]view.UpdateView, len(updates))
// 	for i, u := range updates {
// 		res[i] = view.UpdateView{
// 			OperationType: u.OperationType,
// 			UploadedBy:    u.UploadedBy,
// 			UploadedDate:  u.UploadedDate,
// 		}
// 	}
// 	return res
// }

// // AllCityPrices преобразует мапу доменных allCityPrices в мапу CityPriceDetailView.
// func (p *Presenter) AllCityPrices(prices map[string]domain.CityPriceDetail) map[string]view.CityPriceDetailView {
// 	res := make(map[string]view.CityPriceDetailView, len(prices))
// 	for key, cp := range prices {
// 		res[key] = view.CityPriceDetailView{
// 			Price:    cp.Price,
// 			OldPrice: cp.OldPrice,
// 		}
// 	}
// 	return res
// }

// // Stocks преобразует срез доменных stocks в срез StockView.
// func (p *Presenter) Stocks(stocks []domain.Stock) []view.StockView {
// 	res := make([]view.StockView, len(stocks))
// 	for i, s := range stocks {
// 		res[i] = view.StockView{
// 			CityId:     s.CityId,
// 			StockLevel: p.StockLevel(s.StockLevel),
// 		}
// 	}
// 	return res
// }

// // StockLevel преобразует мапу доменных StockLevelDetail в мапу StockLevelDetailView.
// func (p *Presenter) StockLevel(level map[string]domain.StockLevelDetail) map[string]view.StockLevelDetailView {
// 	res := make(map[string]view.StockLevelDetailView, len(level))
// 	for key, d := range level {
// 		res[key] = view.StockLevelDetailView{
// 			Value:   d.Value,
// 			Enabled: d.Enabled,
// 		}
// 	}
// 	return res
// }

// // func (p *Presenter) Product(product *domain.Products) *view.ProductView {
// // 	if product == nil {
// // 		return nil
// // 	}

// // 	return &view.ProductView{
// // 		Data:  p.ProductData(product.Data),
// // 		Total: product.Total,
// // 	}
// // }

// // func (p *Presenter) ProductData(offers []domain.Offer) []view.OfferView {
// // 	if offers == nil {
// // 		return nil
// // 	}

// // 	items := make([]view.OfferView, len(offers))
// // 	for i, offer := range offers {
// // 		items[i] = *p.Offer(&offer)
// // 	}

// // 	return items
// // }

// // func (*Presenter) Offer(offer *domain.Offer) *view.OfferView {
// // 	if offer == nil {
// // 		return nil
// // 	}

// // 	return &view.OfferView{
// // 		MasterSku : offer.MasterSku,
// // 		MasterSku   : offer.MasterSku,
// // 		Images   : offer.MasterSku,
// // 		ShopLink        : offer.MasterSku,
// // 		MasterTitle       : offer.MasterSku,
// // 		OfferId                 : offer.MasterSku,                         `json:"offerId"`
// // 		FileId                  : offer.MasterSku,                         `json:"fileId"`
// // 		MerchantUid             : offer.MasterSku,                         `json:"merchantUid"`
// // 		Sku                     : offer.MasterSku,                         `json:"sku"`
// // 		Model                   : offer.MasterSku,                         `json:"model"`
// // 		Brand                   *: offer.MasterSku,                        `json:"brand"` // может быть null
// // 		Price                   float64                        `json:"price"`
// // 		Oldprice                float64                        `json:"oldprice"`
// // 		Availabilities          []AvailabilityView             `json:"availabilities"`
// // 		CityPrices              []CityPriceView                `json:"cityPrices"`
// // 		DeliveryOptions         interface{}                    `json:"deliveryOptions"` // может быть null
// // 		Version                 interface{}                    `json:"_version"`        // может быть null
// // 		DistinctVersion         int                            `json:"_distinctVersion"`
// // 		OfferSource             : offer.MasterSku,                         `json:"offerSource"`
// // 		ReindexReason           : offer.MasterSku,                         `json:"reindexReason"`
// // 		MatchesLastBatch        bool                           `json:"matchesLastBatch"`
// // 		SkipDistinctCheck       bool                           `json:"skipDistinctCheck"`
// // 		Priority                int                            `json:"priority"`
// // 		RangePrice              RangePriceView                 `json:"rangePrice"`
// // 		Updates                 []UpdateView                   `json:"updates"`
// // 		ArchivedAt              interface{}                    `json:"archivedAt"` // может быть null
// // 		Modification            bool                           `json:"modification"`
// // 		ColourCode              *: offer.MasterSku,                        `json:"colourCode"`      // может быть null
// // 		Attributes              []interface{}                  `json:"attributes"`      // массив, может быть уточнён по необходимости
// // 		EditionDeadline         interface{}                    `json:"editionDeadline"` // может быть null
// // 		VerticalCategory        : offer.MasterSku,                         `json:"verticalCategory"`
// // 		MasterCategory          : offer.MasterSku,                         `json:"masterCategory"`
// // 		Title                   : offer.MasterSku,                         `json:"title"`
// // 		Available               bool                           `json:"available"`
// // 		MinPrice                float64                        `json:"minPrice"`
// // 		MaxPrice                float64                        `json:"maxPrice"`
// // 		Points                  []: offer.MasterSku,                       `json:"points"`
// // 		UpdatedAt               models.CustomTime              `json:"updatedAt"`
// // 		AnyPickup               bool                           `json:"anyPickup"`
// // 		AnyKaspiDelivery        bool                           `json:"anyKaspiDelivery"`
// // 		AnyKaspiDeliveryLocal   bool                           `json:"anyKaspiDeliveryLocal"`
// // 		AnyKaspiDeliveryExpress bool                           `json:"anyKaspiDeliveryExpress"`
// // 		AnyMerchantDelivery     bool                           `json:"anyMerchantDelivery"`
// // 		BrandCode               *: offer.MasterSku,                        `json:"brandCode"`    // может быть null
// // 		SelectedCity            *: offer.MasterSku,                        `json:"selectedCity"` // может быть null
// // 		AllCityPrices           map[string]CityPriceDetailView `json:"allCityPrices"`
// // 		Stocks                  []StockView                    `json:"stocks"`
// // 		OfferError              interface{}                    `json:"offerError"` // может быть null

// // }

// // func (*Presenter) ProductRequest(req *domain.ProductRequest) *view.ProductRequestView {
// // 	if req == nil {
// // 		return nil
// // 	}

// // 	return &view.ProductRequestView{
// // 		ID:      req.ID,
// // 		Product: req.Product,
// // 		Status:  string(req.Status),
// // 		Remark:  req.Remark,
// // 	}
// // }

// // func (p *Presenter) ProductDataList(ctx context.Context, products []*domain.ProductData, totalCount int64) *view.ProductDataListView {
// // 	items := make([]*view.ProductDataView, len(products))
// // 	for i, lim := range products {
// // 		items[i] = p.ProductData(lim)
// // 	}

// // 	return &view.ProductDataListView{
// // 		TotalCount: totalCount,
// // 		HasNext:    pagination.HasNext(ctx, totalCount),
// // 		Items:      items,
// // 	}
// // }
