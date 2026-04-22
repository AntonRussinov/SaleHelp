package client

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	RateContent     RateCategory = "content"     // content-api.wildberries.ru
	RateMarketplace RateCategory = "marketplace" // marketplace-api.wildberries.ru
	RatePrices      RateCategory = "prices"      // discounts-prices-api.wildberries.ru
	RateStatistics  RateCategory = "statistics"  // statistics-api.wildberries.ru (sales, orders, stocks)
	RateFeedbacks   RateCategory = "feedbacks"   // feedbacks-api.wildberries.ru (вопросы и отзывы: 3 req/1 sec, burst 6)
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

var defaultPolicies = map[RateCategory]RatePolicy{
	RateContent: {
		Limit:    100,
		Interval: time.Minute,
		Burst:    5,
	},
	RateMarketplace: {
		Limit:    300,
		Interval: time.Minute,
		Burst:    20,
	},
	RatePrices: {
		Limit:    20,
		Interval: 6 * time.Second,
		Burst:    20,
	},
	RateStatistics: {
		Limit:    1,
		Interval: time.Minute,
		Burst:    3,
	},
	RateFeedbacks: {
		Limit:    3,
		Interval: time.Second,
		Burst:    6,
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
