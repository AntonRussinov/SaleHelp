package domain

import "time"

type AIRecommendation struct {
	ID        int
	ShopID    int
	ProductID *int
	AlertID   *int
	Type      string
	Model     string
	Prompt    string
	Content   string
	Meta      []byte // raw JSON, можно потом декодить в map[string]any
	CreatedAt time.Time
}
