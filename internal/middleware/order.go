package middleware

import (
	"salehelp/internal/storage/params/order"
	"strings"

	"github.com/labstack/echo/v4"
)

func OrderBy() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			params := c.QueryParams()
			value := strings.Join(params["orderBy"], ",")
			fields := strings.Split(value, ",")

			orderParams := make([]order.Param, 0)

			for _, field := range fields {
				direction := order.Asc

				if strings.HasPrefix(field, "-") {
					field = strings.TrimPrefix(field, "-")
					direction = order.Desc
				}

				if field == "" {
					continue
				}

				orderParam := order.Param{
					Field: field,
					Order: direction,
				}
				orderParams = append(orderParams, orderParam)
			}

			if len(orderParams) > 0 {
				ctx := order.SetToContext(c.Request().Context(), orderParams)
				req := c.Request().WithContext(ctx)
				c.SetRequest(req)
			}

			return next(c)
		}
	}
}
