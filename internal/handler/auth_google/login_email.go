package auth_google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginEmail вход по email и паролю.
//
// @Summary     Вход по email
// @Description Проверяет email и пароль, возвращает JWT и данные пользователя
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body auth_google.LoginRequest true "Email и пароль"
// @Success     200  {object} auth_google.LoginByGoogleResponse
// @Failure     400  {object} api.HTTPError
// @Failure     401  {object} api.HTTPError
// @Router      /v1/auth/login [post]
func (h *Handler) LoginEmail(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, accessToken, err := h.service.LoginEmail(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	return c.JSON(http.StatusOK, LoginByGoogleResponse{
		AccessToken: accessToken,
		User: UserResponse{
			ID:    user.ID,
			Email: email,
		},
	})
}
