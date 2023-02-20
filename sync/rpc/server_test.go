package rpc

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.azc.ext.hp.com/Krypton/lfs-edge/proto"
)

var (
	gCtx        context.Context
	gClient     pb.FileProxyIPCClient
	gConnection *grpc.ClientConn
)

func newFileRequestHeader() *pb.FileRequestHeader {
	return &pb.FileRequestHeader{
		ProtocolVersion: "v1",
		RequestId:       uuid.New().String(),
		RequestTime:     timestamppb.Now(),
	}
}

func TestMain(m *testing.M) {
	var s *LocalRPCServer
	var err error

	// Initialize logging for the test run.
	config.InitTestLogger()
	rpcLogger = config.GetLogger()
	defer config.Shutdown()

	if s, err = initTestRpcServer(); err != nil {
		rpcLogger.Error("Failed to init rpc server",
			zap.Error(err))
	}
	defer shutdownTestRpcServer(s)
	if err = initConnection(); err != nil {
		rpcLogger.Error("Failed to connect to rpc server",
			zap.Error(err))
	}

	retCode := m.Run()

	rpcLogger.Info("Finished running local RPC server unit tests!")
	os.Exit(retCode)
}

func initTestRpcServer() (*LocalRPCServer, error) {
	var err error
	// pkg level config
	rpcServerConfig = &config.Server{
		Host:    "127.0.0.1",
		RpcPort: 8001,
	}
	address = fmt.Sprintf("%s:%d", rpcServerConfig.Host,
		rpcServerConfig.RpcPort)
	s := LocalRPCServer{}
	if err = s.NewServer(); err != nil {
		return nil, err
	}
	if err = s.startServing(); err != nil {
		return nil, err
	}
	return &s, nil
}

func initConnection() error {
	var err error
	gCtx = context.Background()
	addr := fmt.Sprintf("%s:%d", rpcServerConfig.Host, rpcServerConfig.RpcPort)
	gConnection, err = grpc.DialContext(gCtx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	gClient = pb.NewFileProxyIPCClient(gConnection)
	rpcLogger.Info("Created rpc client")
	return nil
}

func shutdownTestRpcServer(s *LocalRPCServer) {
	s.server.GracefulStop()
}

func shutdownLogger() {
	if rpcLogger != nil {
		rpcLogger.Sync()
	}
}
