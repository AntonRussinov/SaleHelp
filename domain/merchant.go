package domain

type Merchant struct {
	MerchantUID string // уникальный идентификатор продавца (merchantUid)
	CompanyName string // название компании
	IsEnabled   bool   // статус активности магазина
}
