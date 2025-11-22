package auth

import (
	"fmt"
	"net/http"
	"salehelp/internal/middleware"
	"salehelp/pkg/auth"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	auth           auth.Authenticator
	userService    middleware.UserService
	cookieUpdateFn func(c echo.Context, token string)
}

func NewHandler(cfg auth.Config, userService middleware.UserService) (*Handler, error) {
	var authenticator auth.Authenticator
	var err error
	var cookieUpdateFn func(c echo.Context, token string)

	if cfg.EnableAuth {
		authenticator, err = auth.NewAzureAuthenticator(&cfg)
		cookieUpdateFn = secureCookieUpdate()
	} else {
		//authenticator = auth.NewSkipAuthenticator()
		cookieUpdateFn = insecureCookieUpdate()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to init authenticate handler: %w", err)
	}

	return &Handler{
		auth:           authenticator,
		userService:    userService,
		cookieUpdateFn: cookieUpdateFn,
	}, nil
}

func (h *Handler) QueryMiddleware() echo.MiddlewareFunc {
	return middleware.AuthenticationQuery(h.auth, h.userService, h.cookieUpdateFn)
}

func (h *Handler) WebSocketMiddleware() echo.MiddlewareFunc {
	return middleware.AuthenticationWebSockets(h.auth, h.userService)
}

func (h *Handler) HeaderMiddleware() echo.MiddlewareFunc {
	return middleware.AuthenticationHeader(h.auth, h.userService, h.cookieUpdateFn)
}

func (h *Handler) CookieMiddleware() echo.MiddlewareFunc {
	return middleware.AuthenticationCookie(h.auth, h.userService)
}

func secureCookieUpdate() func(c echo.Context, token string) {
	return func(c echo.Context, token string) {
		currentCookie, _ := c.Cookie(auth.SecureCookieName)
		if currentCookie != nil && currentCookie.Value == token {
			return
		}

		c.SetCookie(&http.Cookie{
			Name:     auth.SecureCookieName,
			Path:     "/",
			Value:    token,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		})
	}
}

func insecureCookieUpdate() func(c echo.Context, token string) {
	return func(c echo.Context, token string) {
		currentCookie, _ := c.Cookie(auth.InsecureCookieName)
		if currentCookie != nil && currentCookie.Value == token {
			return
		}

		c.SetCookie(&http.Cookie{
			Name:     auth.InsecureCookieName,
			Path:     "/",
			Value:    token,
			HttpOnly: false,
			SameSite: http.SameSiteNoneMode,
			Secure:   false,
		})
	}
}
