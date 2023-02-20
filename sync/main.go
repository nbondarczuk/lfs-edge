package main

import (
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/config"
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/db"
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/files"	
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/rpc"
)

func main() {
	var err error
	config.Init()
	defer config.Shutdown()

	// Read and parse config file
	if !config.Load(false) {
		panic("lfs-edge sync: config load failed")
	}

	logger := config.GetLogger()

	if err = db.Init(logger, &config.Settings.Database); err != nil {
		panic("lfs-edge sync: database init failed")
	}
	defer db.Shutdown()

	if err = files.Init(logger, &config.Settings.Files); err != nil {
		panic("lfs-edge sync: files init failed")
	}
	defer files.Shutdown()
	
	if err = rpc.Init(logger, &config.Settings.Server); err != nil {
		panic("lfs-edge sync: rpc init failed")
	}
	
	logger.Info("HP Lfs-Edge sync service: Goodbye!")
}
