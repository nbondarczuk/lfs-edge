// Purpose:
// Implements the GetFile RPC used to get file location if file is in cache
// If file is not in cache, initiate a transfer via sync service.
package rpc

import (
	"context"
	"errors"
	"fmt"

	pb "github.azc.ext.hp.com/Krypton/lfs-edge/proto"
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/files"
	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	DO_NOT_RETRY              = -1
	GET_FILE_PROTOCOL_VERSION = "v1"
)

var (
	ErrorInvalidFileId = errors.New("file id is invalid")
)

// GetFile RPC is used to initiate a file get via files
// or provide the local file path if available..
func (s *LocalRPCServer) GetFile(ctx context.Context,
	request *pb.GetFileRequest) (*pb.GetFileResponse, error) {

	if err := validateGetFileRequest(request); err != nil {
		rpcLogger.Error("Error validating request header",
			zap.Error(err))
		return fileResponseWithCode(codes.InvalidArgument), nil
	}

	// Ask fort file info possibly triggering async download.
	info, err := files.GetFile(request.FileID)
	if err != nil {
		rpcLogger.Error("Error getting file info",
			zap.Error(err))
		return fileResponseWithCode(codes.Internal), nil
	}

	// Handle local file states loaded during info fetch.
	switch info.Status {
	case file.StatusOk:
		return filePresentResponse(info.Path), nil
	case file.StatusNotFound:
		return fileResponseWithCode(codes.NotFound), nil
	case file.StatusPending:
		return filePendingResponse(info.RetryAfterSeconds), nil
	default:
		// Shouldnt be here.
		rpcLogger.Error("Invalid file status",
			zap.Int32("status: ", info.Status))
		return fileResponseWithCode(codes.Internal), nil
	}
}

// validate common header for version
// validate get file request for valid file id
func validateGetFileRequest(request *pb.GetFileRequest) error {
	// check protocol version supported
	if request.Header.ProtocolVersion != GET_FILE_PROTOCOL_VERSION {
		return fmt.Errorf("GetFile does not support protocol version: %s",
			request.Header.ProtocolVersion)
	}
	// check if file id is valid
	if request.FileID <= 0 {
		return fmt.Errorf("GetFile invalid file id: %d, %w",
			request.FileID, ErrorInvalidFileId)
	}

	return nil
}

// base response
func baseFileResponse() *pb.GetFileResponse {
	return &pb.GetFileResponse{
		Header: &pb.FileResponseHeader{
			Status:       uint32(codes.Unknown),
			ResponseTime: timestamppb.Now(),
		},
		RetryAfterSeconds: DO_NOT_RETRY,
	}
}

// get base response with code as provided code
// usually responses that are not ok
func fileResponseWithCode(code codes.Code) *pb.GetFileResponse {
	response := baseFileResponse()
	response.Header.Status = uint32(code)
	return response
}

func filePresentResponse(filePath string) *pb.GetFileResponse {
	response := fileResponseWithCode(codes.OK)
	response.FilePath = filePath
	return response
}

func filePendingResponse(retryAfterSeconds int32) *pb.GetFileResponse {
	response := fileResponseWithCode(codes.Unavailable)
	response.RetryAfterSeconds = retryAfterSeconds	
	return response	
}
