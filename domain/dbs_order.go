package domain

import "time"

// DBSOrder — сборочное задание DBS (Delivery by Seller) из WB API /api/v3/dbs/orders/new.
// Используется для уведомлений о новых заказах.
type DBSOrder struct {
	ID             int
	ShopID         int
	WBOrderID      int64   // id из WB API (уникален в рамках WB, храним для дедупликации по shop_id)
	OrderUID       string  // orderUid из WB
	NmID           int64   // nmId — номенклатура
	ChrtID         int64   // chrtId — размер/вариант
	Article        string  // артикул продавца
	SupplierStatus string  // new, confirm, deliver, receive, reject, cancel...
	WBStatus       string  // waiting, sold, canceled...
	CreatedAt      time.Time
	NotifiedAt     *time.Time // когда отправили уведомление в Telegram (nil = ещё не отправлено)
	Extra          []byte    // JSON: address, price, salePrice, options, skus и т.д.
}
