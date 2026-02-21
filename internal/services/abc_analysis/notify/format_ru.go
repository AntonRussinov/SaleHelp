package notify

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"salehelp/domain"
)

func formatMoney(v float64) string {
	if v >= 1e6 {
		return fmt.Sprintf("%.0fM ₽", v/1e6)
	}
	if v >= 1e3 {
		return fmt.Sprintf("%.0fk ₽", v/1e3)
	}
	return fmt.Sprintf("%.0f ₽", v)
}

func countByCategory(rows []domain.ABCSnapshotRow, cat domain.ABCCategory) (n int, revenue float64) {
	for _, r := range rows {
		if r.Category == cat {
			n++
			revenue += r.Revenue
		}
	}
	return n, revenue
}

// FormatWeeklyDigestRU форматирует краткий еженедельный дайджест ABC.
func FormatWeeklyDigestRU(snap *domain.ABCSnapshot, changes []domain.ABCChange) string {
	if snap == nil {
		return "📊 ABC-анализ за неделю\nНет данных."
	}
	var b strings.Builder
	b.WriteString("📊 ABC-анализ за неделю\n\n")

	total := snap.TotalRevenue
	if total == 0 {
		total = 1
	}
	na, ra := countByCategory(snap.Rows, domain.ABCCategoryA)
	nb, rb := countByCategory(snap.Rows, domain.ABCCategoryB)
	nc, rc := countByCategory(snap.Rows, domain.ABCCategoryC)
	pa := ra / total * 100
	pb := rb / total * 100
	pc := rc / total * 100

	b.WriteString(fmt.Sprintf("💎 A-товары: %d шт (%s, %.0f%%)\n", na, formatMoney(ra), pa))
	b.WriteString(fmt.Sprintf("⚖️ B-товары: %d шт (%s, %.0f%%)\n", nb, formatMoney(rb), pb))
	b.WriteString(fmt.Sprintf("📦 C-товары: %d шт (%s, %.0f%%)\n\n", nc, formatMoney(rc), pc))

	if len(changes) > 0 {
		b.WriteString("🔥 Изменения:\n")
		for _, c := range changes {
			subj := c.Subject
			if subj == "" {
				subj = fmt.Sprintf("nm %d", c.NmID)
			}
			arrow := "→"
			if c.ChangePct > 0 {
				arrow = "↗️"
			} else if c.ChangePct < 0 {
				arrow = "↘️"
			}
			b.WriteString(fmt.Sprintf("%s \"%s\": %s → %s (%+.0f%%)\n", arrow, truncate(subj, 40), c.FromCat, c.ToCat, c.ChangePct))
		}
		b.WriteString("\n")
	}

	// C-товары без продаж >30 дней — считаем по нулевой выручке в снэпшоте за 7 дней (упрощение: просто кол-во C с 0 продаж)
	noSalesC := 0
	for _, r := range snap.Rows {
		if r.Category == domain.ABCCategoryC && r.SalesCount == 0 {
			noSalesC++
		}
	}
	if noSalesC > 0 {
		b.WriteString(fmt.Sprintf("⚠️ %d C-товаров не продаются >30 дней\n", noSalesC))
	}

	return strings.TrimSpace(b.String())
}

