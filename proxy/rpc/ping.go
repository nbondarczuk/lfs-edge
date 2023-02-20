package rpc

import (
	"context"
	"time"

	"github.azc.ext.hp.com/Krypton/lfs-edge/proto"
)

// Ping packs/unpacks the gRPC protocol envelope si it has all in/out paramters.
func Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var request = &proto.PingRequest{
		Message: "ping",
	}

	_, err := client.proxy.Ping(ctx, request)

	return err
}
