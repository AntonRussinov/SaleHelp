package di

import (
	"salehelp/internal/storage/pggorm/repository/airecommendation"
	"salehelp/internal/storage/pggorm/repository/alert"
	"salehelp/internal/storage/pggorm/repository/alerttriggered"
	"salehelp/internal/storage/pggorm/repository/alerttriggereddeliveries"
	"salehelp/internal/storage/pggorm/repository/cronjob"
	"salehelp/internal/storage/pggorm/repository/notificationendpoint"
	"salehelp/internal/storage/pggorm/repository/product"
	"salehelp/internal/storage/pggorm/repository/productsize"
	"salehelp/internal/storage/pggorm/repository/productsizepricesnapshot"
	"salehelp/internal/storage/pggorm/repository/productsizesku"
	"salehelp/internal/storage/pggorm/repository/productsizesnapshot"
	"salehelp/internal/storage/pggorm/repository/productsnapshot"
	"salehelp/internal/storage/pggorm/repository/sellerwarehouse"
	"salehelp/internal/storage/pggorm/repository/shop"
	shopinvitetokens "salehelp/internal/storage/pggorm/repository/shop_invite_tokens"
	"salehelp/internal/storage/pggorm/repository/shopsyncleases"
	"salehelp/internal/storage/pggorm/repository/subscription"
	telegramchatstate "salehelp/internal/storage/pggorm/repository/telegram_chat_states"
	"salehelp/internal/storage/pggorm/repository/telegramaccount"
	"salehelp/internal/storage/pggorm/repository/user"
	"salehelp/internal/storage/pggorm/repository/wbwarehouse"
)

type Repositories struct {
	AIRecommendation         *airecommendation.Repository
	Alert                    *alert.Repository
	AlertTriggered           *alerttriggered.Repository
	AlertTriggeredDelivery   *alerttriggereddeliveries.Repository
	CronJob                  *cronjob.Repository
	NotificationEndpoint     *notificationendpoint.Repository
	ProductSnapshot          *productsnapshot.Repository
	ProductSizeSnapshot      *productsizesnapshot.Repository
	ProductSizePriceSnapshot *productsizepricesnapshot.Repository
	ProductSizesku           *productsizesku.Repository
	ProductSize              *productsize.Repository
	Product                  *product.Repository
	SellerWarehouse          *sellerwarehouse.Repository
	Shop                     *shop.Repository
	Shopinvitetokens         *shopinvitetokens.Repository
	ShopSyncLeases           *shopsyncleases.Repository
	Subscription             *subscription.Repository
	TelegramAccount          *telegramaccount.Repository
	TelegramChatState        *telegramchatstate.Repository
	User                     *user.Repository
	WBWarehouse              *wbwarehouse.Repository
}

func (prov *Provider) InitRepositories() {
	repository := &Repositories{
		AIRecommendation:         airecommendation.NewRepository(prov.log, prov.Database.GormDB),
		Alert:                    alert.NewRepository(prov.log, prov.Database.GormDB),
		AlertTriggered:           alerttriggered.NewRepository(prov.log, prov.Database.GormDB),
		AlertTriggeredDelivery:   alerttriggereddeliveries.NewRepository(prov.log, prov.Database.GormDB),
		CronJob:                  cronjob.NewRepository(prov.log, prov.Database.GormDB),
		NotificationEndpoint:     notificationendpoint.NewRepository(prov.log, prov.Database.GormDB),
		ProductSnapshot:          productsnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizeSnapshot:      productsizesnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizePriceSnapshot: productsizepricesnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizesku:           productsizesku.NewRepository(prov.log, prov.Database.GormDB),
		ProductSize:              productsize.NewRepository(prov.log, prov.Database.GormDB),
		Product:                  product.NewRepository(prov.log, prov.Database.GormDB),
		SellerWarehouse:          sellerwarehouse.NewRepository(prov.log, prov.Database.GormDB),
		Shop:                     shop.NewRepository(prov.log, prov.Database.GormDB),
		Shopinvitetokens:         shopinvitetokens.NewRepository(prov.log, prov.Database.GormDB),
		ShopSyncLeases:           shopsyncleases.NewRepository(prov.log, prov.Database.GormDB),
		Subscription:             subscription.NewRepository(prov.log, prov.Database.GormDB),
		TelegramAccount:          telegramaccount.NewRepository(prov.log, prov.Database.GormDB),
		TelegramChatState:        telegramchatstate.NewRepository(prov.log, prov.Database.GormDB),
		User:                     user.NewRepository(prov.log, prov.Database.GormDB),
		WBWarehouse:              wbwarehouse.NewRepository(prov.log, prov.Database.GormDB),
	}

	prov.Repositories = *repository
}
