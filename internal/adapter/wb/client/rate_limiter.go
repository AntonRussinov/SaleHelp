package client

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	RateContent        RateCategory = "content"         // content-api.wildberries.ru
	RateMarketplace    RateCategory = "marketplace"     // marketplace-api.wildberries.ru (склады, остатки — тяжёлый синк)
	RateFBSOrders      RateCategory = "fbs_orders"      // marketplace-api: GET /api/v3/orders/new — лимит 30/час
	RateDBOrders       RateCategory = "db_orders"       // marketplace-api: DBW + DBS orders/new — лимит 10/час
	RatePrices         RateCategory = "prices"          // discounts-prices-api: GET цены (тяжёлый синк)
	RatePricesWrite    RateCategory = "prices_write"    // discounts-prices-api: POST/PUT цены (восстановление цен)
	RateStatistics     RateCategory = "statistics"      // statistics-api.wildberries.ru (sales, orders, stocks)
	RateFeedbacks      RateCategory = "feedbacks"       // feedbacks-api: GET /feedbacks, GET /questions
	RateFeedbacksWrite RateCategory = "feedbacks_write" // feedbacks-api: POST /answer, PATCH /questions
)

type (
	RateCategory string

	RatePolicy struct {
		Limit             int
		Interval          time.Duration
		Burst             int
		PersistAcrossSync bool // не сбрасывать при RegisterKey (для операций вне синка)
	}
)

// basicPolicies — лимиты для базового токена WB (с 16 апреля 2026).
// Документация: https://dev.wildberries.ru/news/281
//
// Burst = кол-во одновременных вызовов одной категории за один синк-цикл.
// RegisterKey сбрасывает лимитеры в начале каждого синка, поэтому Burst не
// приводит к превышению часового лимита WB: за час проходит 2 синка × Burst вызовов.
var basicPolicies = map[RateCategory]RatePolicy{
	// POST /content/v2/get/cards/list — 1/час (инкрементальный синк карточек, тяжёлый синк)
	// Вызывается 1 раз за тяжёлый синк.
	RateContent: {
		Limit:    1,
		Interval: time.Hour,
		Burst:    1,
	},
	// GET /api/v3/stocks/{id}, GET /api/v3/warehouses, GET /api/v3/warehouses/seller —
	// 10/час (тяжёлый синк). За один тяжёлый синк: GetWarehouses + GetSellerWarehouses +
	// GetStocks (×N складов). Burst=4 покрывает типовой магазин (2 склада × 1 чанк).
	RateMarketplace: {
		Limit:    10,
		Interval: time.Hour,
		Burst:    4,
	},
	// GET /api/v3/orders/new — 30/час (быстрый синк FBS-заказов).
	// Вызывается 1 раз за быстрый синк.
	RateFBSOrders: {
		Limit:    30,
		Interval: time.Hour,
		Burst:    1,
	},
	// GET /api/v3/dbw/orders/new + /api/v3/dbs/orders/new — 10/час.
	// За один быстрый синк: 2 вызова (DBW + DBS). Burst=2 → оба проходят сразу.
	// Итого за час: 2 синка × 2 вызова = 4 req/hr << 10.
	RateDBOrders: {
		Limit:    10,
		Interval: time.Hour,
		Burst:    2,
	},
	// POST /api/v2/list/goods/filter — 1/час (тяжёлый синк, чтение цен).
	RatePrices: {
		Limit:    1,
		Interval: time.Hour,
		Burst:    1,
	},
	// POST /api/v2/upload/task — 1/час (запись цен: восстановление по кнопке из Telegram).
	// Отдельная категория чтобы не блокировать на час после тяжёлого синка GetPrices.
	RatePricesWrite: {
		Limit:    1,
		Interval: time.Hour,
		Burst:    1,
	},
	// statistics-api.wildberries.ru — 1/мин (sales + order_refusals в быстром синке).
	// За один быстрый синк: 2 вызова. Burst=2 → оба проходят сразу.
	// Итого за час: 2 синка × 2 вызова = 4 req/hr << 60.
	RateStatistics: {
		Limit:    1,
		Interval: time.Minute,
		Burst:    2,
	},
	// GET /api/v1/feedbacks + GET /api/v1/questions — 5/час.
	// За один быстрый синк: 2 вызова. Burst=2 → оба проходят сразу.
	// Итого за час: 2 синка × 2 вызова = 4 req/hr ≤ 5.
	RateFeedbacks: {
		Limit:    5,
		Interval: time.Hour,
		Burst:    2,
	},
	// POST /api/v1/feedbacks/answer, PATCH /api/v1/questions — 5/час, всплеск 1 (базовый токен).
	// PersistAcrossSync: AI-джоба отправляет ответы отдельно от синка — нельзя сбрасывать на
	// каждом RegisterKey, иначе burst-токен пересоздаётся каждые 30 мин и блокирует retry-логику.
	RateFeedbacksWrite: {
		Limit:             5,
		Interval:          time.Hour,
		Burst:             1,
		PersistAcrossSync: true,
	},
}

