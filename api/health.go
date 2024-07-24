package api

import "time"

type GetHealthCheckResponse struct {
	Version   string
	GitCommit string
	GoVersion string
	Built     time.Time
}