// FormatMonthlyReportRU форматирует детальный месячный отчёт ABC.
func FormatMonthlyReportRU(snap *domain.ABCSnapshot, prev *domain.ABCSnapshot, productByNm map[int64]*domain.Product) string {
	if snap == nil {
		return "📈 ABC-анализ за месяц\nНет данных."
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("📈 ABC-анализ за %s\n\n", formatMonthYearRU(snap.PeriodEnd)))

	total := snap.TotalRevenue
	if total == 0 {
		total = 1
	}
	na, ra := countByCategory(snap.Rows, domain.ABCCategoryA)
	nb, rb := countByCategory(snap.Rows, domain.ABCCategoryB)
	nc, rc := countByCategory(snap.Rows, domain.ABCCategoryC)
	pa := ra / total * 100
	pb := rb / total * 100
	pc := rc / total * 100

	b.WriteString(fmt.Sprintf("💎 A-товары (%d шт): %s (%.0f%%)\n", na, formatMoney(ra), pa))
	topA := topNByCategory(snap.Rows, domain.ABCCategoryA, 3)
	if len(topA) > 0 {
		b.WriteString("  Топ-3:\n")
		for i, row := range topA {
			title := productTitle(productByNm, row.NmID)
			b.WriteString(fmt.Sprintf("  %d. %s — %s\n", i+1, title, formatMoney(row.Revenue)))
		}
	}
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("⚖️ B-товары (%d шт): %s (%.0f%%)\n", nb, formatMoney(rb), pb))
	topB := topNByCategory(snap.Rows, domain.ABCCategoryB, 3)
	if len(topB) > 0 {
		b.WriteString("  Топ-3 по выручке:\n")
		for i, row := range topB {
			title := productTitle(productByNm, row.NmID)
			b.WriteString(fmt.Sprintf("  %d. %s — %s\n", i+1, title, formatMoney(row.Revenue)))
		}
	}
	if prev != nil {
		candidates := bCandidatesToA(snap, prev, productByNm, 3)
		if len(candidates) > 0 {
			b.WriteString("  🔥 Кандидаты в A (рост к прошлому периоду):\n")
			for _, c := range candidates {
				b.WriteString(fmt.Sprintf("  • %s (%+.0f%%)\n", c.title, c.changePct))
			}
		}
	}
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("📦 C-товары (%d шт): %s (%.0f%%)\n", nc, formatMoney(rc), pc))
	noSalesC := 0
	for _, r := range snap.Rows {
		if r.Category == domain.ABCCategoryC && r.SalesCount == 0 {
			noSalesC++
		}
	}
	if noSalesC > 0 {
		b.WriteString(fmt.Sprintf("  ⚠️ Застой: %d товаров не продаются >45 дней\n", noSalesC))
	}
	b.WriteString("\n")

	if prev != nil {
		paPrev, _ := countByCategory(prev.Rows, domain.ABCCategoryA)
		pbPrev, _ := countByCategory(prev.Rows, domain.ABCCategoryB)
		pcPrev, _ := countByCategory(prev.Rows, domain.ABCCategoryC)
		b.WriteString("📊 Динамика категорий:\n")
		b.WriteString(fmt.Sprintf("  A: %d шт (было %d) %s\n", na, paPrev, trend(na, paPrev)))
		b.WriteString(fmt.Sprintf("  B: %d шт (было %d) %s\n", nb, pbPrev, trend(nb, pbPrev)))
		b.WriteString(fmt.Sprintf("  C: %d шт (было %d) %s\n", nc, pcPrev, trend(nc, pcPrev)))
	} else {
		b.WriteString("📊 Динамика категорий появится после следующего отчёта (сравнение с предыдущим периодом).\n")
	}

	return strings.TrimSpace(b.String())
}

// bCandidateToA — B-товар с приростом выручки к прошлому периоду (кандидат в A).
type bCandidateToA struct {
	title    string
	changePct float64
}

// bCandidatesToA возвращает B-товары из текущего снэпшота с наибольшим % роста выручки к prev (топ N).
func bCandidatesToA(snap, prev *domain.ABCSnapshot, productByNm map[int64]*domain.Product, n int) []bCandidateToA {
	prevByNm := make(map[int64]float64)
	for _, r := range prev.Rows {
		prevByNm[r.NmID] = r.Revenue
	}
	var withGrowth []struct {
		nmID     int64
		changePct float64
	}
	for _, r := range snap.Rows {
		if r.Category != domain.ABCCategoryB {
			continue
		}
		prevRev := prevByNm[r.NmID]
		if prevRev <= 0 {
			continue
		}
		changePct := (r.Revenue - prevRev) / prevRev * 100
		if changePct <= 0 {
			continue
		}
		withGrowth = append(withGrowth, struct {
			nmID     int64
			changePct float64
		}{r.NmID, changePct})
	}
	sort.Slice(withGrowth, func(i, j int) bool { return withGrowth[i].changePct > withGrowth[j].changePct })
	if len(withGrowth) > n {
		withGrowth = withGrowth[:n]
	}
	out := make([]bCandidateToA, 0, len(withGrowth))
	for _, w := range withGrowth {
		title := productTitle(productByNm, w.nmID)
		out = append(out, bCandidateToA{title: title, changePct: w.changePct})
	}
	return out
}

var ruMonths = []string{"январь", "февраль", "март", "апрель", "май", "июнь", "июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь"}

func formatMonthYearRU(t time.Time) string {
	m := int(t.Month()) - 1
	if m < 0 || m >= len(ruMonths) {
		return t.Format("January 2006")
	}
	return ruMonths[m] + " " + strconv.Itoa(t.Year())
}

func trend(curr, prev int) string {
	if curr > prev {
		return "↗️"
	}
	if curr < prev {
		return "↘️"
	}
	return "→"
}

func topNByCategory(rows []domain.ABCSnapshotRow, cat domain.ABCCategory, n int) []domain.ABCSnapshotRow {
	var out []domain.ABCSnapshotRow
	for _, r := range rows {
		if r.Category == cat {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Revenue > out[j].Revenue })
	if len(out) > n {
		out = out[:n]
	}
	return out
}

func productTitle(byNm map[int64]*domain.Product, nmID int64) string {
	if p := byNm[nmID]; p != nil {
		return p.Title
	}
	return "nm " + strconv.FormatInt(nmID, 10)
}

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-3]) + "..."
}
