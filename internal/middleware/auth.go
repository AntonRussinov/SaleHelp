package middleware

import (
	"context"
	"salehelp/domain"
	"salehelp/pkg/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserService interface {
	Login(ctx context.Context, data *auth.UserData) (*domain.User, error)
}

const (
	websocketProtocolsHeadersCount = 3
	websocketProtocolMainHeader    = "sandstorm"
	websocketAccessTokenPrefix     = "accessToken_"
	websocketIdtokenPrefix         = "idToken_"
)

func AuthenticationQuery(authenticator auth.Authenticator, service UserService, cookieUpdateFn func(c echo.Context, token string)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userData, token, err := authenticator.AuthenticateWithQuery(c.Request())
			if err != nil {
				return echo.ErrUnauthorized
			}

			user, err := service.Login(c.Request().Context(), userData)
			if err != nil {
				return echo.ErrInternalServerError
			}

			if cookieUpdateFn != nil {
				cookieUpdateFn(c, token)
			}

			return nextWithUser(c, next, user)
		}
	}
}

func AuthenticationHeader(authenticator auth.Authenticator, service UserService, cookieUpdateFn func(c echo.Context, token string)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userData, token, err := authenticator.AuthenticateWithHeader(c.Request())
			if err != nil {
				return echo.ErrUnauthorized
			}

			user, err := service.Login(c.Request().Context(), userData)
			if err != nil {
				return echo.ErrInternalServerError
			}

			if cookieUpdateFn != nil {
				cookieUpdateFn(c, token)
			}

			accessToken := c.Request().Header.Get("Accesstoken")

			return nextWithUserAndAccessToken(c, next, user, accessToken)
		}
	}
}

func AuthenticationCookie(authenticator auth.Authenticator, service UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userData, _, err := authenticator.AuthenticateWithCookie(c.Request())
			if err != nil {
				return echo.ErrUnauthorized
			}

			user, err := service.Login(c.Request().Context(), userData)
			if err != nil {
				return echo.ErrInternalServerError
			}

			return nextWithUser(c, next, user)
		}
	}
}

func AuthenticationWebSockets(authenticator auth.Authenticator, service UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			ctx := req.Context()

			protocols := req.Header.Get("Sec-WebSocket-Protocol")

			parts := strings.Split(protocols, ",")
			if len(parts) != websocketProtocolsHeadersCount {
				return echo.ErrUnauthorized
			}

			accessToken := parts[1]
			accessToken = strings.TrimSpace(accessToken)

			accessToken, ok := strings.CutPrefix(accessToken, websocketAccessTokenPrefix)
			if !ok {
				return echo.ErrUnauthorized
			}

			idToken := parts[2]
			idToken = strings.TrimSpace(idToken)

			idToken, ok = strings.CutPrefix(idToken, websocketIdtokenPrefix)
			if !ok {
				return echo.ErrUnauthorized
			}

			userData, _, err := authenticator.Authenticate(ctx, idToken)
			if err != nil {
				return echo.ErrUnauthorized
			}

			user, err := service.Login(ctx, userData)
			if err != nil {
				return echo.ErrInternalServerError
			}

			ctx = auth.SetUserToContext(ctx, user)
			ctx = auth.SetAccessTokenToContext(ctx, accessToken)
			ctx = auth.SetIDTokenToContext(ctx, idToken)
			newReq := c.Request().WithContext(ctx)
			c.SetRequest(newReq)

			return next(c)
		}
	}
}

func nextWithUser(c echo.Context, next echo.HandlerFunc, user *domain.User) error {
	newCtx := auth.SetUserToContext(c.Request().Context(), user)
	newReq := c.Request().WithContext(newCtx)
	c.SetRequest(newReq)

	return next(c)
}

func nextWithUserAndAccessToken(c echo.Context, next echo.HandlerFunc, user *domain.User, accessToken string) error {
	newCtx := auth.SetUserToContext(c.Request().Context(), user)
	newCtx = auth.SetAccessTokenToContext(newCtx, accessToken)
	newReq := c.Request().WithContext(newCtx)
	c.SetRequest(newReq)

	return next(c)
}
