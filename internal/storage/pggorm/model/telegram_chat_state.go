package model

import (
	"time"

	"salehelp/domain"

	"gorm.io/datatypes"
)

type TelegramChatState struct {
	ChatID    int64          `gorm:"column:chat_id;primaryKey"`
	State     string         `gorm:"column:state;not null"`
	Payload   datatypes.JSON `gorm:"column:payload;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
}

func (TelegramChatState) TableName() string { return "telegram_chat_states" }

func (m TelegramChatState) ToDomain() *domain.TelegramChatState {
	return &domain.TelegramChatState{
		ChatID:    m.ChatID,
		State:     m.State,
		Payload:   domain.JSONRaw(m.Payload), // см. helper ниже
		UpdatedAt: m.UpdatedAt,
	}
}
