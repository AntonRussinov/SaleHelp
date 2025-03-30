package domain

type CityPrice struct {
	ID          int     // внутренний идентификатор записи
	MerchantUID string  // идентификатор магазина
	OfferID     int     // идентификатор предложения
	CityID      string  // идентификатор города
	Price       float64 // цена в городе
	OldPrice    float64 // старая цена (если необходимо)
}
