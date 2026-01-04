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
	}

	ShopsRepo interface {
		Create(ctx context.Context, s *domain.Shop) error
		ListByUser(ctx context.Context, userID int, limit, offset int) ([]*domain.Shop, error)
	}

	UsersRepo interface {
		EnsureByTelegramChatID(ctx context.Context, chatID int64) (*domain.User, error)
	}

	ShopNotificationEndpointsRepo interface {
		Exists(ctx context.Context, shopID int, channel, address string) (bool, error)
		Create(ctx context.Context, e *domain.NotificationEndpoint) error
	}

	// MonitoringService — чтобы валидировать WB token или использовать существующую логику.
	MonitoringService interface {
		// ValidateWBToken(ctx, token string) error  // рекомендую добавить
		// или можно временно дернуть SyncShop(...) после создания, но для проверки токена лучше отдельный метод
	}

	Service struct {
		log *slog.Logger

		tg TelegramAPI

		statesRepo TelegramChatStatesRepo

		// на будущее
		invitesRepo   ShopInviteTokensRepo
		usersRepo     UsersRepo
		endpointsRepo ShopNotificationEndpointsRepo
		shopsRepo     ShopsRepo
		monitoring    MonitoringService
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
) *Service {
	return &Service{
		log: log.With(sl.Module("TelegramBotService")),

		tg: tg,

		statesRepo: statesRepo,

		invitesRepo:   invitesRepo,
		usersRepo:     usersRepo,
		endpointsRepo: endpointsRepo,
		shopsRepo:     shopsRepo,
		monitoring:    monitoring,
	}
}
