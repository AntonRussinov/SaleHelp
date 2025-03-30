package domain

type Stock struct {
	ID          int    // внутренний идентификатор записи
	MerchantUID string // идентификатор магазина
	OfferID     int    // идентификатор предложения
	CityID      string // идентификатор города (если применимо)
	StoreID     string // идентификатор склада (например, "30345212_PP1")
	StockCount  int    // количество товара на складе
	Available   bool   // доступность (преобразованное из "yes" / "no")
	PreOrder    int    // количество для предзаказа
	Enabled     bool   // флаг, указывающий, активен ли склад
}
