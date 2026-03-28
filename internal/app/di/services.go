package di

import (
	"context"
	"log/slog"
	"os"
	"salehelp/internal/services/abc_analysis"
	"salehelp/internal/services/ai_reply"
	"salehelp/internal/services/alert"
	auth_google "salehelp/internal/services/auth_google"
	"salehelp/internal/services/monitoring"
	notifications "salehelp/internal/services/notification"
	"salehelp/internal/services/notification_settings"
	"salehelp/internal/services/price_restore"
	"salehelp/internal/services/product"
	"salehelp/internal/services/stock_replenish"
	"salehelp/internal/services/subscription"
	subscriptionmonitor "salehelp/internal/services/subscription_monitor"
	"salehelp/internal/services/telegrambot"
	"salehelp/pkg/ai"
	"salehelp/pkg/ai/client"
	"salehelp/pkg/notify"
	"salehelp/pkg/notify/telegram"
	"time"
)

type Services struct {
	ABCAnalysis *abc_analysis.Service
	//User *user.Service
	//Limit         *limit.Service
	AlertService         *alert.AlertService
	AlertEngine          *alert.AlertEngine
	Monitoring           *monitoring.Service
	Notifier             *notifications.Service //Dispatcher
	AlertWorker          *alert.Worker
	Product              *product.Service
	ProductList          *product.ListService
	TelegramBot          *telegrambot.Service
	Subscription         *subscription.Service
	SubscriptionMonitor  *subscriptionmonitor.Monitor
	StockReplenish       *stock_replenish.Service
	PriceRestore         *price_restore.Service
	AIReply              *ai_reply.Service
	AuthGoogle           *auth_google.Service
	NotificationSettings *notification_settings.Service
}

func (prov *Provider) InitServices() {
	prov.initIndependentServices()
	prov.initDependentServices()
}

func (prov *Provider) initIndependentServices() {
	prov.Services.AlertService = alert.NewAlertService(
		prov.Repositories.Alert,
		prov.log,
	)

	prov.Services.Subscription = subscription.New(
		prov.Repositories.TxRunner,
		prov.Repositories.Subscription,
		prov.Repositories.SubscriptionAddon,
		prov.Repositories.NotificationEndpoint,
	)

	prov.Services.NotificationSettings = notification_settings.New(prov.Repositories.ShopNotificationSetting)

}

