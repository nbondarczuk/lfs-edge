package rpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/config"
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/metrics"
	"go.uber.org/zap"

	pb "github.azc.ext.hp.com/Krypton/lfs-edge/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	rpcLogger       *zap.Logger
	rpcServerConfig *config.Server
	address         string
)

// LocalRPCServer - Connection and other state information for the HP CEM
// large files server local rpc between sync and web components.
type LocalRPCServer struct {
	// sync to web local grpc server
	server *grpc.Server

	// Signal handling to support SIGTERM and SIGINT.
	errChannel  chan error
	stopChannel chan os.Signal
}

// Init - initialize and start the HP CEM lfs-edge local gRPC server
func Init(logger *zap.Logger, serverConfig *config.Server) error {
	rpcLogger = logger
	rpcServerConfig = serverConfig

	s := LocalRPCServer{}

	err := s.NewServer()
	if err != nil {
		rpcLogger.Error("Unable to configure GRPC server. Error!",
			zap.Error(err),
		)
		return err
	}

	address = fmt.Sprintf("%s:%d", rpcServerConfig.Host,
		rpcServerConfig.RpcPort)
	err = s.startServing()
	if err != nil {
		rpcLogger.Error("lfs-edge gRPC server failed to start up.",
			zap.String("Address", address),
			zap.Error(err),
		)
		return err
	}

	rpcLogger.Info("GRPC server initialized",
		zap.String("Address", address))

	s.awaitTermination()
	return nil
}

// NewServer creates and registers a new gRPC server instance for local rpc.
func (s *LocalRPCServer) NewServer() error {
	// Handle SIGTERM and SIGINT.
	s.errChannel = make(chan error)
	s.stopChannel = make(chan os.Signal, 1)
	signal.Notify(s.stopChannel, syscall.SIGINT, syscall.SIGTERM)

	var defaultKeepAliveParams = keepalive.ServerParameters{
		Time:    20 * time.Second,
		Timeout: 5 * time.Second,
	}

	// Initialize and register the gRPC server.
	s.server = grpc.NewServer(
		grpc.KeepaliveParams(defaultKeepAliveParams),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	pb.RegisterFileProxyIPCServer(s.server, s)
	return nil
}

// Start listening on the configured port. Creates a separate goroutine to
// serve gRPC requests.
func (s *LocalRPCServer) startServing() error {
	metrics.RegisterPrometheusMetrics()

	go s.listenAndServe()
	rpcLogger.Info("GRPC server started serving requests",
		zap.String("address", address),
	)

	return nil
}

// Goroutine to listen for and serve gRPC requests.
func (s *LocalRPCServer) listenAndServe() {
	// Start the server and listen to the specified port.
	listener, err := net.Listen("tcp", address)
	if err != nil {
		rpcLogger.Error("Failed to initialize a listener for the gRPC server!",
			zap.Error(err),
			zap.String("address", address),
		)
		s.errChannel <- err
		return
	}

	rpcLogger.Info("GRPC server started listening for client connections",
		zap.String("address", address),
	)

	// Start accepting incoming connection requests.
	err = s.server.Serve(listener)
	if err != nil {
		rpcLogger.Error("Failed to start serving incoming gRPC requests!",
			zap.Error(err),
			zap.String("address", address),
		)
		s.errChannel <- err
		return
	}

	rpcLogger.Info("HP CEM lfs-edge local rpc: Done with serving gRPC requests.",
		zap.String("address", address),
	)
}

// Wait for a signal to shutdown the gRPC server and cleanup.
func (s *LocalRPCServer) awaitTermination() {
	// Block until we receive either an OS signal, or encounter a server
	// fatal error and need to terminate.
	select {
	case err := <-s.errChannel:
		rpcLogger.Error("HP CEM lfs-edge local rpc: Shutting down due to a fatal error.",
			zap.Error(err),
		)
	case sig := <-s.stopChannel:
		rpcLogger.Error("HP CEM lfs-edge local rpc: Received an OS signal and shutting down.",
			zap.String("Signal:", sig.String()),
		)
	}

	// Cleanup.
	s.server.GracefulStop()
}
