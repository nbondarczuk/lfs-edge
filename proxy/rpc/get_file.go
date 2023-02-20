package rpc

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.azc.ext.hp.com/Krypton/lfs-edge/proto"
)

var (
	ErrInvalidRequestId       = errors.New("request id returned not the same as sent")
	ErrInvalidProtocolVersion = errors.New("wrong value in protocol version")
)

// GetFile packs/unpacks the gRPC protocol envelope si it has all in/out paramters.
// It validates returned values.
func GetFile(requestID string, deviceID string, fileID uint64) (status uint32,
	retry int32, path string, err error) {
	var (
		connection *grpc.ClientConn
		proxy      proto.FileProxyIPCClient
		request    = &proto.GetFileRequest{
			Header: &proto.FileRequestHeader{
				ProtocolVersion: GET_FILE_PROTOCOL_VERSION,
				RequestId:       requestID,
				RequestTime:     timestamppb.Now(),
			},
			// Just busines logic attributes provided here.
			DeviceID: deviceID,
			FileID:   fileID,
		}
		response *proto.GetFileResponse
	)

	// Open TCP connection to gRPC server
	connection, err = grpc.Dial(client.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("The GRPC server is not accessible, no way to dial",
			zap.Error(err),
			zap.String("address", client.address),
		)

		return
	}
	defer connection.Close()

	// Make a client
	proxy = proto.NewFileProxyIPCClient(connection)

	// Just 1 sec for a local call shall be enoung.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// The file is not available locally so sync must download it in async mode.
	response, err = proxy.GetFile(ctx, request)
	if err != nil {
		logger.Error("GRPC service GetFile() error occured",
			zap.Error(err),
			zap.String("request_id", requestID),
			zap.String("device_id", deviceID),
			zap.Uint64("file_id", fileID),
		)

		return
	}

	err = validateGetFileResonse(response, request)
	if err != nil {
		logger.Error("Invalid GRPC service GetFile() response",
			zap.Error(err),
			zap.String("request_id", requestID),
			zap.String("device_id", deviceID),
			zap.Uint64("file_id", fileID),
		)

		return
	}

	// Status info determines treatment in REST layer.
	return response.Header.Status,
		response.RetryAfterSeconds,
		response.FilePath,
		err
}

// validateGetFileResonse checks the response
func validateGetFileResonse(response *proto.GetFileResponse,
	request *proto.GetFileRequest) error {
	// Check request id if the same as sent
	if response.Header.RequestId != request.Header.RequestId {
		logger.Error("RPC service GetFile() wrong request id returned",
			zap.String("request_id", response.Header.RequestId),
			zap.String("expected", request.Header.RequestId),
		)

		return ErrInvalidRequestId
	}

	// Check protocil version if the same as required
	if response.Header.ProtocolVersion != GET_FILE_PROTOCOL_VERSION {
		logger.Error("RPC service GetFile() wrong protocol version returned",
			zap.String("protocol_version", response.Header.ProtocolVersion),
			zap.String("expected", GET_FILE_PROTOCOL_VERSION),
		)

		return ErrInvalidProtocolVersion
	}

	return nil
}
