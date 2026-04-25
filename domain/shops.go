package domain

import "time"

type Shop struct {
	ID                  int
	UserID              int
	Marketplace         string
	ExternalID          string
	Name                string
	APIToken            string
	IsActive            bool
	// WBReportTimezone — часовой пояс для отчётов WB (sales, orders). IANA, например Europe/Moscow, Asia/Almaty. Пусто = из конфига.
	WBReportTimezone    string
	LastSyncAt          *time.Time
	FirstSyncNotifiedAt *time.Time
	LastSyncStatusCode  int
	LastSyncError       string
	LastSyncFailedAt    *time.Time
	LastSyncOkAt        *time.Time
	// Поля для разделения синка на быстрый (30 мин) и тяжёлый (1 час).
	LastHeavySyncAt     *time.Time
	// Курсор для инкрементального синка карточек WB (updatedAt + nmID из последнего ответа).
	CardsCursorUpdatedAt string
	CardsCursorNmID      int
	// Время последнего полного (не инкрементального) синка карточек — раз в 7 дней.
	LastFullCardsSyncAt *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
