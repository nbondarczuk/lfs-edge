package db

import (
	"os"
	"testing"

	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/config"
)

func TestMain(m *testing.M) {
	var err error

	// Initialize logging for the test run.
	config.InitTestLogger()
	defer config.Shutdown()

	dbFile := "/tmp/sync_test.db"
	
	dbConfig := &config.Database{
		Path: dbFile,
	}

	if err = Init(config.GetLogger(), dbConfig); err != nil {
		os.Exit(1)
	}

	retCode := m.Run()

	dbLogger.Info("Finished running db init test!")

	os.Remove(dbFile)
	
	os.Exit(retCode)
}
