package di

import (
	"salehelp/internal/storage/pggorm/repository/abcsnapshot"
	aireplyrepo "salehelp/internal/storage/pggorm/repository/ai_feedback_reply"
	aipackages "salehelp/internal/storage/pggorm/repository/ai_packages"
	"salehelp/internal/storage/pggorm/repository/airecommendation"
	"salehelp/internal/storage/pggorm/repository/alert"
	"salehelp/internal/storage/pggorm/repository/alerttriggered"
	"salehelp/internal/storage/pggorm/repository/alerttriggereddeliveries"
	"salehelp/internal/storage/pggorm/repository/cronjob"
	"salehelp/internal/storage/pggorm/repository/dbsorder"
	"salehelp/internal/storage/pggorm/repository/dbworder"
	"salehelp/internal/storage/pggorm/repository/fbsorder"
	linktokens "salehelp/internal/storage/pggorm/repository/link_tokens"
	"salehelp/internal/storage/pggorm/repository/notificationendpoint"
	"salehelp/internal/storage/pggorm/repository/orderrefusal"
	"salehelp/internal/storage/pggorm/repository/ordersynccursor"
	"salehelp/internal/storage/pggorm/repository/pricerestorebatch"
	"salehelp/internal/storage/pggorm/repository/product"
	"salehelp/internal/storage/pggorm/repository/productsize"
	"salehelp/internal/storage/pggorm/repository/productsizepricesnapshot"
	"salehelp/internal/storage/pggorm/repository/productsizesku"
	"salehelp/internal/storage/pggorm/repository/productsizesnapshot"
	"salehelp/internal/storage/pggorm/repository/productsnapshot"
	"salehelp/internal/storage/pggorm/repository/saleevent"
	"salehelp/internal/storage/pggorm/repository/sellerwarehouse"
	"salehelp/internal/storage/pggorm/repository/shop"
	shopairedger "salehelp/internal/storage/pggorm/repository/shop_ai_reply_ledger"
	shopaireplyquota "salehelp/internal/storage/pggorm/repository/shop_ai_reply_quota"
	aisettings "salehelp/internal/storage/pggorm/repository/ai_settings"
	shopinvitetokens "salehelp/internal/storage/pggorm/repository/shop_invite_tokens"
	"salehelp/internal/storage/pggorm/repository/shopabcsettings"
	"salehelp/internal/storage/pggorm/repository/shopnotificationsetting"
	"salehelp/internal/storage/pggorm/repository/shopsyncleases"
	"salehelp/internal/storage/pggorm/repository/subscription"
	telegramchatstate "salehelp/internal/storage/pggorm/repository/telegram_chat_states"
	"salehelp/internal/storage/pggorm/repository/telegramaccount"
	"salehelp/internal/storage/pggorm/repository/user"
	"salehelp/internal/storage/pggorm/repository/wbfeedback"
	"salehelp/internal/storage/pggorm/repository/wborder"
	"salehelp/internal/storage/pggorm/repository/wbquestion"
	"salehelp/internal/storage/pggorm/repository/wbwarehouse"
)

type Repositories struct {
	ABCSnapshot              *abcsnapshot.Repository
	AIRecommendation         *airecommendation.Repository
	AIPackages               *aipackages.Repository
	Alert                    *alert.Repository
	AlertTriggered           *alerttriggered.Repository
	AlertTriggeredDelivery   *alerttriggereddeliveries.Repository
	CronJob                  *cronjob.Repository
	DBWOrder                 *dbworder.Repository
	DBSOrder                 *dbsorder.Repository
	FBSOrder                 *fbsorder.Repository
	SaleEvent                *saleevent.Repository
	NotificationEndpoint     *notificationendpoint.Repository
	OrderRefusal             *orderrefusal.Repository
	OrderSyncCursor          *ordersynccursor.Repository
	PriceRestoreBatch        *pricerestorebatch.Repository
	WBOrder                  *wborder.Repository
	WBFeedback               *wbfeedback.Repository
	WBQuestion               *wbquestion.Repository
	ProductSnapshot          *productsnapshot.Repository
	ProductSizeSnapshot      *productsizesnapshot.Repository
	ProductSizePriceSnapshot *productsizepricesnapshot.Repository
	ProductSizesku           *productsizesku.Repository
	ProductSize              *productsize.Repository
	Product                  *product.Repository
	SellerWarehouse          *sellerwarehouse.Repository
	Shop                     *shop.Repository
	ShopABCSettings          *shopabcsettings.Repository
	AISettings               *aisettings.Repository
	ShopAIReplyLedger        *shopairedger.Repository
	ShopAIReplyQuota         *shopaireplyquota.Repository
	AIFeedbackReply          *aireplyrepo.Repository
	LinkTokens               *linktokens.Repository
	ShopInviteTokens         *shopinvitetokens.Repository
	ShopNotificationSetting  *shopnotificationsetting.Repository
	ShopSyncLeases           *shopsyncleases.Repository
	Subscription             *subscription.Repository
	TelegramAccount          *telegramaccount.Repository
	TelegramChatState        *telegramchatstate.Repository
	User                     *user.Repository
	WBWarehouse              *wbwarehouse.Repository
}

