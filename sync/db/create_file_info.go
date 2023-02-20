package db

import (
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func CreateFileInfo(id uint64, name string, status int32) error {
	s, err := db.Prepare("INSERT INTO files(id, name, status) VALUES (?, ?, ?)")
	if err != nil {
		dbLogger.Error("Error db CreateFile Prepare", zap.Error(err))
		return err
	}
	defer s.Close()

	_, err = s.Exec(id, name, status)
	if err != nil {
		dbLogger.Error("Error db CreateFile Exec", zap.Error(err))
		return err
	}
	
	return nil
}
