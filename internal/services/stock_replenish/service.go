package stock_replenish

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"

	"salehelp/domain"
	"salehelp/pkg/logger/sl"
)

const (
	daysUrgent = 3
	daysWarning = 7
	lookbackDays = 7
)

type ProductRepo interface {
	ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Product, error)
}

type ProductSnapshotRepo interface {
	GetLatestByProductIDs(ctx context.Context, productIDs []int) (map[int]*domain.ProductSnapshot, error)
}

type SaleEventRepo interface {
	DailyConsumptionByNmID(ctx context.Context, shopID int, since time.Time) (map[int64]float64, error)
}

type ShopRepo interface {
	ListActiveIDs(ctx context.Context, limit, offset int) ([]int, error)
}

type EndpointRepo interface {
	ListEnabledTelegramChatIDsByShopID(ctx context.Context, shopID int) ([]int64, error)
}

type TelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type Service struct {
	log           *slog.Logger
	shops         ShopRepo
	products      ProductRepo
	snapshots     ProductSnapshotRepo
	saleEvents    SaleEventRepo
	endpoints     EndpointRepo
	tg            TelegramSender
}

func New(
	shops ShopRepo,
	products ProductRepo,
	snapshots ProductSnapshotRepo,
	saleEvents SaleEventRepo,
	endpoints EndpointRepo,
	tg TelegramSender,
	log *slog.Logger,
) *Service {
	return &Service{
		log:        log.With(sl.Module("StockReplenish")),
		shops:      shops,
		products:   products,
		snapshots:  snapshots,
		saleEvents: saleEvents,
		endpoints:  endpoints,
		tg:         tg,
	}
}

type productStockInfo struct {
	product   *domain.Product
	stock     int
	daysLeft  float64
	consumption float64
}

// Run запускает проверку остатков по всем активным магазинам и отправляет уведомления.
func (s *Service) Run(ctx context.Context) error {
	const op = "StockReplenish.Run"

	offset := 0
	batchSize := 50
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
				s.log.Error("stock replenish shop failed", "shop_id", shopID, sl.Err(err))
				// продолжаем с другими магазинами
			}
		}

		offset += batchSize
		if len(shopIDs) < batchSize {
			break
		}
	}

	s.log.Info("stock replenish done", "op", op)
	return nil
}

func (s *Service) processShop(ctx context.Context, shopID int) error {
	products, err := s.products.ListByShop(ctx, shopID, 0, 0)
	if err != nil {
		return fmt.Errorf("list products: %w", err)
	}
	if len(products) == 0 {
		return nil
	}

	productIDs := make([]int, 0, len(products))
	for _, p := range products {
		if p != nil {
			productIDs = append(productIDs, p.ID)
		}
	}

	snapshots, err := s.snapshots.GetLatestByProductIDs(ctx, productIDs)
	if err != nil {
		return fmt.Errorf("get snapshots: %w", err)
	}

	since := time.Now().AddDate(0, 0, -lookbackDays)
	consumptionByNmID, err := s.saleEvents.DailyConsumptionByNmID(ctx, shopID, since)
	if err != nil {
		return fmt.Errorf("daily consumption: %w", err)
	}

	var urgent, warning []productStockInfo
	for _, p := range products {
		if p == nil || p.ExternalID == "" {
			continue
		}
		nmID, err := strconv.ParseInt(p.ExternalID, 10, 64)
		if err != nil {
			continue
		}
		snap := snapshots[p.ID]
		if snap == nil {
			continue
		}
		stock := 0
		if snap.StockTotal != nil {
			stock = *snap.StockTotal
		}
		if stock <= 0 {
			continue
		}

		daily := consumptionByNmID[nmID]
		if daily <= 0 {
			// fallback: sales_7d из снимка
			if snap.Sales7d != nil && *snap.Sales7d > 0 {
				daily = float64(*snap.Sales7d) / lookbackDays
			}
		}
		if daily <= 0 {
			continue
		}

		daysLeft := float64(stock) / daily
		info := productStockInfo{
			product:     p,
			stock:       stock,
			daysLeft:    daysLeft,
			consumption: daily,
		}

		if daysLeft < daysUrgent {
			urgent = append(urgent, info)
		} else if daysLeft < daysWarning {
			warning = append(warning, info)
		}
	}

	if len(urgent) == 0 && len(warning) == 0 {
		return nil
	}

	text := s.formatMessage(urgent, warning)
	if text == "" {
		return nil
	}

	chatIDs, err := s.endpoints.ListEnabledTelegramChatIDsByShopID(ctx, shopID)
	if err != nil {
		return fmt.Errorf("list endpoints: %w", err)
	}
	for _, chatID := range chatIDs {
		if err := s.tg.SendMessage(ctx, chatID, text); err != nil {
			s.log.Error("send telegram failed", "shop_id", shopID, "chat_id", chatID, sl.Err(err))
		}
	}

	return nil
}

func (s *Service) formatMessage(urgent, warning []productStockInfo) string {
	var parts []string

	if len(urgent) > 0 {
		lines := formatProductLines(urgent)
		parts = append(parts, "🔴 Срочно пополнить (осталось меньше 3 дней):\n"+lines)
	}
	if len(warning) > 0 {
		lines := formatProductLines(warning)
		parts = append(parts, "🟡 Рекомендуется пополнить (осталось меньше 7 дней):\n"+lines)
	}

	if len(parts) == 0 {
		return ""
	}

	title := "📦 Отчёт по остаткам на складе"
	return title + "\n\n" + strings.Join(parts, "\n\n")
}

func formatProductLines(items []productStockInfo) string {
	lines := make([]string, 0, len(items))
	for _, it := range items {
		days := int(math.Ceil(it.daysLeft))
		if days < 1 {
			days = 1
		}
		name := it.product.Title
		if name == "" {
			name = fmt.Sprintf("nm %s", it.product.ExternalID)
		}
		lines = append(lines, fmt.Sprintf("• %s — остаток %d шт, хватит на ~%d дн.", 
			truncateTitle(name, 50), it.stock, days))
	}
	return strings.Join(lines, "\n")
}

func truncateTitle(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-3]) + "..."
}
