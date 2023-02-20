package files

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/db"
)

// registerFileStatus changes the file name and sets its status in db.
func registerFileStatus(id uint64, name, path string, status int32) (string, error) {
	// In this case we only change the status in db.
	if status != file.StatusOk {
		return "", db.SetFileStatus(id, status)
	}

	// The file status is Ok so it has to be renamed and metadata saved.

	resultPath := storagePath + pathSeparator + name

	err := os.Rename(path, resultPath)
	if err != nil {
		logger.Error("Unable rename file",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("from", path),
			zap.String("to", resultPath),
		)

		return "", db.SetFileStatus(id, file.StatusError)
	}

	// All errors in metadata load are to be treated as file error.

	size, modified, err := file.StaticInfo(resultPath)
	if err != nil {
		logger.Error("Error getting static file info",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("path", resultPath),
		)

		return "", db.SetFileStatus(id, file.StatusError)
	}

	md5sum, err := file.ChecksumMD5Info(resultPath)
	if err != nil {
		logger.Error("Error getting MD5 file info",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("path", resultPath),
		)

		return "", db.SetFileStatus(id, file.StatusError)
	}

	err = db.SetFileStatusWithMetadata(id, status, name, size, md5sum, modified)
	if err != nil {
		logger.Error("Error setting file status",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.Int32("status", status),
		)

		// Last resort, update with metadata failed
		return "", db.SetFileStatus(id, file.StatusError)
	}

	return resultPath, nil
}

// downloadFile handles pending file and starts download from file service.
func downloadFile(id uint64) {
	name := fmt.Sprintf("%d", id)
	pendingFilePath := storagePath + pathSeparator +
		name + ".pending"

	err := db.CreateFileInfo(id, name, file.StatusPending)
	if err != nil {
		logger.Error("Error creating file info",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("name", name),
		)

		return
	}

	logger.Info("Downloading pending file from file service",
		zap.Uint64("id", id),
		zap.String("file", pendingFilePath),
		zap.String("url", fileServerURL),
	)

	var status int32 = file.StatusOk
	err = DownloadFileWithPresignedUrl(fileServerURL, id, pendingFilePath)
	if err != nil {
		logger.Error("Error in pending file download",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("path", pendingFilePath),
		)

		// Handle specific error codes ie. not found
		if err == ErrNotFound {
			status = file.StatusNotFound
		} else {
			status = file.StatusError
		}
	}

	// Finalize file download renaming the pending file to final.
	resultPath, err := registerFileStatus(id, name, pendingFilePath, status)
	if err != nil {
		logger.Error("Error finalizing pending file download",
			zap.Error(err),
			zap.Uint64("id", id),
			zap.String("path", pendingFilePath),
		)

		return
	}

	logger.Info("Pending file download finished",
		zap.Uint64("id", id),
		zap.String("file", resultPath),
		zap.Int32("status", status),
	)
}

// Run reading loop for pending files from client(s) until done.
func runDownloadAgent() {
	var (
		id   uint64
		more bool
	)

	logger.Info("Starting download loop for pending files")

	for {
		// Read from channel id of the file to download till channel close.
		id, more = <-pending
		if more {
			logger.Info("Got new pending file to download", zap.Uint64("id", id))
			// TBD: limit parallelity
			go downloadFile(id)
		} else { // When closed by shutdown we must get inn sync via channel.
			doneWithDownload <- true
			break // the loop
		}
	}

	logger.Info("Finished download agent")
}
