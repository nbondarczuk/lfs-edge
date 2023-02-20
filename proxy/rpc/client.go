package rpc

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.azc.ext.hp.com/Krypton/lfs-edge/proto"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/config"
)

const (
	GET_FILE_PROTOCOL_VERSION = "v1"
)

// Container for all gRPC connection variables.
// Remark: Proxy is used as a general term, not a proper name with reference
// to this module.
type LocalRPCClient struct {
	address    string
	connection *grpc.ClientConn
	proxy      proto.FileProxyIPCClient
}

var (
	logger *zap.Logger
	client *LocalRPCClient
)

// Initialize the storage where the files are kept
func Init(usedLogger *zap.Logger, rpcConfig *config.RpcServer) error {
	logger = usedLogger

	// Start gRPC connecvtion to server.
	address := fmt.Sprintf("%s:%d", rpcConfig.Host, rpcConfig.Port)

	// Save it for later use.
	client = &LocalRPCClient{
		address: address,
	}

	logger.Info("Started GRPC client",
		zap.String("address", address),
	)

	return nil
}
