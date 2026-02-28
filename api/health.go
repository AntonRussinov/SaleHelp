package api

import "time"

type GetHealthCheckResponse struct {
	Version   string    `json:"version"`
	GitCommit string    `json:"gitCommit"`
	GoVersion string    `json:"goVersion"`
	Built     time.Time `json:"built"`
}
