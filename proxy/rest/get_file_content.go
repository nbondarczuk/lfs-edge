package rest

import (
	"net/http"
	
	"go.uber.org/zap"
	
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/files"	
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/metrics"
)

// parseGetFileContentArgs checks and loads all required request arguments
func parseGetFileContentArgs(w http.ResponseWriter, r *http.Request) (requestID,
	deviceID, id string, err error) {

	// No contraint on optional request id.
	requestID = r.Header.Get(headerRequestID)

	// Retrieve the mandatory file identifier from request path.
	id, err = getPathVariable(r, paramFileID, true)
	if err != nil {
		logger.Error("The required file id path variable was not specified in the request",
			zap.Error(err),	
			zap.String("request_id", requestID),
		)
		sendBadRequestErrorResponse(w)
		metrics.MetricGetFileBadRequests.Inc()
		err = ErrPathVariableMissing
		return
	}

	// Retrieve the mandatory device id from query paramters.
	deviceID = r.URL.Query().Get(paramDeviceID)
	if deviceID == "" {
		logger.Error("The device id provided in the request is empty",
			zap.String("request_id", requestID),
		)
		sendBadRequestErrorResponse(w)
		metrics.MetricGetFileBadRequests.Inc()
		err = ErrQueryParameterMissing
		return
	}

	return
}

// Retrieves content of file with the specified file ID used as
// storage file name. The http.ServeContent handles the Range tag/val
// under wraps providing a chunk from the right offset with the given length.
// In case of no local foung gRPC sync is contected to download the file
// in async mode. 
func GetFileContentHandler(w http.ResponseWriter, r *http.Request) {
	// Gather all request args.
	requestID, deviceID, id, err := parseGetFileContentArgs(w, r)
	if err != nil {
		logger.Error("Failed to parse request arguments",
			zap.Error(err),
		)
		sendBadRequestErrorResponse(w)
		metrics.MetricGetFileBadRequests.Inc()
		return		
	}
	
	// Retrieve content of the file corresponding to this id.
	content, err := files.GetFileContent(requestID, deviceID, id)
	if err != nil {
		logger.Error("Failed to get file content",
			zap.Error(err),
			zap.String("request_id", requestID),
			zap.String("device_id", deviceID),
			zap.String("id", id),
		)		
		sendInternalServerErrorResponse(w)
		metrics.MetricGetFileInternalErrors.Inc()
		return
	}

	// Send the required contents chunk back with appropriate code 200 or 206 or some other
	// in case of gRPC error.
	sendContentResponse(w, r, content)
}
