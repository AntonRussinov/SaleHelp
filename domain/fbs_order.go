package domain

import "time"

// FBSOrder — сборочное задание FBS (Fulfillment by Seller) из WB API.
// Используется для уведомлений о новых заказах.
type FBSOrder struct {
	ID             int
	ShopID         int
	WBOrderID      int64   // id из WB API (уникален в рамках WB, храним для дедупликации по shop_id)
	OrderUID       string  // orderUid из WB
	NmID           int64   // nmId — номенклатура
	ChrtID         int64   // chrtId — размер/вариант
	Article        string  // артикул продавца
	SupplierStatus string  // new, confirm, complete, cancel
	WBStatus       string  // waiting, sorted, sold, canceled, ...
	CreatedAt      time.Time // время создания заказа в WB
	NotifiedAt     *time.Time // когда отправили уведомление в Telegram (nil = ещё не отправлено)
	Extra          []byte  // JSON: address, price, salePrice, offices, skus и т.д.
}
