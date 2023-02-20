package db

import (
	"database/sql"
	
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// GetFile checks for file in the db. The case of not found is handles as not error
// and the info is populated with Exist = false.
func GetFile(id uint64) (*FileInfo, error) {
	info := FileInfo{ID: id}
	err := db.QueryRow("SELECT name, size, status FROM files WHERE id=$1", id).
		Scan(&info.Name, &info.Size, &info.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return &info, nil
		}
		dbLogger.Error("Error db getfile", zap.Error(err))
		return nil, err
	}
	info.Exist = true
	return &info, nil
}
