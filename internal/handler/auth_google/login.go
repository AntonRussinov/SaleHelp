package auth_google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginGoogle вход по Google id_token.
//
// @Summary     Вход через Google
// @Description Принимает id_token от Google Sign-In, проверяет его, находит или создаёт пользователя, возвращает JWT
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body auth_google.LoginByGoogleRequest true "Google id_token"
// @Success     200  {object} auth_google.LoginByGoogleResponse
// @Failure     400  {object} api.HTTPError
// @Failure     401  {object} api.HTTPError
// @Router      /v1/auth/google [post]
func (h *Handler) LoginGoogle(c echo.Context) error {
	var req LoginByGoogleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, accessToken, err := h.service.LoginByGoogle(c.Request().Context(), req.IDToken)
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
