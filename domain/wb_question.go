package domain

import "time"

// WBQuestion — вопрос покупателя из WB API (Вопросы и отзывы).
// Используется для уведомлений о новых вопросах в Telegram.
type WBQuestion struct {
	ID          int
	ShopID      int
	WBQuestionID string    // id из WB API (уникален в рамках WB)
	NmID        int64
	ProductName string
	Text        string
	CreatedAt   time.Time
	NotifiedAt  *time.Time
	Extra       []byte // JSON при необходимости
}
