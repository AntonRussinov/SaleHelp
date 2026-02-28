package auth_google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginGoogle POST /api/v1/auth/google — вход по Google id_token, возвращает JWT и данные пользователя.
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
