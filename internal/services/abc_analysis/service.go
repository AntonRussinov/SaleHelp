package abc_analysis

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"time"

	"salehelp/domain"
	saleevent "salehelp/internal/storage/pggorm/repository/saleevent"
	"salehelp/pkg/logger/sl"
)

// Пороги накопительного процента выручки для категорий A/B/C (Парето).
const (
	thresholdA = 80
	thresholdB = 95
)

type SaleEventRepo interface {
	RevenueAndCountByNmIDInPeriod(ctx context.Context, shopID int, from, to time.Time) ([]saleevent.RevenueByNmIDRow, error)
}

type ProductRepo interface {
	ListByShop(ctx context.Context, shopID int, limit, offset int) ([]*domain.Product, error)
}

type ABCSnapshotRepo interface {
	CreateBatch(ctx context.Context, rows []domain.ABCSnapshotRow) error
	GetLatestByShopAndPeriod(ctx context.Context, shopID, periodDays int) (*domain.ABCSnapshot, error)
	ListTwoLatestByShopAndPeriod(ctx context.Context, shopID, periodDays int) ([]*domain.ABCSnapshot, error)
}

type Service struct {
	log     *slog.Logger
	sales   SaleEventRepo
	products ProductRepo
	snapshots ABCSnapshotRepo
}

func New(
	sales SaleEventRepo,
	products ProductRepo,
	snapshots ABCSnapshotRepo,
	log *slog.Logger,
) *Service {
	return &Service{
		log:       log.With(sl.Module("ABCAnalysis")),
		sales:     sales,
		products:  products,
		snapshots: snapshots,
	}
}

// rowForSort — промежуточная структура для сортировки и расчёта категорий.
type rowForSort struct {
	nmID       int64
	productID  *int
	revenue    float64
	salesCount int
}

// CalculateABC считает ABC на основе выкупов (sale_events, event_type = 'buyout') за последние periodDays дней.
// Товары без продаж попадают в C. Возвращает снэпшот без сохранения в БД.
func (s *Service) CalculateABC(ctx context.Context, shopID, periodDays int) (*domain.ABCSnapshot, error) {
	periodEnd := time.Now().UTC()
	periodStart := periodEnd.AddDate(0, 0, -periodDays)

	salesRows, err := s.sales.RevenueAndCountByNmIDInPeriod(ctx, shopID, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("revenue by nm_id: %w", err)
	}

	products, err := s.products.ListByShop(ctx, shopID, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}

	nmToProductID := make(map[int64]int)
	nmToSubject := make(map[int64]string)
	for _, p := range products {
		nmID, err := strconv.ParseInt(p.ExternalID, 10, 64)
		if err != nil {
			continue
		}
		nmToProductID[nmID] = p.ID
		nmToSubject[nmID] = p.Title
	}

	inSales := make(map[int64]bool)
	var rows []rowForSort
	for _, r := range salesRows {
		inSales[r.NmID] = true
		pid, _ := nmToProductID[r.NmID]
		var pidPtr *int
		if pid != 0 {
			pidPtr = &pid
		}
		subject := r.Subject
		if subject == "" {
			subject = nmToSubject[r.NmID]
		}
		_ = subject
		rows = append(rows, rowForSort{nmID: r.NmID, productID: pidPtr, revenue: r.Revenue, salesCount: r.SalesCount})
	}
	for _, p := range products {
		nmID, err := strconv.ParseInt(p.ExternalID, 10, 64)
		if err != nil || inSales[nmID] {
			continue
		}
		pid := p.ID
		rows = append(rows, rowForSort{nmID: nmID, productID: &pid, revenue: 0, salesCount: 0})
	}

	sort.Slice(rows, func(i, j int) bool { return rows[i].revenue > rows[j].revenue })

	var totalRevenue float64
	for _, r := range rows {
		totalRevenue += r.revenue
	}

	snap := &domain.ABCSnapshot{
		PeriodDays:   periodDays,
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
		TotalRevenue: totalRevenue,
		Rows:         make([]domain.ABCSnapshotRow, 0, len(rows)),
	}

	if totalRevenue == 0 {
		for _, r := range rows {
			snap.Rows = append(snap.Rows, domain.ABCSnapshotRow{
				ShopID:     shopID,
				ProductID:  r.productID,
				NmID:       r.nmID,
				Category:   domain.ABCCategoryC,
				Revenue:    r.revenue,
				SalesCount: r.salesCount,
				PeriodDays: periodDays,
				PeriodStart: periodStart,
				PeriodEnd:   periodEnd,
			})
		}
		return snap, nil
	}

	cumulative := 0.0
	for _, r := range rows {
		cumulative += r.revenue
		pct := cumulative / totalRevenue * 100
		cat := domain.ABCCategoryC
		if pct <= thresholdA {
			cat = domain.ABCCategoryA
		} else if pct <= thresholdB {
			cat = domain.ABCCategoryB
		}
		snap.Rows = append(snap.Rows, domain.ABCSnapshotRow{
			ShopID:     shopID,
			ProductID:  r.productID,
			NmID:       r.nmID,
			Category:   cat,
			Revenue:    r.revenue,
			SalesCount: r.salesCount,
			PeriodDays: periodDays,
			PeriodStart: periodStart,
			PeriodEnd:   periodEnd,
		})
	}

	return snap, nil
}

