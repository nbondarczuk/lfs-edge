package db

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// SetFileStatus changes the status of a file identified by id
func SetFileStatus(id uint64, status int32) error {
	s, err := db.Prepare("UPDATE files SET status = ? where id = ?")
	if err != nil {
		dbLogger.Error("Error db SetFileStatus Prepare", zap.Error(err))
		return err
	}
	defer s.Close()

	_, err = s.Exec(status, id)
	if err != nil {
		dbLogger.Error("Error db SetFileStatus Exec", zap.Error(err))
		return err
	}

	dbLogger.Info("Set file status",
		zap.Uint64("id", id),
		zap.Int32("status", status),
	)

	return nil
}

// SetFileStatusWithMetadata changes the status of a file identified by id with extended info
func SetFileStatusWithMetadata(id uint64, status int32, name string, size int64, md5sum string,
	modified time.Time) error {
	s, err := db.Prepare("UPDATE files SET status = ?, name = ?, size = ?, md5sum = ? where id = ?")
	if err != nil {
		dbLogger.Error("Error db SetFileStatusWithMetadata Prepare", zap.Error(err))
		return err
	}
	defer s.Close()

	_, err = s.Exec(status, name, size, md5sum, id)
	if err != nil {
		dbLogger.Error("Error db SetFileStatusWithMetadata Exec", zap.Error(err))
		return err
	}

	dbLogger.Info("Set file status",
		zap.Uint64("id", id),
		zap.Int32("status", status),
	)

	return nil
}
