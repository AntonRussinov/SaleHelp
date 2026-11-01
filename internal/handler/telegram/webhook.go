package telegram

import (
	"context"
	"net/http"

	tgdto "salehelp/internal/handler/telegram/dto"

	"github.com/labstack/echo/v4"
)

// Webhook — endpoint для Telegram webhook.
// Telegram шлёт POST JSON update сюда.
//
// @Summary     Telegram webhook
// @Description Принимает Update от Telegram Bot API. Опционально проверяет заголовок X-Telegram-Bot-Api-Secret-Token.
// @Tags        telegram
// @Accept      json
// @Param       body body object true "Telegram Update (JSON)"
// @Success     200
// @Failure     400
// @Failure     403
// @Router      /v1/telegram/webhook [post]
func (h *Handler) Webhook(c echo.Context) error {
	// минимальная безопасность через secret header
	if h.secret != "" {
		got := c.Request().Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if got != h.secret {
			return c.NoContent(http.StatusForbidden)
		}
	}

	var upd tgdto.Update
	if err := c.Bind(&upd); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// Telegram ждёт 200 OK как можно быстрее — обрабатываем в фоне.
	// Используем context.Background(): request-контекст завершится сразу после ответа.
	go h.bot.HandleUpdate(context.Background(), &upd)

	return c.NoContent(http.StatusOK)
}