// SaveSnapshot вызывает CalculateABC и сохраняет результат в abc_snapshots (immutable INSERT).
func (s *Service) SaveSnapshot(ctx context.Context, shopID, periodDays int) (*domain.ABCSnapshot, error) {
	snap, err := s.CalculateABC(ctx, shopID, periodDays)
	if err != nil {
		return nil, err
	}
	if len(snap.Rows) == 0 {
		return snap, nil
	}
	createdAt := time.Now().UTC()
	for i := range snap.Rows {
		snap.Rows[i].CreatedAt = createdAt
	}
	rows := make([]domain.ABCSnapshotRow, len(snap.Rows))
	copy(rows, snap.Rows)
	if err := s.snapshots.CreateBatch(ctx, rows); err != nil {
		return nil, fmt.Errorf("save abc snapshot: %w", err)
	}
	snap.CreatedAt = createdAt
	s.log.Info("abc snapshot saved", "shop_id", shopID, "period_days", periodDays, "rows", len(snap.Rows))
	return snap, nil
}

// GetLatestSnapshot возвращает последний сохранённый снэпшот для магазина и периода.
func (s *Service) GetLatestSnapshot(ctx context.Context, shopID, periodDays int) (*domain.ABCSnapshot, error) {
	return s.snapshots.GetLatestByShopAndPeriod(ctx, shopID, periodDays)
}

// GetTwoLatestSnapshots возвращает два последних снэпшота (новый первый, предыдущий второй) для сравнения.
func (s *Service) GetTwoLatestSnapshots(ctx context.Context, shopID, periodDays int) ([]*domain.ABCSnapshot, error) {
	return s.snapshots.ListTwoLatestByShopAndPeriod(ctx, shopID, periodDays)
}

// CompareSnapshots возвращает изменения категорий между двумя последними снэпшотами.
// Требуется подставить названия товаров (subject) снаружи по NmID/ProductID.
func (s *Service) CompareSnapshots(ctx context.Context, shopID, periodDays int) ([]domain.ABCChange, error) {
	two, err := s.snapshots.ListTwoLatestByShopAndPeriod(ctx, shopID, periodDays)
	if err != nil {
		return nil, err
	}
	if len(two) < 2 {
		return nil, nil
	}
	prev := two[1]
	curr := two[0]
	prevByNm := make(map[int64]domain.ABCCategory)
	for _, row := range prev.Rows {
		prevByNm[row.NmID] = row.Category
	}
	var changes []domain.ABCChange
	for _, row := range curr.Rows {
		from, ok := prevByNm[row.NmID]
		if !ok || from == row.Category {
			continue
		}
		var changePct float64
		if prev.TotalRevenue > 0 {
			prevRow := findRowByNmID(prev.Rows, row.NmID)
			if prevRow != nil && prevRow.Revenue > 0 {
				changePct = (row.Revenue - prevRow.Revenue) / prevRow.Revenue * 100
			}
		}
		changes = append(changes, domain.ABCChange{
			NmID:      row.NmID,
			ProductID: row.ProductID,
			FromCat:   from,
			ToCat:     row.Category,
			Revenue:   row.Revenue,
			ChangePct: changePct,
		})
	}
	return changes, nil
}

func findRowByNmID(rows []domain.ABCSnapshotRow, nmID int64) *domain.ABCSnapshotRow {
	for i := range rows {
		if rows[i].NmID == nmID {
			return &rows[i]
		}
	}
	return nil
}
