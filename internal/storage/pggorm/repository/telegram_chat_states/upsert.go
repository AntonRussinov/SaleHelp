package telegramchatstates

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

func (r *Repository) Upsert(ctx context.Context, st *domain.TelegramChatState, now time.Time) error {
	payload := st.Payload
	if payload == nil {
		payload = []byte(`{}`)
	}

	m := model.TelegramChatState{
		ChatID:    st.ChatID,
		State:     st.State,
		Payload:   datatypes.JSON(payload),
		UpdatedAt: now,
	}

	// ON CONFLICT (chat_id) DO UPDATE
	if err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "chat_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"state", "payload", "updated_at"}),
		}).
		Create(&m).Error; err != nil {
		return fmt.Errorf("telegram_chat_states upsert: %w", err)
	}

	return nil
}
