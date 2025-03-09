package middleware

import (
	"salehelp/internal/storage/params/querystring"

	"github.com/labstack/echo/v4"
)

func QueryString() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			qs := c.QueryParam("query")

			if qs != "" {
				ctx := querystring.SetToContext(c.Request().Context(), qs)
				req := c.Request().WithContext(ctx)
				c.SetRequest(req)
			}

			return next(c)
		}
	}
}
