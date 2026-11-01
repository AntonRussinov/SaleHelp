package auth_google

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// ConfirmTelegramLinkRequest тело POST /auth/confirm-telegram-link.
type ConfirmTelegramLinkRequest struct {
	Code string `json:"code" example:"A3F7XQ"`
}

// ConfirmTelegramLink — пользователь в боте ввёл email, получил код; заходит на сайт под этим email и вводит код.
//
// @Summary     Подтвердить привязку Telegram (код из бота)
// @Description Ввести код, полученный в боте после /link_email. Требует JWT (владелец email).
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body auth_google.ConfirmTelegramLinkRequest true "Код из бота"
// @Success     200
// @Failure     400  {object} api.HTTPError
// @Failure     401  {object} api.HTTPError
// @Router      /v1/auth/confirm-telegram-link [post]
func (h *Handler) ConfirmTelegramLink(c echo.Context) error {
	userID := UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req ConfirmTelegramLinkRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	if req.Code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code is required"})
	}

	if err := h.service.ConfirmTelegramLink(c.Request().Context(), userID, req.Code); err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "expired") {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Код не найден или истёк. Получи новый код в боте (/link_email)."})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
