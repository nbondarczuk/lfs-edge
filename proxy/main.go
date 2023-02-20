package main

import (
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/config"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/files"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/rest"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/rpc"
)

func main() {
	// init logging and load config
	config.Init()
	defer config.Shutdown()

	// Read and parse the configuration file.
	if !config.Load(false) {
		panic("config load failed.")
	}

	logger := config.GetLogger()

	// Initialize the connection to the storage system.
	err := files.Init(logger, &config.Settings.Files)
	if err != nil {
		panic(err)
	}

	// Initialize the connection to the storage system.
	err = rpc.Init(logger, &config.Settings.RpcServer)
	if err != nil {
		panic(err)
	}

	// Initialize the REST server and start serving requests to the local files
	rest.Init(logger, &config.Settings.Server)
}
