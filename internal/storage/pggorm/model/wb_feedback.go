package model

import (
	"salehelp/domain"
	"time"
)

type WBFeedback struct {
	ID               int        `gorm:"primaryKey;column:id"`
	ShopID           int        `gorm:"column:shop_id;not null"`
	WBFeedbackID     string     `gorm:"column:wb_feedback_id;type:text;not null"`
	NmID             int64      `gorm:"column:nm_id;not null"`
	ProductName      string     `gorm:"column:product_name;type:text;not null"`
	Text             string     `gorm:"column:text;type:text;not null"`
	Pros             string     `gorm:"column:pros;type:text;not null"`
	Cons             string     `gorm:"column:cons;type:text;not null"`
	ProductValuation int        `gorm:"column:product_valuation;not null"`
	UserName         string     `gorm:"column:user_name;type:text;not null"`
	CreatedAt        time.Time  `gorm:"column:created_at;not null"`
	NotifiedAt       *time.Time `gorm:"column:notified_at"`
	IsReplied        bool       `gorm:"column:is_replied;not null;default:false"`
	RepliedAt        *time.Time `gorm:"column:replied_at"`
	OrderStatus      string     `gorm:"column:order_status;type:text;not null;default:''"`
	Extra            []byte     `gorm:"column:extra;type:jsonb"`
}

func (m *WBFeedback) ToDomain() *domain.WBFeedback {
	if m == nil {
		return nil
	}
	return &domain.WBFeedback{
		ID:               m.ID,
		ShopID:           m.ShopID,
		WBFeedbackID:     m.WBFeedbackID,
		NmID:             m.NmID,
		ProductName:      m.ProductName,
		Text:             m.Text,
		Pros:             m.Pros,
		Cons:             m.Cons,
		ProductValuation: m.ProductValuation,
		UserName:         m.UserName,
		CreatedAt:        m.CreatedAt,
		NotifiedAt:       m.NotifiedAt,
		IsReplied:        m.IsReplied,
		RepliedAt:        m.RepliedAt,
		OrderStatus:      m.OrderStatus,
		Extra:            m.Extra,
	}
}

func (*WBFeedback) TableName() string {
	return "wb_feedbacks"
}
