package linktokens

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Consume для бота: токены, созданные на сайте (user_id задан). Возвращает user_id.
func (r *Repository) Consume(ctx context.Context, token string, now time.Time) (userID int, err error) {
	row := struct {
		UserID int `gorm:"column:user_id"`
	}{}
	err = r.DB.WithContext(ctx).
		Table("link_tokens").
		Select("user_id").
		Where("token = ?", token).
		Where("expires_at > ?", now).
		Where("user_id IS NOT NULL").
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("link_tokens consume get: %w", err)
	}
	if err := r.DB.WithContext(ctx).Exec("DELETE FROM link_tokens WHERE token = ?", token).Error; err != nil {
		return 0, fmt.Errorf("link_tokens consume delete: %w", err)
	}
	return row.UserID, nil
}

// ConsumeForWeb для сайта: токены, созданные в боте (telegram_chat_id задан). Возвращает chatID.
func (r *Repository) ConsumeForWeb(ctx context.Context, token string, now time.Time) (telegramChatID int64, err error) {
	row := struct {
		TelegramChatID int64 `gorm:"column:telegram_chat_id"`
	}{}
	err = r.DB.WithContext(ctx).
		Table("link_tokens").
		Select("telegram_chat_id").
		Where("token = ?", token).
		Where("expires_at > ?", now).
		Where("telegram_chat_id IS NOT NULL").
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, fmt.Errorf("link_tokens consume for web get: %w", err)
	}
	if err := r.DB.WithContext(ctx).Exec("DELETE FROM link_tokens WHERE token = ?", token).Error; err != nil {
		return 0, fmt.Errorf("link_tokens consume for web delete: %w", err)
	}
	return row.TelegramChatID, nil
}
