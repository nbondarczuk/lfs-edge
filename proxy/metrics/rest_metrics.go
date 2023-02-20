package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// REST request processing latency is partitioned by the REST method. It uses
	// custom buckets based on the expected request duration.
	MetricRestLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "lfsproxy_rest_latency_milliseconds",
			Help:       "A latency histogram for REST requests served",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method"},
	)

	// Number of REST requests received.
	MetricRequestCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "lfsproxy_rest_requests",
		Help:        "Number of requests received",
		ConstLabels: prometheus.Labels{"version": "1"},
	})

	// Number of internal errors encountered when processing get file requests.
	MetricGetFileInternalErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "lfsproxty_rest_get_file_internal_errors",
			Help: "Total number of internal errors encountered processing get file requests",
		})

	// Number of bad get file requests encountered.
	MetricGetFileBadRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "lfsproxy_rest_get_file_bad_requests",
			Help: "Total number of bad get file requests",
		})

	// Number of get file requests where the requested file was not found.
	MetricGetFileNotFoundErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "lfsproxy_rest_get_file_not_found_errors",
			Help: "Total number of get file requests where file was not found",
		})

	// Number of successful get file requests served.
	MetricGetFileResponses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "lfsproxy_rest_get_file_requests",
			Help: "Total number of successful get file requests served",
		})
)
