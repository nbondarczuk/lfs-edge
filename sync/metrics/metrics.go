// Purpose:
// The metrics package is used to initialize and track Prometheus metrics for
// various components in the lfs-edge sync service.
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// RegisterPrometheusMetrics - register prometheus metrics.
func RegisterPrometheusMetrics() {
	prometheus.MustRegister(MetricRPCLatency)
}

// ReportLatencyMetric reports the latency of the specified operation to the
// specified summary vector metric. The label is used to partition the resulting
// histogram.
func ReportLatencyMetric(metric *prometheus.SummaryVec,
	startTime time.Time, label string) {
	duration := time.Since(startTime)
	metric.WithLabelValues(label).Observe(float64(duration.Milliseconds()))
}

// Chronograph is used to measure the time taken by the specified function to
// execute
func Chronograph(logger *zap.Logger, startTime time.Time, functionName string) {
	logger.Info("Execution completed in: ",
		zap.String("Function: ", functionName),
		zap.Duration("Duration (msec): ", time.Since(startTime)),
	)
}
