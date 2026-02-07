package domain

import "time"

// WBFeedback — отзыв покупателя из WB API (Вопросы и отзывы).
// Используется для уведомлений о новых отзывах в Telegram.
type WBFeedback struct {
	ID               int
	ShopID           int
	WBFeedbackID     string    // id из WB API (уникален в рамках WB)
	NmID             int64
	ProductName      string
	Text             string
	Pros             string
	Cons             string
	ProductValuation int
	UserName         string
	CreatedAt        time.Time
	NotifiedAt       *time.Time
	Extra            []byte // JSON при необходимости
}
