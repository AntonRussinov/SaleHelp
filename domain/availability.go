package domain

type Availability struct {
	ID           int    // внутренний идентификатор записи
	MerchantUID  string // идентификатор магазина (продавца)
	OfferID      int    // идентификатор предложения
	Available    bool   // доступность (преобразованное из "yes" / "no")
	StoreID      string // идентификатор склада или точки продаж
	PreOrder     int    // количество для предзаказа
	StockCount   int    // фактическое количество на складе
	StockEnabled bool   // флаг, указывающий, активен ли склад
}
