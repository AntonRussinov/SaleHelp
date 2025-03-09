package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"salehelp/internal/middleware"
	"salehelp/internal/storage/params/order"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type OrderByTestCase struct {
	name   string
	query  string
	expect []order.Param
}

func TestOrderBy(t *testing.T) {
	tests := []OrderByTestCase{
		{name: "empty", query: "", expect: nil},
		{name: "one field", query: "orderBy=field1", expect: []order.Param{
			{Field: "field1", Order: order.Asc},
		}},
		{name: "two fields", query: "orderBy=field1,-field2", expect: []order.Param{
			{Field: "field1", Order: order.Asc},
			{Field: "field2", Order: order.Desc},
		}},
		{name: "empty field", query: "orderBy=field1,,-field2", expect: []order.Param{
			{Field: "field1", Order: order.Asc},
			{Field: "field2", Order: order.Desc},
		}},
		{name: "wrong param name", query: "order=field1", expect: nil},
		{name: "multiple param name", query: "orderBy=field1,-field2&orderBy=field3", expect: []order.Param{
			{Field: "field1", Order: order.Asc},
			{Field: "field2", Order: order.Desc},
			{Field: "field3", Order: order.Asc},
		}},
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

			mw := middleware.OrderBy()

			handler := mw(func(c echo.Context) error {
				order := order.GetFromContext(c.Request().Context())
				if tt.expect == nil {
					require.Nil(t, order)
				}

				require.ElementsMatch(t, tt.expect, order)

				return nil
			})

			_ = handler(c)
		})
	}
}