func (prov *Provider) InitRepositories() {
	repository := &Repositories{
		ABCSnapshot:              abcsnapshot.NewRepository(prov.log, prov.Database.GormDB),
		AIRecommendation:         airecommendation.NewRepository(prov.log, prov.Database.GormDB),
		AIPackages:               aipackages.NewRepository(prov.log, prov.Database.GormDB),
		Alert:                    alert.NewRepository(prov.log, prov.Database.GormDB),
		AlertTriggered:           alerttriggered.NewRepository(prov.log, prov.Database.GormDB),
		AlertTriggeredDelivery:   alerttriggereddeliveries.NewRepository(prov.log, prov.Database.GormDB),
		CronJob:                  cronjob.NewRepository(prov.log, prov.Database.GormDB),
		DBWOrder:                 dbworder.NewRepository(prov.log, prov.Database.GormDB),
		DBSOrder:                 dbsorder.NewRepository(prov.log, prov.Database.GormDB),
		FBSOrder:                 fbsorder.NewRepository(prov.log, prov.Database.GormDB),
		SaleEvent:                saleevent.NewRepository(prov.log, prov.Database.GormDB),
		NotificationEndpoint:     notificationendpoint.NewRepository(prov.log, prov.Database.GormDB),
		OrderRefusal:             orderrefusal.NewRepository(prov.log, prov.Database.GormDB),
		OrderSyncCursor:          ordersynccursor.NewRepository(prov.log, prov.Database.GormDB),
		PriceRestoreBatch:        pricerestorebatch.NewRepository(prov.log, prov.Database.GormDB),
		WBOrder:                  wborder.NewRepository(prov.log, prov.Database.GormDB),
		WBFeedback:               wbfeedback.NewRepository(prov.log, prov.Database.GormDB),
		WBQuestion:               wbquestion.NewRepository(prov.log, prov.Database.GormDB),
		ProductSnapshot:          productsnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizeSnapshot:      productsizesnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizePriceSnapshot: productsizepricesnapshot.NewRepository(prov.log, prov.Database.GormDB),
		ProductSizesku:           productsizesku.NewRepository(prov.log, prov.Database.GormDB),
		ProductSize:              productsize.NewRepository(prov.log, prov.Database.GormDB),
		Product:                  product.NewRepository(prov.log, prov.Database.GormDB),
		SellerWarehouse:          sellerwarehouse.NewRepository(prov.log, prov.Database.GormDB),
		Shop:                     shop.NewRepository(prov.log, prov.Database.GormDB),
		ShopABCSettings:          shopabcsettings.NewRepository(prov.log, prov.Database.GormDB),
		AISettings:               aisettings.NewRepository(prov.log, prov.Database.GormDB),
		ShopAIReplyLedger:        shopairedger.NewRepository(prov.log, prov.Database.GormDB),
		ShopAIReplyQuota:         shopaireplyquota.NewRepository(prov.log, prov.Database.GormDB),
		AIFeedbackReply:          aireplyrepo.NewRepository(prov.log, prov.Database.GormDB),
		LinkTokens:               linktokens.NewRepository(prov.log, prov.Database.GormDB),
		ShopInviteTokens:         shopinvitetokens.NewRepository(prov.log, prov.Database.GormDB),
		ShopNotificationSetting:  shopnotificationsetting.NewRepository(prov.log, prov.Database.GormDB),
		ShopSyncLeases:           shopsyncleases.NewRepository(prov.log, prov.Database.GormDB),
		Subscription:             subscription.NewRepository(prov.log, prov.Database.GormDB),
		TelegramAccount:          telegramaccount.NewRepository(prov.log, prov.Database.GormDB),
		TelegramChatState:        telegramchatstate.NewRepository(prov.log, prov.Database.GormDB),
		User:                     user.NewRepository(prov.log, prov.Database.GormDB),
		WBWarehouse:              wbwarehouse.NewRepository(prov.log, prov.Database.GormDB),
	}

	prov.Repositories = *repository
}
