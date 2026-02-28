package telegram

import (
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

	// важно: Telegram ждёт быстрый 200 OK
	// поэтому обработку делаем синхронно, но без тяжелых операций
	// (а тяжёлое в addshop мы всё равно будем делать быстро)
	h.bot.HandleUpdate(c.Request().Context(), &upd)

	return c.NoContent(http.StatusOK)
}
