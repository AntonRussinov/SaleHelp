package notifications

import (
	"sort"
	"strings"
)

func uniqueSorted(ss []string) []string {
	m := map[string]struct{}{}
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := m[s]; ok {
			continue
		}
		m[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func humanAlertTypeRU(t string) string {
	switch t {
	case "stock_low":
		return "Низкий остаток"
	case "stock_zero":
		return "Товар закончился"
	case "price_drop":
		return "Цена снизилась"
	default:
		return "Алерт"
	}
}
