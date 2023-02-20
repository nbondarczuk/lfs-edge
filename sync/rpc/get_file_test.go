package rpc

import (
	"testing"

	pb "github.azc.ext.hp.com/Krypton/lfs-edge/proto"
	"google.golang.org/grpc/codes"
)

// check get_file with invalid protocol version
// check return status code.
func TestGetFileProtocolVersion(t *testing.T) {
	request := &pb.GetFileRequest{
		Header: newFileRequestHeader(),
		FileID: 0,
	}
	// update header version to invalid
	request.Header.ProtocolVersion = "invalid"

	response, err := gClient.GetFile(gCtx, request)
	if err != nil {
		t.Errorf("GetFile: RPC failed %v", err)
		return
	}

	if response.Header.Status != uint32(codes.InvalidArgument) {
		t.Fatalf("GetFile Status: expected %d, found: %d",
			codes.InvalidArgument, response.Header.Status)
	}
}

// get file with invalid file id
// call must succeed
// return status must be codes.InvalidArgument.
func TestGetFileInvalidFileIdReturnsInvalidArgument(t *testing.T) {
	request := &pb.GetFileRequest{
		Header: newFileRequestHeader(),
		FileID: 0,
	}

	response, err := gClient.GetFile(gCtx, request)
	if err != nil {
		t.Errorf("GetFile: RPC failed %v", err)
		return
	}

	if response.Header.Status != uint32(codes.InvalidArgument) {
		t.Fatalf("GetFile Status: expected %d, found: %d",
			codes.InvalidArgument, response.Header.Status)
	}
}

// get file with invalid file id
// call must succeed
// return RetryAfterSeconds must be -1 as the call is not retriable.
func TestGetFileRetryAfterSecondsSetToNegativeOne(t *testing.T) {
	request := &pb.GetFileRequest{
		Header: newFileRequestHeader(),
		FileID: 0,
	}

	response, err := gClient.GetFile(gCtx, request)
	if err != nil {
		t.Errorf("GetFile: RPC failed %v", err)
		return
	}

	if response.RetryAfterSeconds != DO_NOT_RETRY {
		t.Fatalf("GetFile RetryAfterSeconds: expected %d, found: %d",
			DO_NOT_RETRY, response.RetryAfterSeconds)
	}
}
