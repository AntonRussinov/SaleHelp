package model

import (
	"salehelp/domain"
	"time"
)

type WBQuestion struct {
	ID           int        `gorm:"primaryKey;column:id"`
	ShopID       int        `gorm:"column:shop_id;not null"`
	WBQuestionID string    `gorm:"column:wb_question_id;type:text;not null"`
	NmID         int64      `gorm:"column:nm_id;not null"`
	ProductName  string     `gorm:"column:product_name;type:text;not null"`
	Text         string     `gorm:"column:text;type:text;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null"`
	NotifiedAt   *time.Time `gorm:"column:notified_at"`
	Extra        []byte     `gorm:"column:extra;type:jsonb"`
}

func (m *WBQuestion) ToDomain() *domain.WBQuestion {
	if m == nil {
		return nil
	}
	return &domain.WBQuestion{
		ID:           m.ID,
		ShopID:       m.ShopID,
		WBQuestionID: m.WBQuestionID,
		NmID:         m.NmID,
		ProductName:  m.ProductName,
		Text:         m.Text,
		CreatedAt:    m.CreatedAt,
		NotifiedAt:   m.NotifiedAt,
		Extra:        m.Extra,
	}
}

func (*WBQuestion) TableName() string {
	return "wb_questions"
}
