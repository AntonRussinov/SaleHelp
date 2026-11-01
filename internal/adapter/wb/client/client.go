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
	xClientSecret   string         // X-Client-Secret; пустая строка — заголовок не выставляется
}

// SetTokenType регистрирует тип токена для apiKey, выбирая соответствующий набор rate-limit политик.
// Вызывать перед первым использованием ключа (например, в начале синка магазина).
func (c *Client) SetTokenType(apiKey, tokenType string) {
	c.limiters.RegisterKey(apiKey, tokenType)
}

// New создаёт WB-клиент. reportLocation — часовой пояс для параметра dateFrom в отчётах (sales, orders);
// если nil, используется Europe/Moscow. xClientSecret — сервисный секрет WB (X-Client-Secret);
// если пустая строка, заголовок не выставляется.
func New(httpClient *http.Client, reportLocation *time.Location, xClientSecret string) *Client {
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
		limiters:       NewLimiterManager(),
		reportLocation: reportLocation,
		xClientSecret:  xClientSecret,
	}
}
