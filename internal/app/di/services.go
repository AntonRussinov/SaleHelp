package di

import (
	"log/slog"
	"os"
	"salehelp/internal/services/alert"
	"salehelp/internal/services/monitoring"
	notifications "salehelp/internal/services/notification"
	"salehelp/internal/services/product"
	"salehelp/internal/services/subscription"
	subscriptionmonitor "salehelp/internal/services/subscription_monitor"
	"salehelp/internal/services/telegrambot"
	"salehelp/pkg/notify"
	"salehelp/pkg/notify/telegram"
	"time"
)

type Services struct {
	//User *user.Service
	//Limit         *limit.Service
	AlertEngine         *alert.AlertEngine
	Monitoring          *monitoring.Service
	Notifier            *notifications.Service //Dispatcher
	AlertWorker         *alert.Worker
	Product             *product.Service
	TelegramBot         *telegrambot.Service
	Subscription        *subscription.Service
	SubscriptionMonitor *subscriptionmonitor.Monitor
}

func (prov *Provider) InitServices() {
	prov.initIndependentServices()
	prov.initDependentServices()
}

func (prov *Provider) initIndependentServices() {

	prov.Services.Subscription = subscription.New(
		prov.Repositories.TxRunner,
		prov.Repositories.Subscription,
		prov.Repositories.SubscriptionAddon,
		prov.Repositories.NotificationEndpoint,
	)

}

func (prov *Provider) initDependentServices() {

	// 1) senders registry
	tgSender, err := telegram.NewSender(prov.cfg.TelegramToken)
	if err != nil {
		panic(err) // или вернуть ошибку из InitServices
	}
	senders := []notify.Sender{
		tgSender,
	}

	tgAdapter := telegrambot.NewTelegramAdapter(tgSender)
	registry := notifications.NewSenderRegistry(senders)

	prov.Services.Monitoring = monitoring.NewMonitoringService(
		prov.Repositories.Shop,
		prov.Repositories.Product,
		prov.Repositories.ProductSnapshot,
		prov.Repositories.ProductSizeSnapshot,
		prov.Repositories.ProductSizePriceSnapshot,
		prov.Repositories.ProductSize,
		prov.Repositories.ProductSizesku,
		prov.Repositories.WBWarehouse,
		prov.Repositories.SellerWarehouse,
		prov.Repositories.User,
		tgAdapter,
		prov.WBClient,
		prov.log,
	)

	dispatcher := notifications.NewDispatcher(
		prov.Repositories.AlertTriggeredDelivery, // deliveries repo
		prov.Repositories.AlertTriggered,         // triggered repo
		prov.Repositories.Product,                // products repo
		prov.Repositories.ProductSize,            // product sizes repo
		prov.Repositories.NotificationEndpoint,   // endpoints repo
		registry,
		getWorkerID(),
	)

	prov.Services.Notifier = notifications.NewService(dispatcher, 5*time.Second)
	// prov.Services.Product = product.NewService(
	// 	prov.log, prov.Repositories.Stock, prov.Repositories.Product,
	// 	prov.Repositories.CityPrice, prov.Repositories.Offer,
	// 	prov.Repositories.Merchant, prov.Repositories.Image,
	// 	prov.KaspiPortalAPI)

	prov.Services.AlertEngine = alert.NewAlertEngine(
		prov.Services.Subscription,
		prov.Repositories.Alert,
		prov.Repositories.AlertTriggered,
		prov.Repositories.Product,
		prov.Repositories.ProductSize,
		prov.Repositories.ProductSizeSnapshot,
		prov.Repositories.ProductSnapshot,
		prov.Repositories.AlertTriggeredDelivery,
		prov.log,
	)

	prov.Services.AlertWorker = alert.NewWorker(
		prov.log,
		prov.Repositories.Shop,
		prov.Services.AlertEngine,
		1*time.Minute,
	)

	prov.Services.TelegramBot = telegrambot.NewService(
		prov.log,
		tgAdapter,
		prov.Repositories.TelegramChatState,
		prov.Repositories.ShopInviteTokens,
		prov.Repositories.User,
		prov.Repositories.NotificationEndpoint,
		prov.Repositories.Shop,
		prov.Services.Monitoring,
		prov.cfg.TelegramBotUsername,
		prov.Services.Subscription,
		prov.Repositories.Subscription,
		prov.cfg.TelegramAdminID,
	)

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)
	prov.log.Info("subscription monitor time check",
		slog.String("location", loc.String()),
		slog.String("now", now.Format(time.RFC3339)),
	)

	prov.Services.SubscriptionMonitor = subscriptionmonitor.New(
		prov.Repositories.Subscription,
		prov.Repositories.SubscriptionNotificationEvent,
		prov.Repositories.NotificationEndpoint,
		tgAdapter,
		loc,
	)

}

// Настроить поулчение корректного worker ID
// когда будет запускаться несколько воркеров
// (например, в k8s/azure)
func getWorkerID() string {
	// 1) из env (лучше всего в k8s/azure)
	if v := os.Getenv("WORKER_ID"); v != "" {
		return v
	}
	// 2) hostname
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	// 3) fallback
	return "worker-unknown"
}
