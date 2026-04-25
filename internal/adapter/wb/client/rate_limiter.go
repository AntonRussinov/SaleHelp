package client

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	RateContent     RateCategory = "content"     // content-api.wildberries.ru
	RateMarketplace RateCategory = "marketplace" // marketplace-api.wildberries.ru (склады, остатки — тяжёлый синк)
	RateFBSOrders   RateCategory = "fbs_orders"  // marketplace-api: GET /api/v3/orders/new — лимит 30/час
	RateDBOrders    RateCategory = "db_orders"   // marketplace-api: DBW + DBS orders/new — лимит 10/час
	RatePrices      RateCategory = "prices"      // discounts-prices-api.wildberries.ru
	RateStatistics  RateCategory = "statistics"  // statistics-api.wildberries.ru (sales, orders, stocks)
	RateFeedbacks      RateCategory = "feedbacks"       // feedbacks-api: GET /feedbacks, GET /questions
	RateFeedbacksWrite RateCategory = "feedbacks_write" // feedbacks-api: POST /answer, PATCH /questions
)

type (
	RateCategory string

	RatePolicy struct {
		Limit    int
		Interval time.Duration
		Burst    int
	}
)

// basicPolicies — лимиты для базового токена WB (с 16 апреля 2026).
// Документация: https://dev.wildberries.ru/news/281
var basicPolicies = map[RateCategory]RatePolicy{
	// POST /content/v2/get/cards/list — 1/час (инкрементальный синк карточек, тяжёлый синк)
	RateContent: {
		Limit:    1,
		Interval: time.Hour,
		Burst:    1,
	},
	// GET /api/v3/stocks/{id}, GET /api/v3/warehouses — 10/час (тяжёлый синк)
	RateMarketplace: {
		Limit:    10,
		Interval: time.Hour,
		Burst:    1,
	},
	// GET /api/v3/orders/new — 30/час (быстрый синк FBS-заказов)
	RateFBSOrders: {
		Limit:    30,
		Interval: time.Hour,
		Burst:    1,
	},
	// GET /api/v3/dbw/orders/new, /api/v3/dbs/orders/new — 10/час (быстрый синк)
	RateDBOrders: {
		Limit:    10,
		Interval: time.Hour,
		Burst:    1,
	},
	// POST /api/v2/list/goods/filter, upload/task — 1/час (тяжёлый синк цен)
	RatePrices: {
		Limit:    1,
		Interval: time.Hour,
		Burst:    1,
	},
	// statistics-api.wildberries.ru — не вошёл в новую таблицу лимитов WB, оставляем консервативно
	RateStatistics: {
		Limit:    1,
		Interval: time.Minute,
		Burst:    1,
	},
	// GET /api/v1/feedbacks, GET /api/v1/questions — 5/час (чтение, базовый токен WB)
	RateFeedbacks: {
		Limit:    5,
		Interval: time.Hour,
		Burst:    1,
	},
	// POST /api/v1/feedbacks/answer, PATCH /api/v1/questions — 5/час, burst 1 (базовый токен WB)
	RateFeedbacksWrite: {
		Limit:    5,
		Interval: time.Hour,
		Burst:    1,
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
// Вызывать при старте синка магазина — сбрасывает существующие лимитеры для ключа.
func (m *LimiterManager) RegisterKey(apiKey, tokenType string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	policies := basicPolicies
	if tokenType == "service" {
		policies = servicePolicies
	}
	m.keyPolicies[apiKey] = policies
	delete(m.lim, apiKey) // сбросить существующие лимитеры — пересоздадутся с новыми политиками
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
