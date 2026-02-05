package mapper

import (
	"time"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

// OrderItemToRefusal конвертирует DTO заказа WB в домен OrderRefusal.
// Вызывать только для записей с isCancel == true (отказ / не забран).
func OrderItemToRefusal(item dto.OrderItem, shopID int) domain.OrderRefusal {
	date, _ := item.ParseDate()
	lastChange, _ := item.ParseLastChangeDate()
	cancelDate, _ := item.ParseCancelDate()
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
