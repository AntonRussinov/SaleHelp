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
var basicPolicies = map[RateCategory]RatePolicy{
	// POST /content/v2/get/cards/list — 1/час.
	RateContent: {Limit: 1, Interval: time.Hour, Burst: 1},
	// GET /api/v3/stocks/{id}, /warehouses/seller — 10/час.
	// Burst=4: warehouses + stocks×N складов за один тяжёлый синк.
	RateMarketplace: {Limit: 10, Interval: time.Hour, Burst: 4},
	// GET /api/v3/orders/new — 30/час.
	RateFBSOrders: {Limit: 30, Interval: time.Hour, Burst: 1},
	// GET /api/v3/dbw/orders/new + /dbs/orders/new — 10/час. Burst=2: оба вызова за синк.
	RateDBOrders: {Limit: 10, Interval: time.Hour, Burst: 2},
	// GET+POST /api/v2/list/goods/filter и /upload/task — 1/час (одна категория WB).
	RatePrices:      {Limit: 1, Interval: time.Hour, Burst: 1},
	RatePricesWrite: {Limit: 1, Interval: time.Hour, Burst: 1, PersistAcrossSync: true},
	// statistics-api — 1 запрос в 2 часа (базовый токен).
	// Sales, orders, supplier/stocks — все в одной категории.
	RateStatistics: {Limit: 1, Interval: 2 * time.Hour, Burst: 1},
	// GET feedbacks/questions — 5/час, burst=2 (feedbacks+questions запускаются параллельно, нужно 2 токена за тик).
	RateFeedbacks: {Limit: 5, Interval: time.Hour, Burst: 2},
	// POST answer — 5/час. PersistAcrossSync: AI-джоба работает вне синка.
	RateFeedbacksWrite: {Limit: 5, Interval: time.Hour, Burst: 1, PersistAcrossSync: true},
}

// servicePolicies — лимиты для сервисного/персонального токена WB (подтверждены май 2026).
//
// FBS/DBW/DBS/Marketplace: 300 req/min, интервал 200 мс, burst 20.
// Content:                 100 req/min, интервал 600 мс, burst 5.
// Prices (read+write):      10 req/6sec, интервал 600 мс, burst 5 (одна категория WB).
// Feedbacks (read+write):    3 req/sec,  интервал 333 мс, burst 6.
// Statistics (sales/orders/stocks): 1 req/min, burst 1 — ограничение одинаково для всех типов токенов.
var servicePolicies = map[RateCategory]RatePolicy{
	// 100 req/min, interval 600ms, burst 5
	RateContent: {Limit: 100, Interval: time.Minute, Burst: 5},
	// 300 req/min, interval 200ms, burst 20 (stocks + warehouses)
	RateMarketplace: {Limit: 300, Interval: time.Minute, Burst: 20},
	// 300 req/min, interval 200ms, burst 20
	RateFBSOrders: {Limit: 300, Interval: time.Minute, Burst: 20},
	// 300 req/min, interval 200ms, burst 20
	RateDBOrders: {Limit: 300, Interval: time.Minute, Burst: 20},
	// 10 req/6sec, interval 600ms, burst 5 (одна категория: GET + POST цены)
	RatePrices:      {Limit: 10, Interval: 6 * time.Second, Burst: 5},
	RatePricesWrite: {Limit: 10, Interval: 6 * time.Second, Burst: 5, PersistAcrossSync: true},
	// 1 req/min — лимит одинаков для всех типов токенов. Daily sync: sales+refusals+stocks.
	RateStatistics: {Limit: 1, Interval: time.Minute, Burst: 1},
	// 3 req/sec, interval 333ms, burst 6
	RateFeedbacks: {Limit: 3, Interval: time.Second, Burst: 6},
	// 3 req/sec, burst 6. PersistAcrossSync: AI-джоба работает вне синка.
	RateFeedbacksWrite: {Limit: 3, Interval: time.Second, Burst: 6, PersistAcrossSync: true},
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
