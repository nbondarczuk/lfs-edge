// Purpose:
// Implements the Ping RPC used to perform health/uptime checks for
// lfs-edge sync local RPC service.
package rpc

import (
	"context"
	"fmt"

	pb "github.azc.ext.hp.com/Krypton/lfs-edge/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Maximum allowable length of ping messages.
	maxLengthPingMessage = 25
)

// Ping RPC is used to perform health/uptime checks.
func (s *LocalRPCServer) Ping(ctx context.Context,
	request *pb.PingRequest) (*pb.PingResponse, error) {

	// Reject overly long ping requests.
	if len(request.Message) > maxLengthPingMessage {
		return nil, fmt.Errorf("invalid ping request - message too long")
	}

	// Respond with the caller's ping message and the current timestamp to
	// indicate liveness.
	return &pb.PingResponse{
		Message:      request.Message,
		ResponseTime: timestamppb.Now()}, nil
}
