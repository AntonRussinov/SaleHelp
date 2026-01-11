package telegrambot

import (
	"context"
	"log/slog"
	"time"

	"salehelp/domain" // или куда ты положил Update DTO
	"salehelp/pkg/logger/sl"
)

type (
	// TelegramAPI — низкоуровневый клиент Telegram (sendMessage).
	TelegramAPI interface {
		SendMessage(ctx context.Context, chatID int64, text string) error
		SendMessageWithKeyboard(ctx context.Context, chatID int64, text string, keyboard map[string]any) error
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
		ListByUser(ctx context.Context, userID int, limit, offset int) ([]*domain.Shop, error)
		ListByTelegramChat(ctx context.Context, chatID int64) ([]*domain.Shop, error)
		SetActive(ctx context.Context, shopID int, active bool) error
	}

	UsersRepo interface {
		EnsureByTelegramChatID(ctx context.Context, chatID int64) (*domain.User, error)
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
	}

	Service struct {
		log             *slog.Logger
		tg              TelegramAPI
		statesRepo      TelegramChatStatesRepo
		invitesRepo     ShopInviteTokensRepo
		usersRepo       UsersRepo
		endpointsRepo   ShopNotificationEndpointsRepo
		shopsRepo       ShopsRepo
		monitoring      MonitoringService
		subEntitlements SubscriptionEntitlementsService
		botUsername     string
	}
)

func NewService(
	log *slog.Logger,
	tg TelegramAPI,
	statesRepo TelegramChatStatesRepo,
	invitesRepo ShopInviteTokensRepo,
	usersRepo UsersRepo,
	endpointsRepo ShopNotificationEndpointsRepo,
	shopsRepo ShopsRepo,
	monitoring MonitoringService,
	botUsername string,
	subEntitlements SubscriptionEntitlementsService,
) *Service {
	return &Service{
		log:             log.With(sl.Module("TelegramBotService")),
		tg:              tg,
		statesRepo:      statesRepo,
		invitesRepo:     invitesRepo,
		usersRepo:       usersRepo,
		endpointsRepo:   endpointsRepo,
		shopsRepo:       shopsRepo,
		monitoring:      monitoring,
		botUsername:     botUsername,
		subEntitlements: subEntitlements,
	}
}
