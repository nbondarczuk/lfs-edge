package files

import (
	"os"
	"time"
)

// File content container serves as a reader interface to file parts
// which in basic case is the whole file.
type FileContent struct {
	DeviceID          string
	FileID            uint64
	Status            uint32
	FileName          string
	RetryAfterSeconds int32
	ServedFile        *os.File
}

// File raw system static metadata with UpdatedAt being the modified
// value of the Linux file system.
// Checksum is SHA256 encoded value calulated using the file content.
type FileInfo struct {
	FileID    uint64
	Checksum  string
	Size      int64
	UpdatedAt time.Time
}
