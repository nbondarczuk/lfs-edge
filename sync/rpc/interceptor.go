// Implements a common interceptor used to intercept all unary RPC requests
// received by the gRPC server. This interceptor is used to calculate
// request latencies while processing RPC requests, and track RPC error metrics
// and RPC served metrics.
package rpc

import (
	"context"
	"time"

	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Interceptor for unary gRPCs served.
func unaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Calculate and report RPC latency metric when the interceptor is done.
	defer metrics.ReportLatencyMetric(metrics.MetricRPCLatency, start,
		info.FullMethod)

	// Invoke the handler to process the gRPC request and update RPC metrics.
	h, err := handler(ctx, req)
	if err != nil {
		rpcLogger.Error("GRPC request error occured",
			zap.Error(err),
			zap.String("Method:", info.FullMethod),
			zap.String("Duration:", time.Since(start).String()),
		)
		metrics.MetricRPCErrors.Inc()
		return h, err
	}

	rpcLogger.Info("GRPC request processed",
		zap.String("Method:", info.FullMethod),
		zap.String("Duration:", time.Since(start).String()),
	)

	metrics.MetricRPCsServed.Inc()

	return h, err
}