func (prov *Provider) initDependentServices() {

	prov.Services.ProductList = product.NewListService(prov.Repositories.Product)

	// 1) senders registry
	tgSender, err := telegram.NewSender(prov.cfg.TelegramToken)
	if err != nil {
		panic(err) // или вернуть ошибку из InitServices
	}
	if err := tgSender.SetMyCommands(context.Background(), telegram.DefaultRussianCommands()); err != nil {
		prov.log.Warn("telegram setMyCommands failed", "err", err)
	}
	senders := []notify.Sender{
		tgSender,
	}

	tgAdapter := telegrambot.NewTelegramAdapter(tgSender)
	prov.telegramAdapter = tgAdapter
	registry := notifications.NewSenderRegistry(senders)

	// Пояс по умолчанию для отчётов WB (sales, orders); при отсутствии у магазина wb_report_timezone.
	wbReportTZ := prov.cfg.WBReportTimezone
	if wbReportTZ == "" {
		wbReportTZ = "Europe/Moscow"
	}
	wbReportLoc, errLoc := time.LoadLocation(wbReportTZ)
	if errLoc != nil {
		prov.log.Warn("wb default report timezone load failed, using Moscow", "tz", wbReportTZ, "err", errLoc)
		wbReportLoc, _ = time.LoadLocation("Europe/Moscow")
		if wbReportLoc == nil {
			wbReportLoc = time.FixedZone("MSK", 3*3600)
		}
	}

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
		prov.Repositories.FBSOrder,
		prov.Repositories.DBWOrder,
		prov.Repositories.DBSOrder,
		prov.Repositories.SaleEvent,
		prov.Repositories.OrderRefusal,
		prov.Repositories.WBOrder,
		prov.Repositories.OrderSyncCursor,
		prov.Repositories.WBFeedback,
		prov.Repositories.WBQuestion,
		prov.Repositories.NotificationEndpoint,
		prov.Repositories.User,
		tgAdapter,
		prov.WBClient,
		wbReportLoc,
		prov.log,
	)

	dispatcher := notifications.NewDispatcher(
		prov.Repositories.AlertTriggeredDelivery,
		prov.Repositories.AlertTriggered,
		prov.Repositories.Product,
		prov.Repositories.ProductSize,
		prov.Repositories.NotificationEndpoint,
		prov.Repositories.PriceRestoreBatch,
		prov.Services.NotificationSettings,
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
		prov.Services.NotificationSettings,
		prov.Repositories.Alert,
		prov.Repositories.AlertTriggered,
		prov.Repositories.Product,
		prov.Repositories.ProductSize,
		prov.Repositories.ProductSizeSnapshot,
		prov.Repositories.ProductSnapshot,
		prov.Repositories.AlertTriggeredDelivery,
		prov.Repositories.FBSOrder,
		prov.Repositories.DBWOrder,
		prov.Repositories.DBSOrder,
		prov.Repositories.SaleEvent,
		prov.Repositories.OrderRefusal,
		prov.Repositories.WBFeedback,
		prov.Repositories.WBQuestion,
		prov.Repositories.NotificationEndpoint,
		tgAdapter,
		prov.Repositories.ABCSnapshot,
		prov.log,
	)

	prov.Services.AlertWorker = alert.NewWorker(
		prov.log,
		prov.Repositories.Shop,
		prov.Services.AlertEngine,
		1*time.Minute,
	)

	prov.Services.ABCAnalysis = abc_analysis.New(
		prov.Repositories.SaleEvent,
		prov.Repositories.Product,
		prov.Repositories.ABCSnapshot,
		prov.log,
	)

	prov.Services.PriceRestore = price_restore.NewService(
		prov.log,
		prov.Repositories.PriceRestoreBatch,
		prov.Repositories.AlertTriggered,
		prov.Repositories.ProductSnapshot,
		prov.Repositories.ProductSizePriceSnapshot,
		prov.Repositories.Product,
		prov.Repositories.ProductSize,
		prov.Repositories.Shop,
		prov.WBClient,
	)

	prov.Services.TelegramBot = telegrambot.NewService(
		prov.log,
		tgAdapter,
		prov.Repositories.TelegramChatState,
		prov.Repositories.ShopInviteTokens,
		prov.Repositories.User,
		prov.Repositories.LinkTokens,
		prov.Repositories.NotificationEndpoint,
		prov.Repositories.Shop,
		prov.Services.Monitoring,
		prov.cfg.TelegramBotUsername,
		prov.Services.Subscription,
		prov.Repositories.Subscription,
		prov.Services.ABCAnalysis,
		prov.Repositories.Product,
		prov.Services.PriceRestore,
		prov.Repositories.PriceRestoreBatch,
		prov.Services.NotificationSettings,
		prov.cfg.TelegramAdminID,
	)

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
	// Единая TZ для всего приложения (cron + сервисы, которые опираются на локальное время).
	prov.schedulerLoc = loc

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
		prov.Services.NotificationSettings,
		loc,
	)

	prov.Services.StockReplenish = stock_replenish.New(
		prov.Repositories.Shop,
		prov.Repositories.Product,
		prov.Repositories.ProductSnapshot,
		prov.Repositories.SaleEvent,
		prov.Repositories.WBOrder,
		prov.Repositories.NotificationEndpoint,
		tgAdapter,
		prov.Services.NotificationSettings,
		prov.log,
	)

	// AI Reply: OpenAI-клиент и сервис генерации ответов на отзывы
	modelCfg := ai.DefaultModelConfig()
	if prov.cfg.AI.ModelForPositive != "" {
		modelCfg.ModelForPositive = prov.cfg.AI.ModelForPositive
	}
	if prov.cfg.AI.ModelForOther != "" {
		modelCfg.ModelForOther = prov.cfg.AI.ModelForOther
	}
	includedLimit := prov.cfg.AI.IncludedLimit
	if includedLimit <= 0 {
		includedLimit = 50
	}
	aiClient, errAI := client.New(modelCfg, prov.cfg.AI.APIKey)
	if errAI != nil {
		prov.log.Warn("ai_reply: OpenAI client init failed, AI replies disabled", "err", errAI)
	} else {
		prov.Services.AIReply = ai_reply.New(
			prov.Repositories.ShopAISettings,
			prov.Repositories.WBFeedback,
			prov.Repositories.AIFeedbackReply,
			prov.Repositories.ShopAIReplyQuota,
			prov.Repositories.ShopAIReplyLedger,
			prov.Repositories.Shop,
			prov.WBClient,
			aiClient,
			modelCfg,
			includedLimit,
			prov.log,
		)
	}

	// Google Auth + JWT для фронта (вход по Google, выдача своего JWT)
	jwtSecret := prov.cfg.JWTSecret
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	googleClientID := prov.cfg.GoogleAuth.ClientID
	if googleClientID == "" {
		googleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	}
	// Сервис авторизации: достаточно JWT_SECRET (email/пароль). Google — опционально.
	if jwtSecret == "" {
		prov.log.Warn("JWT_SECRET not set: auth routes (/api/v1/auth/*) will not be registered; set env JWT_SECRET or config jwt_secret")
	}
	if jwtSecret != "" {
		authGoogleSvc, errAuth := auth_google.NewService(
			prov.log,
			prov.Repositories.User,
			prov.Repositories.LinkTokens,
			prov.Repositories.Shop,
			googleClientID,
			jwtSecret,
			24*time.Hour,
		)
		if errAuth != nil {
			prov.log.Warn("auth service init failed", "err", errAuth)
		} else {
			prov.Services.AuthGoogle = authGoogleSvc
		}
	}
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
