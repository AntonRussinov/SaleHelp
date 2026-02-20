package notifications

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"salehelp/domain"
	"salehelp/pkg/notify"
)

const digestChunkLimit = 3500 // запас под HTML escape + <b>Title</b>

func FormatShopDigestMessagesRU(
	triggered []*domain.AlertTriggered,
	productsByID map[int]*domain.Product,
	sizesByID map[int]*domain.ProductSize,
	sizesCountByProduct map[int]int,
) []notify.Message {
	if len(triggered) == 0 {
		return nil
	}

	// order of sections
	typeOrder := []string{"stock_zero", "stock_low", "price_drop"}

	// group by alert type
	byType := map[string][]*domain.AlertTriggered{}
	for _, tr := range triggered {
		if tr == nil || tr.AlertType == "" {
			continue
		}
		byType[tr.AlertType] = append(byType[tr.AlertType], tr)
	}

	header := "📦 Отчет"

	// build all lines (then chunk)
	lines := make([]string, 0, len(triggered)+16)

	// deterministic: only types that exist, in priority order
	for _, at := range typeOrder {
		trs := byType[at]
		if len(trs) == 0 {
			continue
		}

		// section title
		lines = append(lines, fmt.Sprintf("%s:", sectionTitleRU(at)))

		// group by product
		byProduct := map[int][]*domain.AlertTriggered{}
		prodIDs := make([]int, 0, 32)
		for _, tr := range trs {
			byProduct[tr.ProductID] = append(byProduct[tr.ProductID], tr)
		}
		for pid := range byProduct {
			prodIDs = append(prodIDs, pid)
		}
		sort.Ints(prodIDs)

		for _, pid := range prodIDs {
			p := productsByID[pid]
			prodLine := formatProductHeaderRU(p, pid)

			items := byProduct[pid]

			// разделим на size-triggered и product-triggered
			sizeItems := make([]*domain.AlertTriggered, 0, len(items))
			prodItems := make([]*domain.AlertTriggered, 0, 1)

			for _, tr := range items {
				if tr != nil && tr.ProductSizeID != nil {
					sizeItems = append(sizeItems, tr)
				} else {
					prodItems = append(prodItems, tr)
				}
			}

			// 1) product-level события (без размеров) — печатаем как есть
			for _, tr := range prodItems {
				if tr == nil {
					continue
				}
				lines = append(lines, fmt.Sprintf("• %s — %s", prodLine, tr.Message))
			}

			// 2) size-level события — компактим
			if len(sizeItems) > 0 {
				// для стабильности
				sort.SliceStable(sizeItems, func(i, j int) bool {
					ai := displaySizeName(sizeItems[i], sizesByID, sizesCountByProduct)
					aj := displaySizeName(sizeItems[j], sizesByID, sizesCountByProduct)
					return ai < aj
				})

				if at == "stock_zero" {
					// максимальная компрессия: только список размеров
					szs := make([]string, 0, len(sizeItems))
					for _, tr := range sizeItems {
						sn := displaySizeName(tr, sizesByID, sizesCountByProduct)
						if sn != "" {
							szs = append(szs, sn)
						}
					}
					szs = uniqueSorted(szs)

					if len(szs) > 0 {
						lines = append(lines, fmt.Sprintf("• %s — %s", prodLine, formatSizeNamesLine(szs)))
					} else {
						// fallback
						lines = append(lines, fmt.Sprintf("• %s — %s", prodLine, sizeItems[0].Message))
					}
				} else {
					// stock_low (и любые другие size-level) — группируем по одинаковому сообщению
					byMsg := map[string][]string{}
					msgKeys := make([]string, 0, 4)

					for _, tr := range sizeItems {
						msg := ""
						if tr != nil {
							msg = strings.TrimSpace(tr.Message)
						}
						sn := displaySizeName(tr, sizesByID, sizesCountByProduct)
						byMsg[msg] = append(byMsg[msg], sn)
					}

					for msg := range byMsg {
						msgKeys = append(msgKeys, msg)
					}
					sort.Strings(msgKeys)

					if len(msgKeys) == 1 {
						msg := msgKeys[0]
						szs := uniqueSorted(byMsg[msg])
						if len(szs) > 0 && msg != "" {
							lines = append(lines, fmt.Sprintf("• %s — %s — %s", prodLine, formatSizeNamesLine(szs), msg))
						} else if len(szs) > 0 {
							lines = append(lines, fmt.Sprintf("• %s — %s", prodLine, formatSizeNamesLine(szs)))
						} else {
							lines = append(lines, fmt.Sprintf("• %s — %s", prodLine, msg))
						}
					} else {
						lines = append(lines, fmt.Sprintf("• %s", prodLine))
						for _, msg := range msgKeys {
							szs := uniqueSorted(byMsg[msg])
							if len(szs) > 0 && msg != "" {
								lines = append(lines, fmt.Sprintf("  – %s — %s", formatSizeNamesLine(szs), msg))
							} else if len(szs) > 0 {
								lines = append(lines, fmt.Sprintf("  – %s", formatSizeNamesLine(szs)))
							} else {
								lines = append(lines, fmt.Sprintf("  – %s", msg))
							}
						}
					}
				}
			}

		}

		lines = append(lines, "") // blank line between sections
	}

	// trim trailing empty lines
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	// chunk by rune length (safe for Telegram/WhatsApp; Email тоже ок)
	chunks := chunkLines(header, lines, digestChunkLimit)

	msgs := make([]notify.Message, 0, len(chunks))
	for i, txt := range chunks {
		// можно добавить " (1/3)" если частей много
		title := "SaleHelp"
		if len(chunks) > 1 {
			title = fmt.Sprintf("SaleHelp (%d/%d)", i+1, len(chunks))
		}
		msgs = append(msgs, notify.Message{
			Title: title,
			Text:  txt,
		})
	}
	return msgs
}

