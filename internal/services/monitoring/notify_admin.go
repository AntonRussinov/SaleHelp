package monitoring

import (
	"context"
	"fmt"
	"time"

	"salehelp/domain"
)

const maxErrLen = 300

func (s *Service) notifyAdminSyncFailed(shop *domain.Shop, err error) {
	if s.adminChatID == 0 || s.tg == nil {
		return
	}

	errText := err.Error()
	if len(errText) > maxErrLen {
		errText = errText[:maxErrLen] + "..."
	}

	msg := fmt.Sprintf(
		"⚠️ Синк упал: «%s» (ID: %d)\n\n%s",
		shop.Name, shop.ID, errText,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = s.tg.SendMessage(ctx, s.adminChatID, msg)
}
