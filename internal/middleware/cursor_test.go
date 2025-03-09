package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"salehelp/internal/middleware"
	"salehelp/internal/storage/params/cursor"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCursorPagination(t *testing.T) {
	t.Parallel()

	e := echo.New()

	middlewareFunc := middleware.CursorPagination()

	tests := []struct {
		name     string
		limit    int
		cursor   paginator.Cursor
		expected *cursor.Params
	}{
		{
			name:   "valid limit and cursor",
			limit:  10,
			cursor: paginator.Cursor{After: nil, Before: nil},
			expected: &cursor.Params{
				Cursor: &paginator.Cursor{After: nil, Before: nil},
				Limit:  10,
			},
		},
		{
			name:     "limit is zero",
			limit:    0,
			cursor:   paginator.Cursor{After: nil, Before: nil},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// add limit query
			q := req.URL.Query()
			q.Add("limit", strconv.Itoa(tt.limit))
			req.URL.RawQuery = q.Encode()

			// set cursor
			after := ""
			before := ""

			if tt.cursor.After != nil {
				after = *tt.cursor.After
			}

			if tt.cursor.Before != nil {
				before = *tt.cursor.Before
			}

			c.SetParamNames("after", "before")
			c.SetParamValues(after, before)

			handler := middlewareFunc(func(c echo.Context) error {
				ctx := c.Request().Context()
				params := cursor.GetFromContext(ctx)
				assert.Equal(t, tt.expected, params)

				return nil
			})

			err := handler(c)
			require.NoError(t, err)
		})
	}
}
