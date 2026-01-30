package mapper

import (
	"encoding/json"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

// FBSOrderNewToDomain конвертирует DTO нового сборочного задания в домен.
// supplier_status и wb_status по умолчанию "new" и "waiting" (для новых заказов из /orders/new).
func FBSOrderNewToDomain(o dto.FBSOrderNew, shopID int) domain.FBSOrder {
	extra := fbsOrderExtra(o)
	return domain.FBSOrder{
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

// fbsOrderExtra сериализует дополнительные поля в JSONB (address, price, offices и т.д.).
func fbsOrderExtra(o dto.FBSOrderNew) []byte {
	m := map[string]interface{}{
		"deliveryType": o.DeliveryType,
		"price":        o.Price,
		"finalPrice":   o.FinalPrice,
		"salePrice":    o.SalePrice,
		"warehouseId":  o.WarehouseID,
		"officeId":     o.OfficeID,
		"offices":      o.Offices,
		"skus":         o.Skus,
		"comment":      o.Comment,
	}
	if o.Address != nil {
		m["address"] = o.Address
	}
	b, _ := json.Marshal(m)
	return b
}
