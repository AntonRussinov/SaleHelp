package notifications

import (
	"fmt"
	"strings"

	"salehelp/domain"
)

// FormatSaleEventsMessageRU формирует текст сообщения о новых выкупах/возвратах для отправки в Telegram.
// productTitleByNmID — название товара по nm_id (из карточки WB); если нет — используется артикул или "nm ID".
func FormatSaleEventsMessageRU(events []*domain.SaleEvent, productTitleByNmID map[int64]string) string {
	if len(events) == 0 {
		return ""
	}
	buyouts := make([]*domain.SaleEvent, 0)
	returns := make([]*domain.SaleEvent, 0)
	for _, e := range events {
		switch e.EventType {
		case domain.SaleEventTypeReturn, domain.SaleEventTypePartialReturn:
			returns = append(returns, e)
		default:
			buyouts = append(buyouts, e)
		}
	}
	var parts []string
	if len(buyouts) > 0 {
		parts = append(parts, formatSaleEventsBlockRU("✅ Выкупы", buyouts, productTitleByNmID))
	}
	if len(returns) > 0 {
		parts = append(parts, formatSaleEventsBlockRU("↩️ Возвраты", returns, productTitleByNmID))
	}
	return strings.Join(parts, "\n\n")
}

func formatSaleEventsBlockRU(header string, events []*domain.SaleEvent, productTitleByNmID map[int64]string) string {
	lines := make([]string, 0, len(events)+1)
	lines = append(lines, header)
	for _, e := range events {
		name := productTitleByNmID[e.NmID]
		if name == "" {
			name = e.SupplierArticle
		}
		if name == "" {
			name = fmt.Sprintf("nm %d", e.NmID)
		}
		date := e.Date.Format("02.01.2006 15:04")
		lines = append(lines, fmt.Sprintf("• %s — %s, %s", name, date, formatPriceRU(e.ForPay)))
	}
	return strings.Join(lines, "\n")
}

func formatPriceRU(v float64) string {
	return fmt.Sprintf("%.2f ₽", v)
}
