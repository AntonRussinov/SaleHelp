package handler

import (
	"net/http"
	"os"
	"runtime"
	"salehelp/api"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthCheck
// @id           HealthCheck
// @Description  Get health check data
// @Security     firebase
// @Security     sandstorm-backend
// @Success      200                  {object}  api.GetHealthCheckResponse
// @Failure      400,401,403,500      {object}  api.HTTPError
// @Router       /version [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, api.GetHealthCheckResponse{
		Version:   os.Getenv("VERSION"),
		GitCommit: os.Getenv("GIT_COMMIT"),
		GoVersion: runtime.Version(),
		Built:     time.Now(),
	})
}
