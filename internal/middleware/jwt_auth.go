package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTValidator проверяет JWT и возвращает userID, email.
type JWTValidator interface {
	ValidateJWT(tokenString string) (userID int, email string, err error)
}

// JWTAuth возвращает middleware: извлекает Bearer token, проверяет JWT, кладёт userID и email в контекст.
// Контекст-ключи задаются в пакете handler (auth_google) — middleware принимает setter для совместимости.
func JWTAuth(validator JWTValidator, setContext func(r *http.Request, userID int, email string) *http.Request) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			token := strings.TrimPrefix(auth, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid authorization header"})
			}
			userID, email, err := validator.ValidateJWT(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			if setContext != nil {
				c.SetRequest(setContext(c.Request(), userID, email))
			}
			return next(c)
		}
	}
}
