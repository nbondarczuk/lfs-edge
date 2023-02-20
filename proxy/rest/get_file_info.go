package rest

import (
	"net/http"
	"time"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/files"	
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/metrics"
	"go.uber.org/zap"
)

// Retrieves metadata of local file with the specified file id used as storage
// file name.
func GetFileInfoHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get(headerRequestID)

	// Retrieve the specified file identifier.
	id, err := getPathVariable(r, paramFileID, true)
	if err != nil {
		logger.Error("The required file id path variable was not specified in the request",
			zap.Error(err),						
		)
		sendBadRequestErrorResponse(w)
		metrics.MetricGetFileBadRequests.Inc()
		return
	}

	// Retrieve information about the file corresponding to this id.
	file, err := files.GetFileInfo(id)
	if err != nil {
		if err == files.ErrNotFound {
			logger.Error("File does not exist",
				zap.Error(err),
				zap.String("id", id),
			)
			sendNotFoundErrorResponse(w)
			metrics.MetricGetFileNotFoundErrors.Inc()
			return
		}

		logger.Error("Failed to read information about file",
			zap.Error(err),
		)
		sendInternalServerErrorResponse(w)
		metrics.MetricGetFileInternalErrors.Inc()
		return
	}

	response := common.CommonFileInfoResponse{
		File: common.FileInformation{
			FileID:    file.FileID,
			Checksum:  file.Checksum,
			Size:      file.Size,
			UpdatedAt: file.UpdatedAt,
		},
		RequestID:    requestID,
		ResponseTime: time.Now(),
	}

	// JSON encode and return information about the file.
	err = sendJsonResponse(w, http.StatusOK, response)
	if err != nil {
		metrics.MetricGetFileInternalErrors.Inc()
	}

	metrics.MetricGetFileResponses.Inc()
}
