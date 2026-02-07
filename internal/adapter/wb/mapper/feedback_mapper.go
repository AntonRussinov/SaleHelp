package mapper

import (
	"time"

	"salehelp/domain"
	"salehelp/internal/adapter/wb/dto"
)

// WBFeedbackToDomain конвертирует DTO отзыва WB в домен.
func WBFeedbackToDomain(f dto.WBFeedback, shopID int) domain.WBFeedback {
	createdAt := parseWBDateTime(f.CreatedDate)
	return domain.WBFeedback{
		ShopID:           shopID,
		WBFeedbackID:     f.ID,
		NmID:             f.ProductDetails.NmID,
		ProductName:      f.ProductDetails.ProductName,
		Text:             f.Text,
		Pros:             f.Pros,
		Cons:             f.Cons,
		ProductValuation: f.ProductValuation,
		UserName:         f.UserName,
		CreatedAt:        createdAt,
		NotifiedAt:       nil,
	}
}

// WBQuestionToDomain конвертирует DTO вопроса WB в домен.
func WBQuestionToDomain(q dto.WBQuestion, shopID int) domain.WBQuestion {
	createdAt := parseWBDateTime(q.CreatedDate)
	return domain.WBQuestion{
		ShopID:       shopID,
		WBQuestionID: q.ID,
		NmID:         q.ProductDetails.NmID,
		ProductName:  q.ProductDetails.ProductName,
		Text:         q.Text,
		CreatedAt:    createdAt,
		NotifiedAt:   nil,
	}
}

// parseWBDateTime парсит дату из WB API (RFC3339 или с наносекундами).
func parseWBDateTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999Z07:00",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
