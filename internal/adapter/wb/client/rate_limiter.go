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

	LimiterManager struct {
		mu       sync.Mutex
		policies map[RateCategory]RatePolicy
		lim      map[string]map[RateCategory]*rate.Limiter // apiKey -> category -> limiter
	}
)

// Лимиты для базового токена WB (вступили в силу с 16 апреля 2026).
// Документация: https://dev.wildberries.ru/news/281
// Burst=1 везде — WB не разрешает всплески для базового токена.
// После получения сервисного секрета (заявка на business-solutions@rwb.ru)
// лимиты можно вернуть к прежним значениям.
var defaultPolicies = map[RateCategory]RatePolicy{
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

func NewLimiterManager(policies map[RateCategory]RatePolicy) *LimiterManager {
	return &LimiterManager{
		policies: policies,
		lim:      make(map[string]map[RateCategory]*rate.Limiter),
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

	p := m.policies[cat]

	interval := m.policies[cat].Interval / time.Duration(m.policies[cat].Limit)
	l := rate.NewLimiter(rate.Every(interval), p.Burst)

	m.lim[apiKey][cat] = l
	return l
}
