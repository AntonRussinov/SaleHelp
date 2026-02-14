package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestText_empty(t *testing.T) {
	assert.Equal(t, "", Text(""))
}

func TestText_phone_formats(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"+7 999 123-45-67", "Позвоните +7 999 123-45-67 пожалуйста", "Позвоните [redacted] пожалуйста"},
		{"8 999 1234567", "Телефон 8 999 1234567", "Телефон [redacted]"},
		{"+7(999)123-45-67", "Контакты +7(999)123-45-67", "Контакты [redacted]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Text(tt.in))
		})
	}
}

func TestText_email(t *testing.T) {
	in := "Пишите на support@example.com или info@test.ru"
	want := "Пишите на [redacted] или [redacted]"
	assert.Equal(t, want, Text(in))
}

func TestText_urls(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Сайт https://example.com/page", "Сайт [redacted]"},
		{"Перейдите на www.wildberries.ru", "Перейдите на [redacted]"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, Text(tt.in))
	}
}

func TestText_combined(t *testing.T) {
	in := "Текст  +7 999 111-22-33  и  mail@test.com  и  https://link.ru  конец."
	want := "Текст [redacted] и [redacted] и [redacted] конец."
	assert.Equal(t, want, Text(in))
}

func TestText_no_sensitive(t *testing.T) {
	in := "Просто благодарю за заказ, всё супер!"
	assert.Equal(t, in, Text(in))
}
