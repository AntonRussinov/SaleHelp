package client

import "net/http"

type Client struct {
	httpClient     *http.Client
	baseURL        string
	marketplaceURL string
	priceURL       string
	limiters       *LimiterManager
}

func New(httpClient *http.Client) *Client {
	return &Client{
		httpClient:     httpClient,
		baseURL:        "https://content-api.wildberries.ru",
		marketplaceURL: "https://marketplace-api.wildberries.ru",
		priceURL:       "https://discounts-prices-api.wildberries.ru",
		limiters:       NewLimiterManager(defaultPolicies),
	}
}
