package handler

import (
	"net/http"
	"os"
	"runtime"
	"salehelp/api"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthCheck возвращает версию и метаданные сборки.
//
// @Summary     Health check
// @Description Версия приложения, Go, коммит, время сборки
// @Tags        system
// @Produce     json
// @Success     200 {object} api.GetHealthCheckResponse
// @Router      /version [get]
// Health — liveness probe для Cloud Run.
//
// @Summary     Liveness probe
// @Description Возвращает 200 OK если сервис жив
// @Tags        system
// @Produce     json
// @Success     200 {object} map[string]string
// @Router      /health [get]
func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, api.GetHealthCheckResponse{
		Version:   os.Getenv("VERSION"),
		GitCommit: os.Getenv("GIT_COMMIT"),
		GoVersion: runtime.Version(),
		Built:     time.Now(),
	})
}
