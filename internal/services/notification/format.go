package notifications

import (
	"fmt"
	"sort"
	"strings"

	"salehelp/domain"
	"salehelp/pkg/notify"
)

func FormatAggregatedMessageRU(
	alertType string,
	p *domain.Product,
	triggered []*domain.AlertTriggered,
	sizesByID map[int]*domain.ProductSize,
) notify.Message {

	title := "SaleHelp"

	prodTitle := "Товар"
	prodArt := "—"
	prodURL := ""

	if p != nil {
		if p.Title != "" {
			prodTitle = p.Title
		}
		if p.ExternalID != "" {
			prodArt = p.ExternalID
		}
		if p.URL != "" {
			prodURL = p.URL
		}
	}

	// Стабильный порядок: sizeName, потом без size в конце
	sort.SliceStable(triggered, func(i, j int) bool {
		a := sizeName(triggered[i], sizesByID)
		b := sizeName(triggered[j], sizesByID)
		return a < b
	})

	lines := make([]string, 0, len(triggered))
	for _, tr := range triggered {
		sn := sizeName(tr, sizesByID)
		if sn != "" {
			lines = append(lines, fmt.Sprintf("• %s — %s", sn, tr.Message))
		} else {
			lines = append(lines, fmt.Sprintf("• %s", tr.Message))
		}
	}

	header := fmt.Sprintf("🛎 %s\n\n%s (арт. %s)", humanAlertTypeRU(alertType), prodTitle, prodArt)
	if prodURL != "" {
		header += "\n" + prodURL
	}

	body := header + "\n\n" + strings.Join(lines, "\n")

	return notify.Message{
		Title: title,
		Text:  body,
	}
}

func sizeName(tr *domain.AlertTriggered, sizesByID map[int]*domain.ProductSize) string {
	if tr == nil || tr.ProductSizeID == nil {
		return ""
	}
	sz := sizesByID[*tr.ProductSizeID]
	if sz == nil {
		return fmt.Sprintf("size_id=%d", *tr.ProductSizeID)
	}
	if sz.TechSizeName != nil && *sz.TechSizeName != "" {
		return *sz.TechSizeName
	}
	// fallback: external_id (chrtID)
	if sz.ExternalID != "" {
		return sz.ExternalID
	}
	return fmt.Sprintf("size_id=%d", *tr.ProductSizeID)
}

func humanAlertTypeRU(t string) string {
	switch t {
	case "stock_low":
		return "Низкий остаток"
	case "stock_zero":
		return "Размер закончился"
	case "price_drop":
		return "Цена снизилась"
	default:
		return "Алерт"
	}
}
