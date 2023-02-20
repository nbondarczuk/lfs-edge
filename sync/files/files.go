package files

import (
	"fmt"
	"os"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"
	"github.azc.ext.hp.com/Krypton/lfs-edge/sync/config"

	"go.uber.org/zap"
)

const InvalidRetryAfterSeconds int32 = -1

var (
	logger *zap.Logger

	// No need to access it from dependency
	pathSeparator = string(os.PathSeparator)

	// Pendig id values sent by client for download
	pending          chan uint64
	// Semaphore for waiting on stopped download agent
	doneWithDownload chan bool

	// Dir where local files are stored
	storagePath string
	// Config value to be returned in the grpc call
	retryAfterSeconds int32
	// File Source URL to get files by id
	fileServerURL string
	// Size of pending channel
	pendingChannelSize int
)

// loadAndCheckConfig checks, stores locally and formats all needed parameters.
func loadAndCheckConfig(filesConfig *config.Files) error {
	// File storage directory
	storagePath = filesConfig.StoragePath	
	// Check if the file storage folder exists
	ok, err := file.PathExists(storagePath, file.DirPath)
	if err != nil {
		logger.Error("The storage dir path is not accessible",
			zap.String("path", storagePath),
			zap.Error(err),
		)
		return ErrNotAllowed
	}
	if !ok {
		logger.Error("The storage dir path does not exist",
			zap.String("path", storagePath),
		)
		return ErrNotFound
	}

	// Template is created using config pattern with %%s for file id.
	fileServerURL = fmt.Sprintf(filesConfig.FileServerURLTemplate,
		filesConfig.FileServerHost, filesConfig.FileServerPort)

	// Size of the buffered channel for pending file id.
	pendingChannelSize = filesConfig.PendingChannelSize

	return nil
}

// Init is main entry point for the package.
func Init(usedLogger *zap.Logger, filesConfig *config.Files) error {
	logger = usedLogger

	if err := loadAndCheckConfig(filesConfig); err != nil {
		return err
	}

	// Make a channel for pending file ids and for termination.
	pending = make(chan uint64, pendingChannelSize)
	doneWithDownload = make(chan bool)
	
	// The download agent - process pending files creating final ones from pending.
	go runDownloadAgent()

	return nil
}

// InitForTest mocker is used in the integration test when start of the sub-servers is not feasible.
func InitForTest(usedLogger *zap.Logger, path string, url string, retry int32, size int, startDownload, startSync bool) {
	logger = usedLogger

	storagePath = path
	fileServerURL = url
	retryAfterSeconds = retry

	pending = make(chan uint64, size)
	doneWithDownload = make(chan bool)

	if startDownload {
		go runDownloadAgent()
	}
}

// Shutdown stops agents sending signal on their done channels.
func Shutdown() {
	// Stop download agent loop
	close(pending)
	<-doneWithDownload
	
	// All clean
	logger.Info("Agent stopped")
}
