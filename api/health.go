package api

import "time"

type GetHealthCheckResponse struct {
	Version   string    `json:"version"   example:"1.0.0"`
	GitCommit string    `json:"gitCommit" example:"84b789e"`
	GoVersion string    `json:"goVersion" example:"go1.22.0"`
	Built     time.Time `json:"built"     example:"2024-01-15T10:00:00Z"`
}
