package handler_test

import (
	"net/http"
	"net/http/httptest"
	"salehelp/internal/handler"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HealthCheck(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}
