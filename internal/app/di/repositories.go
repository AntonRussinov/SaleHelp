package di

import (
	"salehelp/internal/storage/pggorm/repository/airecommendation"
	"salehelp/internal/storage/pggorm/repository/alert"
	"salehelp/internal/storage/pggorm/repository/alerttriggered"
	"salehelp/internal/storage/pggorm/repository/alerttriggereddeliveries"
	"salehelp/internal/storage/pggorm/repository/cronjob"
	"salehelp/internal/storage/pggorm/repository/dbsorder"
	"salehelp/internal/storage/pggorm/repository/dbworder"
	"salehelp/internal/storage/pggorm/repository/fbsorder"
	"salehelp/internal/storage/pggorm/repository/notificationendpoint"
	"salehelp/internal/storage/pggorm/repository/orderrefusal"
	"salehelp/internal/storage/pggorm/repository/ordersynccursor"
	"salehelp/internal/storage/pggorm/repository/product"
	"salehelp/internal/storage/pggorm/repository/productsize"
	"salehelp/internal/storage/pggorm/repository/productsizepricesnapshot"
	"salehelp/internal/storage/pggorm/repository/productsizesku"
	"salehelp/internal/storage/pggorm/repository/productsizesnapshot"
	"salehelp/internal/storage/pggorm/repository/productsnapshot"
	"salehelp/internal/storage/pggorm/repository/saleevent"
	"salehelp/internal/storage/pggorm/repository/sellerwarehouse"
	"salehelp/internal/storage/pggorm/repository/shop"
	shopinvitetokens "salehelp/internal/storage/pggorm/repository/shop_invite_tokens"
	"salehelp/internal/storage/pggorm/repository/shopsyncleases"
	"salehelp/internal/storage/pggorm/repository/subscription"
	subscriptionaddon "salehelp/internal/storage/pggorm/repository/subscription_addon"
	subscriptionnotificationevent "salehelp/internal/storage/pggorm/repository/subscription_notification_event"
	telegramchatstate "salehelp/internal/storage/pggorm/repository/telegram_chat_states"
	"salehelp/internal/storage/pggorm/repository/telegramaccount"
	"salehelp/internal/storage/pggorm/repository/user"
	"salehelp/internal/storage/pggorm/repository/wbfeedback"
	"salehelp/internal/storage/pggorm/repository/wbquestion"
	"salehelp/internal/storage/pggorm/repository/wborder"
	"salehelp/internal/storage/pggorm/repository/wbwarehouse"
	"salehelp/internal/storage/pggorm/txrunner"
)

type Repositories struct {
	AIRecommendation              *airecommendation.Repository
	Alert                         *alert.Repository
	AlertTriggered                *alerttriggered.Repository
	AlertTriggeredDelivery        *alerttriggereddeliveries.Repository
	CronJob                       *cronjob.Repository
	DBWOrder                      *dbworder.Repository
	DBSOrder                      *dbsorder.Repository
	FBSOrder                      *fbsorder.Repository
	SaleEvent                     *saleevent.Repository
	NotificationEndpoint          *notificationendpoint.Repository
	OrderRefusal                  *orderrefusal.Repository
	OrderSyncCursor               *ordersynccursor.Repository
	WBOrder                       *wborder.Repository
	WBFeedback                    *wbfeedback.Repository
	WBQuestion                    *wbquestion.Repository
	ProductSnapshot               *productsnapshot.Repository
	ProductSizeSnapshot           *productsizesnapshot.Repository
	ProductSizePriceSnapshot      *productsizepricesnapshot.Repository
	ProductSizesku                *productsizesku.Repository
	ProductSize                   *productsize.Repository
	Product                       *product.Repository
	SellerWarehouse               *sellerwarehouse.Repository
	Shop                          *shop.Repository
	ShopInviteTokens              *shopinvitetokens.Repository
	ShopSyncLeases                *shopsyncleases.Repository
	Subscription                  *subscription.Repository
	SubscriptionAddon             *subscriptionaddon.Repository
	SubscriptionNotificationEvent *subscriptionnotificationevent.Repository
	TelegramAccount               *telegramaccount.Repository
	TelegramChatState             *telegramchatstate.Repository
	User                          *user.Repository
	WBWarehouse                   *wbwarehouse.Repository
	TxRunner                      *txrunner.Runner
}

func (prov *Provider) InitRepositories() {

	subscrRepo := subscription.NewRepository(prov.log, prov.Database.GormDB)
	subscrAddonRepo := subscriptionaddon.NewRepository(prov.log, prov.Database.GormDB)

	repository := &Repositories{
		AIRecommendation:              airecommendation.NewRepository(prov.log, prov.Database.GormDB),
		Alert:                         alert.NewRepository(prov.log, prov.Database.GormDB),
		AlertTriggered:                alerttriggered.NewRepository(prov.log, prov.Database.GormDB),
		AlertTriggeredDelivery:        alerttriggereddeliveries.NewRepository(prov.log, prov.Database.GormDB),
		CronJob:                       cronjob.NewRepository(prov.log, prov.Database.GormDB),
		DBWOrder:                      dbworder.NewRepository(prov.log, prov.Database.GormDB),
		DBSOrder:                      dbsorder.NewRepository(prov.log, prov.Database.GormDB),
		FBSOrder:                      fbsorder.NewRepository(prov.log, prov.Database.GormDB),
		SaleEvent:                     saleevent.NewRepository(prov.log, prov.Database.GormDB),
		NotificationEndpoint:          notificationendpoint.NewRepository(prov.log, prov.Database.GormDB),
		OrderRefusal:                  orderrefusal.NewRepository(prov.log, prov.Database.GormDB),
		OrderSyncCursor:               ordersynccursor.NewRepository(prov.log, prov.Database.GormDB),
		WBOrder:                       wborder.NewRepository(prov.log, prov.Database.GormDB),
		WBFeedback:                    wbfeedback.NewRepository(prov.log, prov.Database.GormDB),
		WBQuestion:                    wbquestion.NewRepository(prov.log, prov.Database.GormDB),
		ProductSnapshot:               productsnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizeSnapshot:           productsizesnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizePriceSnapshot:      productsizepricesnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizesku:                productsizesku.NewRepository(prov.log, prov.Database.GormDB),
		ProductSize:                   productsize.NewRepository(prov.log, prov.Database.GormDB),
		Product:                       product.NewRepository(prov.log, prov.Database.GormDB),
		SellerWarehouse:               sellerwarehouse.NewRepository(prov.log, prov.Database.GormDB),
		Shop:                          shop.NewRepository(prov.log, prov.Database.GormDB),
		ShopInviteTokens:              shopinvitetokens.NewRepository(prov.log, prov.Database.GormDB),
		ShopSyncLeases:                shopsyncleases.NewRepository(prov.log, prov.Database.GormDB),
		Subscription:                  subscrRepo,
		SubscriptionAddon:             subscrAddonRepo,
		SubscriptionNotificationEvent: subscriptionnotificationevent.NewRepository(prov.log, prov.Database.GormDB),
		TelegramAccount:               telegramaccount.NewRepository(prov.log, prov.Database.GormDB),
		TelegramChatState:             telegramchatstate.NewRepository(prov.log, prov.Database.GormDB),
		User:                          user.NewRepository(prov.log, prov.Database.GormDB),
		WBWarehouse:                   wbwarehouse.NewRepository(prov.log, prov.Database.GormDB),
		TxRunner:                      txrunner.New(prov.Database.GormDB, subscrRepo, subscrAddonRepo),
	}

	prov.Repositories = *repository
}
