package middleware

import (
	"salehelp/internal/storage/params/cursor"

	"github.com/labstack/echo/v4"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

func CursorPagination() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			limit := parseLimit(c)

			var currr paginator.Cursor
			_ = c.Bind(&currr)

			if limit == 0 {
				return next(c)
			}

			ctx := cursor.SetToContext(c.Request().Context(), &cursor.Params{
				Cursor: &currr,
				Limit:  limit,
			})
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			return next(c)
		}
	}
}
