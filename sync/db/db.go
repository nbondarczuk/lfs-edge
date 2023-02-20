package db

import (
	"database/sql"

	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/config"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const (
	DB_TYPE = "sqlite3"
)

var (
	dbLogger *zap.Logger
	db       *sql.DB
	dbConfig *config.Database
)

func Init(logger *zap.Logger, config *config.Database) error {
	var err error
	dbLogger = logger
	dbConfig = config

	db, err = sql.Open(DB_TYPE, dbConfig.Path)
	if err != nil {
		dbLogger.Error("Error opening db", zap.Error(err))
		return err
	}
	if err = initDb(); err != nil {
		return err
	}
	dbLogger.Info("Connected to embedded database",
		zap.String("Type: ", DB_TYPE),
		zap.String("Path: ", dbConfig.Path))
	return nil
}

// run initial ddl
func initDb() error {
	sql := `CREATE TABLE IF NOT EXISTS files (
		id INTEGER NOT NULL primary key,
		size INTEGER DEFAULT -1,
		name TEXT NULL,
		md5sum TEXT NULL,
		status INTEGER DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NULL);`
	_, err := db.Exec(sql)
	if err != nil {
		dbLogger.Error("Error initializing db schema", zap.Error(err))
		return err
	}
	return nil
}

// migrate support for sqlite3 from golang-migrate is pending
// see https://github.com/mattes/migrate/issues/165
func migrate() {
	if !dbConfig.Migrate {
		dbLogger.Info("Db config has migrate off. Skipping")
		return
	}
	if dbConfig.Schema == "" {
		dbLogger.Error("Empty schema path in config")
		return
	}
	dbLogger.Info("Applying database schema",
		zap.String("Schema path: ", dbConfig.Schema))
}

func Shutdown() {
	dbLogger.Info("Closing database connection")
	db.Close()
}
