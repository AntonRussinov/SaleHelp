package monitoring

import (
	"context"
	"fmt"

	"salehelp/pkg/wbtoken"
)

// CheckWBToken декодирует и валидирует WB JWT-токен без сетевых запросов.
// Возвращает *wbtoken.ValidationError с русским текстом для показа пользователю.
func (s *Service) CheckWBToken(_ context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("токен не может быть пустым")
	}

	claims, err := wbtoken.Decode(token)
	if err != nil {
		return err
	}

	return wbtoken.Validate(claims)
}
