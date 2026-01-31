package mapper

import (
	"encoding/json"
	"time"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

// SaleItemToDomain конвертирует DTO записи продажи/возврата WB в домен.
// eventType: по умолчанию buyout; если forPay < 0 — return (возврат).
func SaleItemToDomain(item dto.SaleItem, shopID int) domain.SaleEvent {
	eventType := domain.SaleEventTypeBuyout
	if item.ForPay < 0 {
		eventType = domain.SaleEventTypeReturn
	}
	date, _ := item.ParseDate()
	lastChange, _ := item.ParseLastChangeDate()
	if date.IsZero() {
		date = time.Now()
	}
	if lastChange.IsZero() {
		lastChange = time.Now()
	}
	extra := saleEventExtra(item)
	return domain.SaleEvent{
		ShopID:          shopID,
		EventType:       eventType,
		WBSaleID:        item.SaleID,
		Srid:            item.Srid,
		Date:            date,
		LastChangeDate:  lastChange,
		NmID:            item.NmID,
		Barcode:         item.Barcode,
		SupplierArticle: item.SupplierArticle,
		WarehouseName:   item.WarehouseName,
		Category:        item.Category,
		Subject:         item.Subject,
		Brand:           item.Brand,
		TechSize:        item.TechSize,
		IncomeID:        item.IncomeID,
		TotalPrice:      item.TotalPrice,
		DiscountPercent: item.DiscountPercent,
		FinishedPrice:   item.FinishedPrice,
		PriceWithDisc:   item.PriceWithDisc,
		ForPay:          item.ForPay,
		CreatedAt:       time.Now(),
		NotifiedAt:      nil,
		Extra:           extra,
	}
}

func saleEventExtra(item dto.SaleItem) []byte {
	m := map[string]interface{}{
		"warehouseType":     item.WarehouseType,
		"countryName":      item.CountryName,
		"oblastOkrugName":  item.OblastOkrugName,
		"regionName":       item.RegionName,
		"spp":              item.Spp,
		"paymentSaleAmount": item.PaymentSaleAmount,
		"sticker":          item.Sticker,
		"gNumber":          item.GNumber,
	}
	b, _ := json.Marshal(m)
	return b
}