// servicePolicies — лимиты для сервисного секрета WB.
// Точные значения уточнить при получении секрета на business-solutions@rwb.ru.
var servicePolicies = map[RateCategory]RatePolicy{
	RateContent:        {Limit: 100, Interval: time.Hour, Burst: 1},
	RateMarketplace:    {Limit: 300, Interval: time.Hour, Burst: 1},
	RateFBSOrders:      {Limit: 300, Interval: time.Hour, Burst: 1},
	RateDBOrders:       {Limit: 300, Interval: time.Hour, Burst: 1},
	RatePrices:         {Limit: 300, Interval: time.Hour, Burst: 1},
	RatePricesWrite:    {Limit: 300, Interval: time.Hour, Burst: 1},
	RateStatistics:     {Limit: 60, Interval: time.Minute, Burst: 1},
	RateFeedbacks:      {Limit: 300, Interval: time.Hour, Burst: 1},
	RateFeedbacksWrite: {Limit: 300, Interval: time.Hour, Burst: 1},
}

type LimiterManager struct {
	mu          sync.Mutex
	keyPolicies map[string]map[RateCategory]RatePolicy  // per-key overrides
	lim         map[string]map[RateCategory]*rate.Limiter
}

func NewLimiterManager() *LimiterManager {
	return &LimiterManager{
		keyPolicies: make(map[string]map[RateCategory]RatePolicy),
		lim:         make(map[string]map[RateCategory]*rate.Limiter),
	}
}

// RegisterKey устанавливает набор политик для конкретного API-ключа.
// Вызывать при старте синка магазина — сбрасывает лимитеры для ключа, кроме тех у которых
// PersistAcrossSync=true (они используются вне синка и не должны терять накопленные токены).
func (m *LimiterManager) RegisterKey(apiKey, tokenType string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	policies := basicPolicies
	if tokenType == "service" {
		policies = servicePolicies
	}
	m.keyPolicies[apiKey] = policies

	if existing, ok := m.lim[apiKey]; ok {
		// Удаляем только те категории, которым можно сбрасываться при синке.
		for cat := range existing {
			if p, found := policies[cat]; !found || !p.PersistAcrossSync {
				delete(existing, cat)
			}
		}
	} else {
		m.lim[apiKey] = make(map[RateCategory]*rate.Limiter)
	}
}

func (m *LimiterManager) Get(apiKey string, cat RateCategory) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.lim[apiKey]; !ok {
		m.lim[apiKey] = make(map[RateCategory]*rate.Limiter)
	}
	if l, ok := m.lim[apiKey][cat]; ok {
		return l
	}

	policies, ok := m.keyPolicies[apiKey]
	if !ok {
		policies = basicPolicies
	}
	p := policies[cat]
	interval := p.Interval / time.Duration(p.Limit)
	l := rate.NewLimiter(rate.Every(interval), p.Burst)

	m.lim[apiKey][cat] = l
	return l
}
