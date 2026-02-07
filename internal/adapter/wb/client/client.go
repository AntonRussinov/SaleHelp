package client

import (
	"net/http"
	"time"
)

type Client struct {
	httpClient      *http.Client
	baseURL         string
	marketplaceURL  string
	priceURL        string
	statisticsURL   string
	feedbacksURL    string
	limiters        *LimiterManager
	reportLocation  *time.Location // для dateFrom в отчётах (sales, orders) — в каком поясе отправлять
}

// New создаёт WB-клиент. reportLocation — часовой пояс для параметра dateFrom в отчётах (sales, orders);
// если nil, используется Europe/Moscow.
func New(httpClient *http.Client, reportLocation *time.Location) *Client {
	if reportLocation == nil {
		reportLocation, _ = time.LoadLocation("Europe/Moscow")
		if reportLocation == nil {
			reportLocation = time.FixedZone("MSK", 3*3600)
		}
	}
	return &Client{
		httpClient:     httpClient,
		baseURL:        "https://content-api.wildberries.ru",
		marketplaceURL: "https://marketplace-api.wildberries.ru",
		priceURL:       "https://discounts-prices-api.wildberries.ru",
		statisticsURL:  "https://statistics-api.wildberries.ru",
		feedbacksURL:   "https://feedbacks-api.wildberries.ru",
		limiters:       NewLimiterManager(defaultPolicies),
		reportLocation: reportLocation,
	}
}
