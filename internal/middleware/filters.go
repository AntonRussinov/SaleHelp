package middleware

import (
	"salehelp/internal/storage/params/filter"
	"strings"

	"github.com/labstack/echo/v4"
)

func Filters() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			params := c.QueryParams()

			filters := make([]filter.Param, 0)

			for key, values := range params {
				if !strings.HasPrefix(key, "filter[") {
					continue
				}

				filter := parseFilter(key, values)
				if filter == nil {
					continue
				}

				filters = append(filters, *filter)
			}

			if len(filters) > 0 {
				ctx := filter.SetToContext(c.Request().Context(), filters)
				req := c.Request().WithContext(ctx)
				c.SetRequest(req)
			}

			return next(c)
		}
	}
}

func parseCompareType(key string) filter.CompareType {
	_, key, ok := strings.Cut(key, "[")
	if !ok {
		return filter.CompareTypeEqual
	}

	compareParam, _, ok := strings.Cut(key, "]")
	if !ok {
		return filter.CompareTypeEqual
	}

	compType := filter.CompareType(compareParam)
	if !compType.IsValid() {
		return filter.CompareTypeEqual
	}

	return compType
}

func parseFilter(key string, values []string) *filter.Param {
	key = strings.TrimPrefix(key, "filter[")

	name, key, ok := strings.Cut(key, "]")
	if !ok {
		return nil
	}

	compType := parseCompareType(key)

	filterParam := filter.Param{
		Field:       name,
		CompareType: compType,
		Value:       strings.Join(values, ","),
		Values:      values,
	}

	return &filterParam
}
