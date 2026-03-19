package telegrambot

import (
	"context"
	"log/slog"
	"time"

	"salehelp/domain"
	"salehelp/internal/services/notification_settings"
	"salehelp/internal/services/price_restore"
	"salehelp/internal/services/subscription"
	"salehelp/internal/storage/pggorm/model"
	"salehelp/pkg/logger/sl"
)

type (
	// TelegramAPI — низкоуровневый клиент Telegram (sendMessage, inline-кнопки, callback).
	TelegramAPI interface {
		SendMessage(ctx context.Context, chatID int64, text string) error
		SendMessageWithKeyboard(ctx context.Context, chatID int64, text string, keyboard map[string]any) error
		SendMessageWithInlineKeyboard(ctx context.Context, chatID int64, text string, inlineKeyboard map[string]any) error
		AnswerCallbackQuery(ctx context.Context, callbackQueryID string) error
	}

	// TelegramChatStatesRepo — твой репозиторий состояний мастера.
	TelegramChatStatesRepo interface {
		GetByChatID(ctx context.Context, chatID int64) (*domain.TelegramChatState, error)
		Upsert(ctx context.Context, st *domain.TelegramChatState, now time.Time) error // если у тебя иначе — смени сигнатуру
		Clear(ctx context.Context, chatID int64) error
	}

	// ShopInviteTokensRepo — на будущее (/invite, /join). Пока можно не использовать, но зависимость заложим.
	ShopInviteTokensRepo interface {
		Create(ctx context.Context, token *domain.ShopInviteToken) error
		GetActiveByShopID(ctx context.Context, shopID int) (*domain.ShopInviteToken, error)
		GetByToken(ctx context.Context, token string) (*domain.ShopInviteToken, error)
	}

	ShopsRepo interface {
		Create(ctx context.Context, s *domain.Shop) error
		GetByID(ctx context.Context, id int) (*domain.Shop, error)
		GetLastSyncAtMap(ctx context.Context, shopIDs []int) (map[int]*time.Time, error)
		ListByUser(ctx context.Context, userID int, limit, offset int) ([]*domain.Shop, error)
		ListByTelegramChat(ctx context.Context, chatID int64) ([]*domain.Shop, error)
		SetActive(ctx context.Context, shopID int, active bool) error
		ReassignShopsToUser(ctx context.Context, fromUserID, toUserID int) error

		CountAll(ctx context.Context) (int, error)
		CountActive(ctx context.Context) (int, error)
		CountByLastSyncStatusCode(ctx context.Context, code int) (int, error)
	}

	UsersRepo interface {
		EnsureByTelegramChatID(ctx context.Context, chatID int64) (*domain.User, error)
		GetByTelegramChatID(ctx context.Context, chatID int64) (*domain.User, error)
		GetByEmail(ctx context.Context, email string) (*domain.User, error)
		SetTelegramChatID(ctx context.Context, userID int, chatID int64) error
		ClearTelegramChatID(ctx context.Context, userID int) error
		Update(ctx context.Context, u *domain.User) error
	}

	LinkTokensRepo interface {
		Consume(ctx context.Context, token string, now time.Time) (userID int, err error)
		CreateFromTelegram(ctx context.Context, chatID int64) (code string, expiresAt time.Time, err error)
	}

	ShopNotificationEndpointsRepo interface {
		Exists(ctx context.Context, shopID int, channel, address string) (bool, error)
		Create(ctx context.Context, e *domain.NotificationEndpoint) error
		SetEnabledByShopAndAddress(ctx context.Context, shopID int, channel, address string, enabled bool) error
		GetByShopAndAddress(ctx context.Context, shopID int, channel string, address string) (*domain.NotificationEndpoint, error)
		ListByShopID(ctx context.Context, shopID int, onlyEnabled bool) ([]*domain.NotificationEndpoint, error)
		CountEnabledByShopAndChannel(ctx context.Context, shopID int, channel string) (int, error)
	}

	// MonitoringService — чтобы валидировать WB token или использовать существующую логику.
	MonitoringService interface {
		CheckWBToken(ctx context.Context, token string) error
	}

	SubscriptionEntitlementsService interface {
		// возвращает allowed telegram chats для магазина
		AllowedTelegramChats(ctx context.Context, shopID int) (int, error)
		GetShopAccessLevel(ctx context.Context, shopID int) (subscription.AccessLevel, *domain.Subscription, error)
	}
	SubscriptionRepo interface {
		CountByState(ctx context.Context, now time.Time) (active, grace, expired int, err error)
		SumRevenueActiveByCurrency(ctx context.Context, now time.Time) ([]model.MoneyByCurrency, error)
		SumRevenueGraceByCurrency(ctx context.Context, now time.Time) ([]model.MoneyByCurrency, error)
	}

	// ABCReportService — для команды /abc_now (отчёт по магазину за 30 дней).
	ABCReportService interface {
		SaveSnapshot(ctx context.Context, shopID, periodDays int) (*domain.ABCSnapshot, error)
		GetTwoLatestSnapshots(ctx context.Context, shopID, periodDays int) ([]*domain.ABCSnapshot, error)
	}

	ProductsRepoForABC interface {
		ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Product, error)
	}

	// PriceRestoreService — восстановление цен по пакету (кнопка «Вернуть цены назад»).
	PriceRestoreService interface {
		Restore(ctx context.Context, batchID int) price_restore.RestoreResult
	}

	// PriceRestoreBatchRepo — проверка restored_at при нажатии «Вернуть цены назад».
	PriceRestoreBatchRepo interface {
		GetByID(ctx context.Context, id int) (*domain.PriceRestoreBatch, error)
	}

	Service struct {
		log              *slog.Logger
		tg               TelegramAPI
		statesRepo       TelegramChatStatesRepo
		invitesRepo      ShopInviteTokensRepo
		usersRepo        UsersRepo
		linkTokensRepo   LinkTokensRepo
		endpointsRepo    ShopNotificationEndpointsRepo
		shopsRepo        ShopsRepo
		monitoring       MonitoringService
		subEntitlements  SubscriptionEntitlementsService
		subscriptionRepo SubscriptionRepo
		abcSvc           ABCReportService
		productsRepo     ProductsRepoForABC
		priceRestore     PriceRestoreService
		restoreBatchRepo PriceRestoreBatchRepo
		notifSettings    *notification_settings.Service
		botUsername      string
		adminChatID      int64
	}
)

