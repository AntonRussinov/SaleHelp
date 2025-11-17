package domain

type Merchant struct {
	MerchantUID string // уникальный идентификатор продавца (merchantUid)
	CompanyName string // название компании
	IsEnabled   bool   // статус активности магазина
	Login       string // логин для доступа к API
	Password    string // пароль для доступа к API
}
