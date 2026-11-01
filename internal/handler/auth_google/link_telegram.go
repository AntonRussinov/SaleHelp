package auth_google

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// LinkTelegramResponse ответ POST /auth/link-telegram.
type LinkTelegramResponse struct {
	LinkCode  string `json:"linkCode"  example:"A3F7XQ"`
	ExpiresAt string `json:"expiresAt" example:"2024-01-15T10:15:00Z"` // RFC3339
	DeepLink  string `json:"deepLink"  example:"https://t.me/salehelp_bot?start=link_A3F7XQ"`
}

// LinkTelegram создаёт одноразовый код для привязки Telegram к аккаунту. Требует JWT.
//
// @Summary     Код для привязки Telegram
// @Description Генерирует код, который нужно ввести в боте командой /link &lt;код&gt;
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object} auth_google.LinkTelegramResponse
// @Failure     401  {object} api.HTTPError
// @Router      /v1/auth/link-telegram [post]
func (h *Handler) LinkTelegram(c echo.Context) error {
	userID := UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	linkCode, expiresAt, err := h.service.CreateLinkToken(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	deepLink := ""
	if h.botUsername != "" {
		deepLink = fmt.Sprintf("https://t.me/%s?start=link_%s", h.botUsername, linkCode)
	}

	return c.JSON(http.StatusOK, LinkTelegramResponse{
		LinkCode:  linkCode,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		DeepLink:  deepLink,
	})
}
