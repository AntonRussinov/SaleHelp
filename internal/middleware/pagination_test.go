package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"salehelp/internal/middleware"
	"salehelp/internal/storage/params/pagination"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type PaginationTestCase struct {
	name   string
	query  string
	expect *pagination.Params
}

func TestPagination(t *testing.T) {
	tests := []PaginationTestCase{
		{name: "valid", query: "limit=10&offset=20", expect: &pagination.Params{Offset: 20, Limit: 10}},
		{name: "offset is missing", query: "limit=10", expect: &pagination.Params{Offset: 0, Limit: 10}},
		{name: "limit missing is invalid", query: "offset=20", expect: nil},
		{name: "limit is not a number", query: "limit=abc&offset=20", expect: nil},
		{name: "offset is not a number", query: "limit=10&offset=abc", expect: nil},
		{name: "offset is empty", query: "limit=10&offset=", expect: &pagination.Params{Offset: 0, Limit: 10}},
		{name: "limit is empty", query: "limit=&offset=20", expect: nil},
		{name: "limit is zero", query: "limit=0&offset=20", expect: nil},
		{name: "offset is zero", query: "limit=10&offset=0", expect: &pagination.Params{Offset: 0, Limit: 10}},
		{name: "limit is negative", query: "limit=-10&offset=20", expect: nil},
		{name: "offset is negative", query: "limit=10&offset=-20", expect: nil},
	}

	t.Parallel()

	e := echo.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.URL.RawQuery = tt.query
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mw := middleware.Pagination()

			handler := mw(func(c echo.Context) error {
				params := pagination.GetFromContext(c.Request().Context())
				require.Equal(t, tt.expect, params)

				return nil
			})

			_ = handler(c)
		})
	}
}

func TestPageNumPagination(t *testing.T) {
	tests := []PaginationTestCase{
		{name: "valid", query: "limit=10&page=3", expect: &pagination.Params{Offset: 20, Limit: 10}},
		{name: "page is missing", query: "limit=10", expect: &pagination.Params{Offset: 0, Limit: 10}},
		{name: "limit missing is invalid", query: "page=22", expect: nil},
		{name: "limit is not a number", query: "limit=abc&page=2", expect: nil},
		{name: "page is not a number", query: "limit=10&page=abc", expect: &pagination.Params{Offset: 0, Limit: 10}},
		{name: "page is empty", query: "limit=10&page=", expect: &pagination.Params{Offset: 0, Limit: 10}},
		{name: "limit is empty", query: "limit=&page=2", expect: nil},
		{name: "limit is zero", query: "limit=0&page=2", expect: nil},
		{name: "page is zero", query: "limit=10&page=0", expect: &pagination.Params{Offset: 0, Limit: 10}},
		{name: "limit is negative", query: "limit=-10&page=2", expect: nil},
		{name: "page is negative", query: "limit=10&page=-2", expect: &pagination.Params{Offset: 0, Limit: 10}},
	}

	t.Parallel()

	e := echo.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.URL.RawQuery = tt.query
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mw := middleware.PageNumPagination()

			handler := mw(func(c echo.Context) error {
				params := pagination.GetFromContext(c.Request().Context())
				require.Equal(t, tt.expect, params)

				return nil
			})

			_ = handler(c)
		})
	}
}
