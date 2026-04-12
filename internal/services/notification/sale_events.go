package notifications

import (
	"fmt"
	"strings"
	"time"

	"salehelp/domain"
)

// FormatSaleEventsMessageRU формирует текст сообщения о новых выкупах/возвратах для отправки в Telegram.
// productTitleByNmID — название товара по nm_id (из карточки WB); если нет — используется артикул или "nm ID".
func FormatSaleEventsMessageRU(events []*domain.SaleEvent, productTitleByNmID map[int64]string, loc *time.Location) string {
	if len(events) == 0 {
		return ""
	}
	if loc == nil {
		loc = time.UTC
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
		parts = append(parts, formatSaleEventsBlockRU("✅ Выкупы", buyouts, productTitleByNmID, loc))
	}
	if len(returns) > 0 {
		parts = append(parts, formatSaleEventsBlockRU("↩️ Возвраты", returns, productTitleByNmID, loc))
	}
	return strings.Join(parts, "\n\n")
}

// FormatOrderRefusalsMessageRU формирует текст блока об отказах (заказы не забраны).
// Группирует по дате отмены (CancelDate), внутри каждого дня — по названию товара со счётчиком.
// Так сообщение остаётся точным, даже если отказы накопились за несколько дней.
func FormatOrderRefusalsMessageRU(refusals []*domain.OrderRefusal, productTitleByNmID map[int64]string, loc *time.Location) string {
	if len(refusals) == 0 {
		return ""
	}
	if loc == nil {
		loc = time.UTC
	}

	type nameEntry struct {
		count int
	}
	type dayEntry struct {
		date      string
		nameOrder []string
		byName    map[string]*nameEntry
	}

	// Группируем по дням, сохраняя порядок первого появления дня и товара.
	dayOrder := make([]string, 0)
	byDay := make(map[string]*dayEntry)

	for _, r := range refusals {
		dateKey := r.CancelDate.In(loc).Format("02.01.2006")
		name := productTitleByNmID[r.NmID]
		if name == "" {
			name = r.SupplierArticle
		}
		if name == "" {
			name = fmt.Sprintf("nm %d", r.NmID)
		}

		day, ok := byDay[dateKey]
		if !ok {
			day = &dayEntry{date: dateKey, byName: make(map[string]*nameEntry)}
			byDay[dateKey] = day
			dayOrder = append(dayOrder, dateKey)
		}
		if _, seen := day.byName[name]; !seen {
			day.nameOrder = append(day.nameOrder, name)
			day.byName[name] = &nameEntry{}
		}
		day.byName[name].count++
	}

	var parts []string
	for _, dateKey := range dayOrder {
		day := byDay[dateKey]
		lines := make([]string, 0, len(day.nameOrder)+1)
		lines = append(lines, day.date+":")
		for _, name := range day.nameOrder {
			e := day.byName[name]
			if e.count > 1 {
				lines = append(lines, fmt.Sprintf("• %s — %d шт.", name, e.count))
			} else {
				lines = append(lines, fmt.Sprintf("• %s", name))
			}
		}
		parts = append(parts, strings.Join(lines, "\n"))
	}

	return "🚫 Отказы (не забрали)\n\n" + strings.Join(parts, "\n\n")
}

// FormatSalesAndRefusalsMessageRU объединяет выкупы, возвраты и отказы в одно сообщение.
func FormatSalesAndRefusalsMessageRU(events []*domain.SaleEvent, refusals []*domain.OrderRefusal, productTitleByNmID map[int64]string, loc *time.Location) string {
	var parts []string
	if len(events) > 0 {
		parts = append(parts, FormatSaleEventsMessageRU(events, productTitleByNmID, loc))
	}
	if len(refusals) > 0 {
		parts = append(parts, FormatOrderRefusalsMessageRU(refusals, productTitleByNmID, loc))
	}
	return strings.Join(parts, "\n\n")
}

func formatSaleEventsBlockRU(header string, events []*domain.SaleEvent, productTitleByNmID map[int64]string, loc *time.Location) string {
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
		date := e.Date.In(loc).Format("02.01.2006 15:04")
		lines = append(lines, fmt.Sprintf("• %s — %s, %s", name, date, formatPriceRU(e.ForPay)))
	}
	return strings.Join(lines, "\n")
}

func formatPriceRU(v float64) string {
	return fmt.Sprintf("%.2f ₽", v)
}
