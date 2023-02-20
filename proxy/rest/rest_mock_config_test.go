package rest_test

import (
	"testing"

	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/config"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/files"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/rest"	
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/rpc"	
)

func mockConfig(t *testing.T, path string) {
	t.Helper()
	
	config.InitTestLogger()
	logger := config.GetLogger()
	if logger == nil {
		panic("logger")
	}
	
	if err := files.Init(logger,
		&config.Files{
			Path: path,
		}); err != nil {
		panic(err)
	}

	if err := rpc.Init(logger,
		&config.RpcServer{
			Host: "localhost",
			Port: 8001,
		}); err != nil {
		panic(err)
	}

	rest.InitForTest(logger)
}
