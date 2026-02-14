package sanitize

import (
	"regexp"
	"strings"
)

const redacted = "[redacted]"

var (
	// Телефоны: +7, 8, скобки, пробелы, дефисы, без букв в номере (упрощённо).
	phoneRe = regexp.MustCompile(`(?i)(\+7|8)[\s\-\(\)]*(\d[\s\-\(\)]*){10,}`)
	// Email: локаль@домен.
	emailRe = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	// Ссылки: http/https, www, типичные домены.
	urlRe = regexp.MustCompile(`(?i)(https?://[^\s]+|www\.[^\s]+|[a-zA-Z0-9][-a-zA-Z0-9]*\.[a-zA-Z]{2,}(/[^\s]*)?)`)
)

// Text заменяет телефоны, email и ссылки на [redacted]. Возвращает строку без лишних пробелов.
func Text(s string) string {
	if s == "" {
		return s
	}
	out := phoneRe.ReplaceAllString(s, redacted)
	out = emailRe.ReplaceAllString(out, redacted)
	out = urlRe.ReplaceAllString(out, redacted)
	return strings.TrimSpace(normalizeSpaces(out))
}

func normalizeSpaces(s string) string {
	for {
		prev := s
		s = strings.ReplaceAll(s, "  ", " ")
		if s == prev {
			break
		}
	}
	return s
}
