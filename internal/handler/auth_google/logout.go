package auth_google

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Logout завершает сессию текущего пользователя.
//
// @Summary     Выход из системы
// @Description Добавляет токен в blocklist до его естественного истечения. После перезапуска сервиса токены из blocklist сбрасываются.
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object} map[string]string
// @Failure     401  {object} api.HTTPError
// @Router      /v1/auth/logout [post]
func (h *Handler) Logout(c echo.Context) error {
	userID := UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	if token != "" {
		_ = h.service.RevokeToken(token)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "logged out"})
}
