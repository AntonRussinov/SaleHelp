package mapper

import (
	"encoding/json"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

// DBSOrderNewToDomain конвертирует DTO нового сборочного задания DBS в домен.
// supplier_status и wb_status по умолчанию "new" и "waiting" (для новых заказов из /dbs/orders/new).
func DBSOrderNewToDomain(o dto.DBSOrderNew, shopID int) domain.DBSOrder {
	extra := dbsOrderExtra(o)
	return domain.DBSOrder{
		ShopID:         shopID,
		WBOrderID:      o.ID,
		OrderUID:       o.OrderUID,
		NmID:           o.NmID,
		ChrtID:         o.ChrtID,
		Article:        o.Article,
		SupplierStatus: "new",
		WBStatus:       "waiting",
		CreatedAt:      o.CreatedAt,
		NotifiedAt:     nil,
		Extra:          extra,
	}
}

// dbsOrderExtra сериализует дополнительные поля в JSONB (address, price, skus и т.д.).
func dbsOrderExtra(o dto.DBSOrderNew) []byte {
	m := map[string]interface{}{
		"warehouseId":       o.WarehouseID,
		"price":             o.Price,
		"finalPrice":        o.FinalPrice,
		"convertedPrice":    o.ConvertedPrice,
		"salePrice":         o.SalePrice,
		"currencyCode":      o.CurrencyCode,
		"cargoType":         o.CargoType,
		"comment":           o.Comment,
		"requiredMeta":      o.RequiredMeta,
		"deliveryType":      o.DeliveryType,
		"skus":              o.Skus,
	}
	if o.Address != nil {
		m["address"] = o.Address
	}
	if o.Options != nil {
		m["options"] = o.Options
	}
	b, _ := json.Marshal(m)
	return b
}
