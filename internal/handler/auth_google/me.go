package auth_google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Me GET /api/v1/me — текущий пользователь из JWT (middleware уже проверил токен и положил userID в context).
func (h *Handler) Me(c echo.Context) error {
	userID := UserIDFromContext(c.Request().Context())
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	// Сервис может отдавать user по ID; пока мы не инжектим user repo в handler, отдаём то, что есть в JWT (userID, email из контекста).
	email := EmailFromContext(c.Request().Context())
	return c.JSON(http.StatusOK, UserResponse{
		ID:    userID,
		Email: email,
	})
}
