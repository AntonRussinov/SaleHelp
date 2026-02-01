package mapper

import (
	"encoding/json"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

// DBWOrderNewToDomain конвертирует DTO нового сборочного задания DBW в домен.
// supplier_status и wb_status по умолчанию "new" и "waiting" (для новых заказов из /dbw/orders/new).
func DBWOrderNewToDomain(o dto.DBWOrderNew, shopID int) domain.DBWOrder {
	extra := dbwOrderExtra(o)
	return domain.DBWOrder{
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

// dbwOrderExtra сериализует дополнительные поля в JSONB (address, price, skus и т.д.).
func dbwOrderExtra(o dto.DBWOrderNew) []byte {
	m := map[string]interface{}{
		"warehouseId":       o.WarehouseID,
		"price":             o.Price,
		"convertedPrice":    o.ConvertedPrice,
		"salePrice":         o.SalePrice,
		"currencyCode":      o.CurrencyCode,
		"cargoType":         o.CargoType,
		"comment":           o.Comment,
		"requiredMeta":      o.RequiredMeta,
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
