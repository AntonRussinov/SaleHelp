package dto

type GetStocksRequest struct {
	// chrtIds — основной способ, массив ID размеров (size IDs).
	ChrtIDs []int64 `json:"chrtIds"`

	// skus — deprecated у WB (будет отключено 9 февраля 2026 года), оставляем только для обратной совместимости.
	Skus []string `json:"skus,omitempty"`
}

type GetStocksResponse struct {
	Stocks []StockItem `json:"stocks"`
}

type StockItem struct {
	ChrtID int64  `json:"chrtId"`
	Sku    string `json:"sku"`
	Amount int    `json:"amount"`
}
