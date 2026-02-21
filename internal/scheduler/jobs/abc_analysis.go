package jobs

import (
	"context"
	"fmt"
	"sync/atomic"

	"salehelp/domain"
	"salehelp/internal/services/abc_analysis"
	"salehelp/internal/services/abc_analysis/notify"
	"salehelp/internal/services/subscription"
)

type ABCShopsRepo interface {
	ListActiveIDs(ctx context.Context, limit, offset int) ([]int, error)
}

type ABCSnapshotRepo interface {
	GetLatestByShopAndPeriod(ctx context.Context, shopID, periodDays int) (*domain.ABCSnapshot, error)
	ListTwoLatestByShopAndPeriod(ctx context.Context, shopID, periodDays int) ([]*domain.ABCSnapshot, error)
}

type ShopABCSettingsRepo interface {
	GetByShopID(ctx context.Context, shopID int) (*domain.ShopABCSettings, error)
}

type ProductRepo interface {
	ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Product, error)
}

type ABCSubscriptionService interface {
	GetShopAccessLevel(ctx context.Context, shopID int) (subscription.AccessLevel, *domain.Subscription, error)
}

type ABCEndpointRepo interface {
	ListEnabledTelegramChatIDsByShopID(ctx context.Context, shopID int) ([]int64, error)
}

type ABCTelegramSender interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}

type ABCJobConfig struct {
	PeriodDays       int
	MinSalesThreshold int
	IsWeekly         bool
}

type ABCAnalysisJob struct {
	log          Logger
	cfg          ABCJobConfig
	shops        ABCShopsRepo
	abcSvc       *abc_analysis.Service
	settingsRepo ShopABCSettingsRepo
	subSvc       ABCSubscriptionService
	endpoints    ABCEndpointRepo
	tg           ABCTelegramSender
	products     ProductRepo
}

func NewABCAnalysisJob(
	log Logger,
	cfg ABCJobConfig,
	shops ABCShopsRepo,
	abcSvc *abc_analysis.Service,
	settingsRepo ShopABCSettingsRepo,
	subSvc ABCSubscriptionService,
	endpoints ABCEndpointRepo,
	tg ABCTelegramSender,
	products ProductRepo,
) *ABCAnalysisJob {
	if cfg.MinSalesThreshold <= 0 {
		cfg.MinSalesThreshold = 10
	}
	return &ABCAnalysisJob{
		log:          log,
		cfg:          cfg,
		shops:        shops,
		abcSvc:       abcSvc,
		settingsRepo: settingsRepo,
		subSvc:       subSvc,
		endpoints:    endpoints,
		tg:           tg,
		products:     products,
	}
}

func (j *ABCAnalysisJob) Run(ctx context.Context) {
	offset := 0
	const batchSize = 50
	var doneCount int64

	for {
		select {
		case <-ctx.Done():
			j.log.Warn("abc_analysis job: canceled", "err", ctx.Err())
			return
		default:
		}

		ids, err := j.shops.ListActiveIDs(ctx, batchSize, offset)
		if err != nil {
			j.log.Error("abc_analysis job: list shops failed", "err", err)
			return
		}
		if len(ids) == 0 {
			break
		}

		for _, shopID := range ids {
			lvl, _, err := j.subSvc.GetShopAccessLevel(ctx, shopID)
			if err != nil {
				j.log.Error("abc_analysis job: get access level failed", "shop_id", shopID, "err", err)
				continue
			}
			if lvl != subscription.AccessActive {
				continue
			}
			settings, _ := j.settingsRepo.GetByShopID(ctx, shopID)
			weeklyEnabled := settings == nil || settings.ABCWeeklyEnabled
			monthlyEnabled := settings == nil || settings.ABCMonthlyEnabled
			if j.cfg.IsWeekly && !weeklyEnabled {
				continue
			}
			if !j.cfg.IsWeekly && !monthlyEnabled {
				continue
			}

			snap, err := j.abcSvc.SaveSnapshot(ctx, shopID, j.cfg.PeriodDays)
			if err != nil {
				j.log.Error("abc_analysis job: save snapshot failed", "shop_id", shopID, "err", err)
				continue
			}

			totalSales := 0
			for _, row := range snap.Rows {
				totalSales += row.SalesCount
			}
			if totalSales < j.cfg.MinSalesThreshold {
				j.log.Info("abc_analysis job: skip report, too few sales", "shop_id", shopID, "total_sales", totalSales)
				atomic.AddInt64(&doneCount, 1)
				continue
			}

			products, _ := j.products.ListByShop(ctx, shopID, 0, 0)
			productByNm := make(map[int64]*domain.Product)
			for _, p := range products {
				var nmID int64
				if n, err := parseInt64(p.ExternalID); err == nil {
					nmID = n
					productByNm[nmID] = p
				}
			}

			var msg string
			if j.cfg.IsWeekly {
				changes, _ := j.abcSvc.CompareSnapshots(ctx, shopID, j.cfg.PeriodDays)
				for i := range changes {
					if p := productByNm[changes[i].NmID]; p != nil {
						changes[i].Subject = p.Title
					}
				}
				msg = notify.FormatWeeklyDigestRU(snap, changes)
			} else {
				two, _ := j.abcSvc.GetTwoLatestSnapshots(ctx, shopID, j.cfg.PeriodDays)
				var prev *domain.ABCSnapshot
				if len(two) >= 2 {
					prev = two[1]
				}
				msg = notify.FormatMonthlyReportRU(snap, prev, productByNm)
			}

			chatIDs, err := j.endpoints.ListEnabledTelegramChatIDsByShopID(ctx, shopID)
			if err != nil || len(chatIDs) == 0 {
				j.log.Warn("abc_analysis job: no telegram endpoints", "shop_id", shopID)
				atomic.AddInt64(&doneCount, 1)
				continue
			}
			for _, chatID := range chatIDs {
				if err := j.tg.SendMessage(ctx, chatID, msg); err != nil {
					j.log.Error("abc_analysis job: send failed", "shop_id", shopID, "chat_id", chatID, "err", err)
				}
			}
			atomic.AddInt64(&doneCount, 1)
		}

		offset += len(ids)
	}

	j.log.Info("abc_analysis job: finished", "period_days", j.cfg.PeriodDays, "shops_processed", doneCount)
}

func parseInt64(s string) (int64, error) {
	var n int64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
