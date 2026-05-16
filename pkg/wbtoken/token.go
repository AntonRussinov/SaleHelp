package wbtoken

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"salehelp/domain"
)

// Типы токенов (поле acc).
const (
	AccBasic    = 1
	AccTest     = 2
	AccPersonal = 3
	AccService  = 4
)

// Позиции бит в маске s (отсчёт от 0).
const (
	BitContent     = 1  // Контент
	BitPrices      = 3  // Цены и скидки
	BitMarketplace = 4  // Маркетплейс
	BitStatistics  = 5  // Статистика
	BitFeedbacks   = 7  // Вопросы и отзывы
	BitReadOnly    = 30 // Только чтение
)

// requiredBits — категории, без которых salehelp не работает.
var requiredBits = []struct {
	Bit  int
	Name string
}{
	{BitContent, "Контент"},
	{BitPrices, "Цены и скидки"},
	{BitMarketplace, "Маркетплейс"},
	{BitStatistics, "Статистика"},
	{BitFeedbacks, "Вопросы и отзывы"},
}

// Claims — payload WB JWT-токена.
type Claims struct {
	ID  string `json:"id"`  // UUIDv4 — ID токена
	SID string `json:"sid"` // UUIDv4 — уникальный ID продавца на WB
	ACC int    `json:"acc"` // тип токена
	S   uint64 `json:"s"`   // битовая маска доступных категорий
	T   bool   `json:"t"`   // true = тестовый
}

// TokenType возвращает доменный тип токена по полю acc.
// AccPersonal и AccService имеют одинаковые лимиты WB (3 req/sec, burst 6).
func (c *Claims) TokenType() string {
	if c.ACC == AccService || c.ACC == AccPersonal {
		return domain.TokenTypeService
	}
	return domain.TokenTypeBasic
}

func (c *Claims) hasBit(pos int) bool {
	return c.S&(1<<uint(pos)) != 0
}

// ValidationError — ошибка с сообщением для показа пользователю.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// Decode декодирует WB JWT без проверки подписи — только payload.
// Возвращает ValidationError с русским текстом, если токен не является WB JWT.
func Decode(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, &ValidationError{Message: "Это не WB-токен: скопируй токен полностью из личного кабинета Wildberries."}
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, &ValidationError{Message: "Не удалось прочитать токен. Убедись, что скопировал его полностью."}
	}

	var c Claims
	if err := json.Unmarshal(payload, &c); err != nil {
		return nil, &ValidationError{Message: "Не удалось разобрать содержимое токена. Попробуй скопировать ещё раз."}
	}

	if c.SID == "" {
		return nil, &ValidationError{Message: "Токен не содержит идентификатор продавца. Убедись, что это токен WB, а не другого сервиса."}
	}

	return &c, nil
}

// Validate проверяет токен на соответствие требованиям приложения.
// Возвращает *ValidationError с понятным русским описанием проблемы.
func Validate(c *Claims) error {
	if c.ACC == AccTest {
		return &ValidationError{Message: "Тестовый токен не подходит — он работает только с тестовыми данными WB. Создай базовый токен в личном кабинете WB → Настройки → Доступ к API."}
	}

	var missing []string
	for _, req := range requiredBits {
		if !c.hasBit(req.Bit) {
			missing = append(missing, req.Name)
		}
	}
	if len(missing) > 0 {
		return &ValidationError{
			Message: fmt.Sprintf(
				"Токену не хватает доступа к категориям: %s. Пересоздай токен и включи эти категории.",
				strings.Join(missing, ", "),
			),
		}
	}

	if c.hasBit(BitReadOnly) {
		return &ValidationError{Message: "Токен создан только для чтения, а для автоответов на отзывы нужны права на запись. Пересоздай токен без галочки «Только чтение»."}
	}

	return nil
}
