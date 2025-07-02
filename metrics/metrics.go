package metrics

import (
	"fmt"

	"github.com/VictoriaMetrics/metrics"
)

const requestDurationLabel = `http_server_request_duration_milliseconds{route="%s"}`

func recordRequestDuration(route string, duration int64) {
	l := fmt.Sprintf(requestDurationLabel, route)
	metrics.GetOrCreateSummary(l).Update(float64(duration))
}
