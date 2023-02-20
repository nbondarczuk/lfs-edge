package files

import (
	"go.uber.org/zap"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/config"
)

var (
	logger      *zap.Logger
	storagePath string
)

// Initialize the storage where the files are kept
func Init(usedLogger *zap.Logger, filesConfig *config.Files) error {
	logger = usedLogger

	// Check if the file storage folder exists
	exist, err := file.PathExists(filesConfig.Path, file.DirPath)
	if err != nil {
		logger.Error("The storage dir path is not accessible",
			zap.String("dir", filesConfig.Path),
		)
		return ErrNotAllowed
	}
	if !exist {
		logger.Error("The storage dir path does not exist",
			zap.String("dir", filesConfig.Path),
		)
		return ErrNotFound
	}

	storagePath = filesConfig.Path

	return nil
}
