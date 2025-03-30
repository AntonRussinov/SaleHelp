package domain

import "time"

type Offer struct {
	ID                      int        // внутренний идентификатор предложения
	OfferExternalID         string     // внешний идентификатор (offerId из JSON)
	SKU                     string     // SKU, используется для идентификации предложения
	ProductID               int        // ссылка на товар (из Product)
	MerchantUID             string     // ссылка на продавца (из Merchant)
	ShopLink                string     // ссылка на страницу товара в магазине
	FileID                  string     // идентификатор файла (file_id)
	OfferSource             string     // источник предложения
	ReindexReason           string     // причина переиндексации
	MatchesLastBatch        bool       // флаг, соответствует ли последней партии
	SkipDistinctCheck       bool       // флаг для проверки уникальности
	Priority                int        // приоритет предложения
	Price                   float64    // цена
	OldPrice                float64    // старая цена
	Available               bool       // доступность товара
	MinPrice                float64    // минимальная цена (если актуально)
	MaxPrice                float64    // максимальная цена (если актуально)
	ArchivedAt              *time.Time // дата архивирования, если применимо
	Modification            bool       // флаг модификации
	ColourCode              string     // код цвета, если есть
	EditionDeadline         *time.Time // срок действия версии, если применяется
	AnyPickup               bool       // доступна ли самовывоз
	AnyKaspiDelivery        bool       // доступна ли доставка Kaspi
	AnyKaspiDeliveryLocal   bool       // локальная доставка Kaspi
	AnyKaspiDeliveryExpress bool       // экспресс-доставка Kaspi
	AnyMerchantDelivery     bool       // доставка от продавца
	BrandCode               string     // код бренда, если используется
	SelectedCity            string     // выбранный город
	OfferError              string     // информация об ошибках, если имеются
	// Вложенные связанные модели (агрегаты)
	Product        *Product       // Связанный товар
	Merchant       *Merchant      // Связанный продавец
	Images         []Image        // Список изображений предложения
	Availabilities []Availability // Наличие товара (по точкам/складам)
	CityPrices     []CityPrice    // Цены для разных городов
	Stocks         []Stock        // Остатки на складах
}
