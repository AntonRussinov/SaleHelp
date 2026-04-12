package refusal_digest

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"salehelp/domain"
	notifications "salehelp/internal/services/notification"
	"salehelp/internal/services/notification_settings"
	"salehelp/pkg/logger/sl"
)

// ShopRepo перечисляет активные магазины.
type ShopRepo interface {
	ListActiveIDs(ctx context.Context, limit, offset int) ([]int, error)
}

// OrderRefusalRepo читает и помечает отказы.
type OrderRefusalRepo interface {
	ListUnnotifiedByShop(ctx context.Context, shopID int, limit int) ([]*domain.OrderRefusal, error)
	MarkNotified(ctx context.Context, ids []int, at time.Time) error
}

// ProductRepo загружает товары магазина для подстановки названий.
type ProductRepo interface {
	ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Product, error)
}

// EndpointRepo отдаёт chat_id Telegram-эндпоинтов магазина.
type EndpointRepo interface {
	ListEnabledTelegramChatIDsByShopID(ctx context.Context, shopID int) ([]int64, error)
}

// TelegramSender отправляет сообщение в чат (с автоматическим разбиением длинных сообщений).
type TelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

// Service — суточный дайджест отказов (заказы не забраны).
// Запускается раз в сутки (07:00 МСК), собирает все необработанные отказы по каждому магазину
// и отправляет одним сообщением (или несколькими, если данных много).
type Service struct {
	log       *slog.Logger
	shops     ShopRepo
	refusals  OrderRefusalRepo
	products  ProductRepo
	endpoints EndpointRepo
	tg        TelegramSender
	settings  *notification_settings.Service
	loc       *time.Location
}

func New(
	shops ShopRepo,
	refusals OrderRefusalRepo,
	products ProductRepo,
	endpoints EndpointRepo,
	tg TelegramSender,
	settings *notification_settings.Service,
	loc *time.Location,
	log *slog.Logger,
) *Service {
	if loc == nil {
		loc = time.UTC
	}
	return &Service{
		log:       log.With(sl.Module("RefusalDigest")),
		shops:     shops,
		refusals:  refusals,
		products:  products,
		endpoints: endpoints,
		tg:        tg,
		settings:  settings,
		loc:       loc,
	}
}

// Run обходит все активные магазины и отправляет дайджест отказов.
func (s *Service) Run(ctx context.Context) error {
	const op = "RefusalDigest.Run"

	offset := 0
	const batchSize = 50
	for {
		shopIDs, err := s.shops.ListActiveIDs(ctx, batchSize, offset)
		if err != nil {
			return fmt.Errorf("%s: list shops: %w", op, err)
		}
		if len(shopIDs) == 0 {
			break
		}
		for _, shopID := range shopIDs {
			if err := s.processShop(ctx, shopID); err != nil {
				s.log.Error("refusal digest shop failed", "shop_id", shopID, sl.Err(err))
			}
		}
		offset += batchSize
		if len(shopIDs) < batchSize {
			break
		}
	}

	s.log.Info("refusal digest done", "op", op)
	return nil
}

func (s *Service) processShop(ctx context.Context, shopID int) error {
	refs, err := s.refusals.ListUnnotifiedByShop(ctx, shopID, 200)
	if err != nil {
		return fmt.Errorf("list refusals: %w", err)
	}
	if len(refs) == 0 {
		return nil
	}

	now := time.Now()
	ids := make([]int, 0, len(refs))
	for _, r := range refs {
		ids = append(ids, r.ID)
	}

	// Если уведомления отключены — просто помечаем как отправленные, не шлём.
	if s.settings != nil {
		enabled, err := s.settings.EffectiveEnabled(ctx, shopID, domain.NotificationEventSalesAndRefusals)
		if err == nil && !enabled {
			return s.refusals.MarkNotified(ctx, ids, now)
		}
	}

	prods, err := s.products.ListByShop(ctx, shopID, 0, 0)
	if err != nil {
		return fmt.Errorf("list products: %w", err)
	}
	titleByNmID := make(map[int64]string, len(prods))
	for _, p := range prods {
		if p == nil || p.ExternalID == "" || p.Title == "" {
			continue
		}
		nmID, parseErr := strconv.ParseInt(p.ExternalID, 10, 64)
		if parseErr != nil {
			continue
		}
		titleByNmID[nmID] = p.Title
	}

	text := notifications.FormatOrderRefusalsMessageRU(refs, titleByNmID, s.loc)
	if text == "" {
		return nil
	}

	chatIDs, err := s.endpoints.ListEnabledTelegramChatIDsByShopID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("list endpoints: %w", err)
	}

	sent := false
	for _, chatID := range chatIDs {
		if err := s.tg.SendMessage(ctx, chatID, text); err != nil {
			s.log.Error("refusal digest send failed", "shop_id", shopID, "chat_id", chatID, sl.Err(err))
			continue
		}
		sent = true
	}

	if sent {
		return s.refusals.MarkNotified(ctx, ids, now)
	}
	return nil
}
