package middleware

import (
	"salehelp/internal/storage/params/pagination"
	"strconv"

	"github.com/labstack/echo/v4"
)

func parseLimit(c echo.Context) uint32 {
	limitParam := c.QueryParam("limit")
	if limitParam == "" {
		return 0
	}

	limit, err := strconv.ParseUint(limitParam, 10, 32)
	if err != nil {
		return 0
	}

	return uint32(limit)
}

func parsePage(c echo.Context) uint32 {
	pageParam := c.QueryParam("page")
	if pageParam == "" {
		return 1
	}

	page, err := strconv.ParseUint(pageParam, 10, 32)
	if err != nil || page == 0 {
		return 1
	}

	return uint32(page)
}

func parsePaginationPageNumParams(c echo.Context) *pagination.Params {
	limit := parseLimit(c)
	if limit == 0 {
		return nil
	}

	page := parsePage(c)

	offset := (page - 1) * limit

	return &pagination.Params{
		Offset: offset,
		Limit:  limit,
	}
}

func parsePaginationParams(c echo.Context) *pagination.Params {
	limit := parseLimit(c)
	if limit == 0 {
		return nil
	}

	offsetParam := c.QueryParam("offset")

	if offsetParam == "" {
		offsetParam = "0"
	}

	offset, err := strconv.ParseUint(offsetParam, 10, 32)
	if err != nil {
		return nil
	}

	return &pagination.Params{
		Offset: uint32(offset),
		Limit:  limit,
	}
}

func Pagination() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			params := parsePaginationParams(c)

			if params == nil {
				return next(c)
			}

			ctx := pagination.SetToContext(c.Request().Context(), params)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			return next(c)
		}
	}
}

func PageNumPagination() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			params := parsePaginationPageNumParams(c)

			if params == nil {
				return next(c)
			}

			ctx := pagination.SetToContext(c.Request().Context(), params)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)

			return next(c)
		}
	}
}