func NewService(
	log *slog.Logger,
	tg TelegramAPI,
	statesRepo TelegramChatStatesRepo,
	invitesRepo ShopInviteTokensRepo,
	usersRepo UsersRepo,
	linkTokensRepo LinkTokensRepo,
	endpointsRepo ShopNotificationEndpointsRepo,
	shopsRepo ShopsRepo,
	monitoring MonitoringService,
	botUsername string,
	subEntitlements SubscriptionEntitlementsService,
	subscriptionRepo SubscriptionRepo,
	abcSvc ABCReportService,
	productsRepo ProductsRepoForABC,
	priceRestore PriceRestoreService,
	restoreBatchRepo PriceRestoreBatchRepo,
	notifSettings *notification_settings.Service,
	adminChatID int64,
) *Service {
	return &Service{
		log:              log.With(sl.Module("TelegramBotService")),
		tg:               tg,
		statesRepo:       statesRepo,
		invitesRepo:      invitesRepo,
		usersRepo:        usersRepo,
		linkTokensRepo:   linkTokensRepo,
		endpointsRepo:    endpointsRepo,
		shopsRepo:        shopsRepo,
		monitoring:       monitoring,
		botUsername:      botUsername,
		subEntitlements:  subEntitlements,
		subscriptionRepo: subscriptionRepo,
		abcSvc:           abcSvc,
		productsRepo:     productsRepo,
		priceRestore:     priceRestore,
		restoreBatchRepo: restoreBatchRepo,
		notifSettings:    notifSettings,
		adminChatID:      adminChatID,
	}
}
