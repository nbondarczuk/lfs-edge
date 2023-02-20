package files

import (
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"	
)

// GetFileInfo - retrieve file stat info using file ID from the storage dir.
func GetFileInfo(id string) (*FileInfo, error) {
	// Check the format of passed file id.
	fileID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logger.Error("The file id value format is invalid, not numeric",
			zap.String("id", id),
		)
		return nil, ErrInvalidArg
	}

	// Use full path name of the file.
	var fileName = storagePath + string(os.PathSeparator) + id
	exist, err := file.PathExists(fileName, file.RegularPath)
	if err != nil {
		logger.Error("The file is not accessible",
			zap.String("path", fileName),
		)
		return nil, ErrNotAllowed
	}
	if !exist {
		logger.Error("The file does not exist",
			zap.String("path", fileName),
		)
		return nil, ErrNotFound
	}

	// Get raw file system info data.
	size, updated, err := file.StaticInfo(fileName)
	if err != nil {
		return nil, ErrNotAllowed
	}

	// Calculate cksum metadata.
	cksum, err := file.ChecksumMD5Info(fileName)
	if err != nil {
		return nil, ErrNotAllowed
	}

	return &FileInfo{
		FileID:    fileID,
		Checksum:  cksum,
		Size:      size,
		UpdatedAt: updated,
	}, nil
}