// formatSizeNamesLine форматирует список названий размеров для отчёта: "размер M, размер L".
func formatSizeNamesLine(szs []string) string {
	if len(szs) == 0 {
		return ""
	}
	return "размер " + strings.Join(szs, ", размер ")
}

func sectionTitleRU(alertType string) string {
	switch alertType {
	case "stock_zero":
		return "❌ " + humanAlertTypeRU(alertType)
	case "stock_low":
		return "⚠️ " + humanAlertTypeRU(alertType)
	case "price_drop":
		return "🔻 " + humanAlertTypeRU(alertType)
	default:
		return "🛎 " + humanAlertTypeRU(alertType)
	}
}

func formatProductHeaderRU(p *domain.Product, productID int) string {
	title := "Товар"
	art := ""
	if p != nil {
		if strings.TrimSpace(p.Title) != "" {
			title = p.Title
		}
		if strings.TrimSpace(p.ExternalID) != "" {
			art = p.ExternalID
		}
	}
	if art != "" {
		return fmt.Sprintf("%s (арт. %s)", title, art)
	}
	// fallback, но без мусора
	return title
}

func displaySizeName(
	tr *domain.AlertTriggered,
	sizesByID map[int]*domain.ProductSize,
	sizesCountByProduct map[int]int,
) string {
	if tr == nil || tr.ProductSizeID == nil {
		return ""
	}
	sz := sizesByID[*tr.ProductSizeID]
	if sz == nil {
		return ""
	}

	name := ""
	if sz.TechSizeName != nil && strings.TrimSpace(*sz.TechSizeName) != "" {
		name = strings.TrimSpace(*sz.TechSizeName)
	} else if strings.TrimSpace(sz.ExternalID) != "" {
		name = strings.TrimSpace(sz.ExternalID)
	}

	// Не показывать «0» как размер — выглядит как остаток и путает; пустое имя тоже скрываем.
	if isPhantomSizeName(name) || strings.TrimSpace(name) == "" {
		return ""
	}
	return name
}

func isPhantomSizeName(name string) bool {
	n := strings.TrimSpace(strings.ToLower(name))
	return n == "" || n == "0"
}

func chunkLines(header string, lines []string, limit int) []string {
	if limit <= 0 {
		return []string{header + "\n\n" + strings.Join(lines, "\n")}
	}

	out := make([]string, 0, 4)
	cur := header + "\n\n"
	curLen := runeLen(cur)

	for _, ln := range lines {
		add := ln
		if cur != header+"\n\n" {
			add = "\n" + ln
		}

		if curLen+runeLen(add) > limit {
			// flush current chunk
			out = append(out, cur)
			// start new chunk with header again (короткий)
			cur = header + "\n\n" + ln
			curLen = runeLen(cur)
			continue
		}

		cur += add
		curLen += runeLen(add)
	}

	if strings.TrimSpace(cur) != "" {
		out = append(out, cur)
	}
	return out
}

func runeLen(s string) int {
	// чуть быстрее/точнее чем []rune для длинных строк
	return utf8.RuneCountInString(s)
}
