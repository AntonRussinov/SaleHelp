package notifications

import (
	"fmt"
	"strings"
	"time"
)

// AssemblyOrderRow — минимальные данные заказа для формирования сообщения.
// Используется для FBS, DBW, DBS.
type AssemblyOrderRow struct {
	NmID         int64
	Article      string
	WBOrderID    int64
	CreatedAt    time.Time
	StockWarning bool // true если по этому nmID активных заказов >= остаток на складе
}

// FormatAssemblyOrdersMessageRU формирует текст сообщения о новых сборочных заданиях для Telegram.
// orderType — "FBS", "DBW" или "DBS" (тип задания на сборку).
// productTitleByNmID — название товара по nm_id (из карточки WB); если нет — артикул или "nm ID".
// Одинаковые товары группируются в одну строку с указанием количества.
func FormatAssemblyOrdersMessageRU(orderType string, orders []AssemblyOrderRow, productTitleByNmID map[int64]string, _ *time.Location) string {
	if len(orders) == 0 {
		return ""
	}

	type group struct {
		name         string
		count        int
		stockWarning bool
	}

	seen := make(map[int64]int) // nmID → index in groups
	var groups []group

	for _, o := range orders {
		name := productTitleByNmID[o.NmID]
		if name == "" {
			name = o.Article
		}
		if name == "" {
			name = fmt.Sprintf("nm %d", o.NmID)
		}
		if idx, ok := seen[o.NmID]; ok {
			groups[idx].count++
			if o.StockWarning {
				groups[idx].stockWarning = true
			}
		} else {
			seen[o.NmID] = len(groups)
			groups = append(groups, group{name: name, count: 1, stockWarning: o.StockWarning})
		}
	}

	header := fmt.Sprintf("📦 Новые сборочные задания %s", orderType)
	lines := make([]string, 0, len(groups))
	for _, g := range groups {
		suffix := ""
		if g.stockWarning {
			suffix = " ⚠️ последний товар!"
		}
		lines = append(lines, fmt.Sprintf("• %s — %d шт.%s", g.name, g.count, suffix))
	}
	return header + "\n\n" + strings.Join(lines, "\n")
}
