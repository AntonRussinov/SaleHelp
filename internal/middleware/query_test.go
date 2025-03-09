package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"salehelp/internal/middleware"
	"salehelp/internal/storage/params/querystring"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type QueryStringUseCase struct {
	name     string
	rawQuery string
	expected string
	message  string
}

func TestQueryString(t *testing.T) {
	tests := []QueryStringUseCase{
		{
			name:     "empty",
			rawQuery: "",
			expected: "",
			message:  "empty query should result in empty query string in context",
		},
		{
			name:     "empty field",
			rawQuery: "query=",
			expected: "",
			message:  "empty field should be ignored",
		},
		{
			name:     "wrong param name",
			rawQuery: "query1=field1",
			expected: "",
			message:  "wrong param name should result in empty query string in context",
		},
		{
			name:     "succes",
			rawQuery: "query=someString",
			expected: "someString",
		},
	}

	t.Parallel()

	e := echo.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.URL.RawQuery = tt.rawQuery
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mw := middleware.QueryString()

			handler := mw(func(c echo.Context) error {
				qs := querystring.GetFromContext(c.Request().Context())

				require.Equal(t, tt.expected, qs, tt.message)

				return nil
			})

			_ = handler(c)
		})
	}
}
