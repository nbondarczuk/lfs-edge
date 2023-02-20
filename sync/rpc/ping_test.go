package rpc

import (
	"testing"

	pb "github.azc.ext.hp.com/Krypton/lfs-edge/proto"
)

func TestPing(t *testing.T) {
	pingRequest := &pb.PingRequest{
		Message: "ping",
	}

	response, err := gClient.Ping(gCtx, pingRequest)
	if err != nil {
		t.Errorf("Ping: RPC failed %v", err)
		return
	}

	if pingRequest.Message != response.Message {
		t.Fatalf("Ping: expected %s, found: %s",
			pingRequest.Message, response.Message)
	}
}
