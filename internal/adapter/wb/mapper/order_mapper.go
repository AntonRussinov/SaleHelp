package mapper

import (
	"time"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

// OrderItemToRefusal конвертирует DTO заказа WB в домен OrderRefusal.
// Вызывать только для записей с isCancel == true (отказ / не забран).
// reportLoc — пояс магазина для парсинга дат WB; nil = глобальный из конфига.
func OrderItemToRefusal(item dto.OrderItem, shopID int, reportLoc *time.Location) domain.OrderRefusal {
	date, _ := item.ParseDate(reportLoc)
	lastChange, _ := item.ParseLastChangeDate(reportLoc)
	cancelDate, _ := item.ParseCancelDate(reportLoc)
	if date.IsZero() {
		date = time.Now()
	}
	if lastChange.IsZero() {
		lastChange = time.Now()
	}
	if cancelDate.IsZero() {
		cancelDate = lastChange
	}
	return domain.OrderRefusal{
		ShopID:          shopID,
		Srid:            item.Srid,
		NmID:            item.NmID,
		Barcode:         item.Barcode,
		SupplierArticle: item.SupplierArticle,
		Date:            date,
		CancelDate:      cancelDate,
		LastChangeDate:  lastChange,
		CreatedAt:       time.Now(),
		NotifiedAt:      nil,
	}
}
