// Defines prometheus metrics used to track the performance of various RPC
// calls exposed by the lfs-edge local rpc over its gRPC server.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Number of gRPC requests served by the rpc service.
	MetricRPCsServed = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "local_rpc_requests",
			Help: "Total number of RPCs served",
		})

	// Number of failed gRPC requests.
	MetricRPCErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "local_rpc_errors",
			Help: "Total number of failed RPC requests",
		})

	// RPC request processing latency is partitioned by the RPC method. It uses
	// custom buckets based on the expected request duration.
	MetricRPCLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "local_rpc_latency_milliseconds",
			Help:       "A latency histogram for RPC requests",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method"},
	)
)
