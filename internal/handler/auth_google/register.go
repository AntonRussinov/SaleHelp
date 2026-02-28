package auth_google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Register регистрация по email и паролю.
//
// @Summary     Регистрация
// @Description Создаёт пользователя по email и паролю (не менее 8 символов), возвращает JWT
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body auth_google.RegisterRequest true "Email и пароль"
// @Success     201  {object} auth_google.LoginByGoogleResponse
// @Failure     400  {object} api.HTTPError
// @Failure     409  {object} api.HTTPError "email already registered"
// @Router      /v1/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, accessToken, err := h.service.RegisterEmail(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "email already registered" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}
	return c.JSON(http.StatusCreated, LoginByGoogleResponse{
		AccessToken: accessToken,
		User: UserResponse{
			ID:    user.ID,
			Email: email,
		},
	})
}
