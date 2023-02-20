package files

import (
	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/db"
)

type FileInfo struct {
	Status            int32
	RetryAfterSeconds int32
	Path              string
	Message           string
}

// GetFile gets information on a file by file id and initiate download if needed.
// It provides interface to RPC method handler. The error case may be handled
// by explicit Status setting or by just error return (caller handles it).  
// As separation of concerns rulez this is the only place where status is applied.
func GetFile(id uint64) (*FileInfo, error) {
	info := &FileInfo{
		// Default case to avoid sloppy error assignement later.
		Status: file.StatusUnknown,
	}

	entry, err := db.GetFile(id)
	if err != nil {
		// Error case
		info.Status = file.StatusError
		info.RetryAfterSeconds = InvalidRetryAfterSeconds
		info.Message = err.Error()
		return info, nil
	}

	if entry.Exist {
		// File already registered so all info in the db. It may error as well.
		info.Status = entry.Status
		if info.Status == file.StatusOk {
			info.Path = storagePath + pathSeparator + entry.Name
		} else {
			info.RetryAfterSeconds = InvalidRetryAfterSeconds
		}
	} else {
		// Async download agent triggered so no name available yet but just a promise.
		info.Status = file.StatusPending
		info.RetryAfterSeconds = retryAfterSeconds
		pending <- id
	}

	return info, nil
}
