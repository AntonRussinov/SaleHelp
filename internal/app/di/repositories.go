package di

import (
	"salehelp/internal/storage/pggorm/repository/airecommendation"
	"salehelp/internal/storage/pggorm/repository/alert"
	"salehelp/internal/storage/pggorm/repository/alerttriggered"
	"salehelp/internal/storage/pggorm/repository/cronjob"
	"salehelp/internal/storage/pggorm/repository/product"
	"salehelp/internal/storage/pggorm/repository/productsnapshot"
	"salehelp/internal/storage/pggorm/repository/shop"
	"salehelp/internal/storage/pggorm/repository/telegramaccount"
	"salehelp/internal/storage/pggorm/repository/user"
)

type Repositories struct {
	AIRecommendation *airecommendation.Repository
	Alert            *alert.Repository
	AlertTriggered   *alerttriggered.Repository
	CronJob          *cronjob.Repository
	ProductSnapshot  *productsnapshot.Repository
	Product          *product.Repository
	Shop             *shop.Repository
	Subscription     *shop.Repository
	TelegrammAccount *telegramaccount.Repository
	User             *user.Repository
}

func (prov *Provider) InitRepositories() *Repositories {
	return &Repositories{
		AIRecommendation: airecommendation.NewRepository(prov.log, prov.Database.GormDB),
		Alert:            alert.NewRepository(prov.log, prov.Database.GormDB),
		AlertTriggered:   alerttriggered.NewRepository(prov.log, prov.Database.GormDB),
		CronJob:          cronjob.NewRepository(prov.log, prov.Database.GormDB),
		ProductSnapshot:  productsnapshot.NewRepository(prov.log, prov.Database.GormDB),
		Product:          product.NewRepository(prov.log, prov.Database.GormDB),
		Shop:             shop.NewRepository(prov.log, prov.Database.GormDB),
		Subscription:     shop.NewRepository(prov.log, prov.Database.GormDB),
		TelegrammAccount: telegramaccount.NewRepository(prov.log,
			prov.Database.GormDB),
		User: user.NewRepository(prov.log, prov.Database.GormDB),
	}
}
