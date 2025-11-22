package domain

import "time"

type Alert struct {
	ID                 int
	ShopID             int
	ProductID          *int
	Type               string
	Severity           int
	Status             string
	Message            string
	Payload            []byte
	TriggeredAt        time.Time
	ResolvedAt         *time.Time
	AIRecommendationID *int
}
