package linktokens

import (
	"context"
	"crypto/rand"
	"fmt"
	"salehelp/domain"
	"time"
)

const (
	fromTelegramTokenTTL    = 10 * time.Minute
	fromTelegramTokenLength = 6
	fromTelegramCharset     = "123456789ABCDEFGHJKLMNPQRSTUVWXYZ"
)

// CreateFromTelegram создаёт одноразовый код для привязки «из Telegram»: пользователь вводит email в боте,
// получает код, затем заходит на сайт под этим email и вводит код там (подтверждение владельца email).
func (r *Repository) CreateFromTelegram(ctx context.Context, chatID int64) (code string, expiresAt time.Time, err error) {
	code, err = generateCode(fromTelegramTokenLength)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("generate code: %w", err)
	}
	expiresAt = time.Now().Add(fromTelegramTokenTTL)
	t := &domain.LinkToken{
		Token:          code,
		TelegramChatID: &chatID,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now(),
	}
	if err := r.Create(ctx, t); err != nil {
		return "", time.Time{}, err
	}
	return code, expiresAt, nil
}

func generateCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = fromTelegramCharset[int(b[i])%len(fromTelegramCharset)]
	}
	return string(b), nil
}
