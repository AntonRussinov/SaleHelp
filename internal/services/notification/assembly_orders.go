package notifications

import (
	"fmt"
	"strings"
	"time"
)

// AssemblyOrderRow — минимальные данные заказа для формирования сообщения.
// Используется для FBS, DBW, DBS.
type AssemblyOrderRow struct {
	NmID      int64
	Article   string
	WBOrderID int64
	CreatedAt time.Time
}

// FormatAssemblyOrdersMessageRU формирует текст сообщения о новых сборочных заданиях для Telegram.
// orderType — "FBS", "DBW" или "DBS" (тип задания на сборку).
// productTitleByNmID — название товара по nm_id (из карточки WB); если нет — артикул или "nm ID".
func FormatAssemblyOrdersMessageRU(orderType string, orders []AssemblyOrderRow, productTitleByNmID map[int64]string, loc *time.Location) string {
	if len(orders) == 0 {
		return ""
	}
	if loc == nil {
		loc = time.UTC
	}
	header := fmt.Sprintf("📦 Новые сборочные задания %s", orderType)
	lines := make([]string, 0, len(orders)+2)
	for _, o := range orders {
		name := productTitleByNmID[o.NmID]
		if name == "" {
			name = o.Article
		}
		if name == "" {
			name = fmt.Sprintf("nm %d", o.NmID)
		}
		created := o.CreatedAt.In(loc).Format("02.01.2006 15:04")
		lines = append(lines, fmt.Sprintf("• %s — заказ #%d, %s", name, o.WBOrderID, created))
	}
	return header + "\n\n" + strings.Join(lines, "\n")
}
