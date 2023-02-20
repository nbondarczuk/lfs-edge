package files

import (
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/rpc"
)

// GetFileContent - provides data for http.ServeContent which handles
// the Range header tags correctly ie. it loads the content using
// header provided values for start and offset of the data part.
// The file descriptor must be closed in defer after usage.
func GetFileContent(requestID, deviceID, id string) (*FileContent, error) {
	// Check the format of passed file id. It must be numeric.
	fileID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logger.Error("The file id value format is invalid, not numeric",
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, ErrInvalidArg
	}

	// Use full path name of the file.
	var (
		status            uint32
		retryAfterSeconds int32
		filePathName      string
	)

	filePathName = storagePath + string(os.PathSeparator) + id
	exist, err := file.PathExists(filePathName, file.RegularPath)
	if err != nil {
		logger.Error("The file is not accessible",
			zap.Error(err),
			zap.String("path", filePathName),
		)
		return nil, ErrNotAllowed
	}
	if !exist {
		// Unhappy path: file is panding and has to be downloaded via sync rpc request.
		logger.Info("The local file does not exist, trying to get it from sync via GRPC",
			zap.String("path", filePathName),
		)

		// Ask for pending file via RPC call to sync.
		status, retryAfterSeconds, filePathName, err = rpc.GetFile(requestID, deviceID, fileID)
		if err != nil {
			logger.Error("The GRPC routing GetFile() error",
				zap.Error(err),
				zap.Uint64("file_id", fileID),
			)
			return nil, ErrGRPCRoutineFailed
		}

		logger.Info("GRPC request feedback received",
			zap.Uint64("file_id", fileID),
			zap.Uint32("status", status),
			zap.Int32("retry", retryAfterSeconds),
			zap.String("path", filePathName),
		)
	}

	// It makes sense to pend file when its name is provided.
	// The file descriptor must be closed after use but in the request handler.
	var fd *os.File
	if status == file.StatusOk {
		var err error
		fd, err = os.Open(filePathName)
		if err != nil {
			logger.Error("The file opening error",
				zap.Error(err),
				zap.String("path", filePathName),
			)
			return nil, ErrNotAllowed
		}
	}

	// Providing values for http.ServeContent so id turns into file name
	// and file descriptor is retuned to be used later.
	return &FileContent{
		DeviceID:          deviceID,
		FileID:            fileID,
		Status:            status,
		FileName:          filePathName,
		RetryAfterSeconds: retryAfterSeconds,
		ServedFile:        fd,
	}, nil
}
