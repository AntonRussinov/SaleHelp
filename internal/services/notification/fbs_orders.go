package notifications

import (
	"fmt"
	"strings"

	"salehelp/domain"
)

// FormatFBSOrdersMessageRU формирует текст сообщения о новых FBS-заказах для отправки в Telegram.
// productTitleByNmID — название товара по nm_id (из карточки WB); если нет — используется артикул или "nm ID".
func FormatFBSOrdersMessageRU(orders []*domain.FBSOrder, productTitleByNmID map[int64]string) string {
	if len(orders) == 0 {
		return ""
	}
	header := "📦 Новые сборочные задания FBS"
	lines := make([]string, 0, len(orders)+2)
	for _, o := range orders {
		name := productTitleByNmID[o.NmID]
		if name == "" {
			name = o.Article
		}
		if name == "" {
			name = fmt.Sprintf("nm %d", o.NmID)
		}
		created := o.CreatedAt.Format("02.01.2006 15:04")
		lines = append(lines, fmt.Sprintf("• %s — заказ #%d, %s", name, o.WBOrderID, created))
	}
	return header + "\n\n" + strings.Join(lines, "\n")
}
