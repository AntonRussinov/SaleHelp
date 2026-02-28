package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// AdminJWTValidator проверяет JWT и возвращает userID, email только если роль admin.
type AdminJWTValidator interface {
	ValidateJWTAndRequireAdmin(tokenString string) (userID int, email string, err error)
}

// AdminAuth возвращает middleware: извлекает Bearer token, вызывает ValidateJWTAndRequireAdmin, при успехе кладёт userID и email в контекст. При отсутствии/невалидном токене — 401, при роли не admin — 403.
func AdminAuth(validator AdminJWTValidator, setContext func(r *http.Request, userID int, email string) *http.Request) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			token := strings.TrimPrefix(auth, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid authorization header"})
			}
			userID, email, err := validator.ValidateJWTAndRequireAdmin(token)
			if err != nil {
				if strings.Contains(err.Error(), "admin role required") {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "admin role required"})
				}
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			if setContext != nil {
				c.SetRequest(setContext(c.Request(), userID, email))
			}
			return next(c)
		}
	}
}
