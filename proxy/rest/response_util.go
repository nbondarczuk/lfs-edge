package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.azc.ext.hp.com/Krypton/lfs-edge/common/file"
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/files"		
	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/metrics"	
)

const (
	headerRetryAfterTag = "Retry-After"
)

func sendInternalServerErrorResponse(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func sendBadRequestErrorResponse(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func sendNotFoundErrorResponse(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func sendUnsupportedMediaTypeResponse(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusUnsupportedMediaType),
		http.StatusUnsupportedMediaType)
}

// sendJsonResponse does JSON encode and send the specified payload
// and the specified HTTP status code.
func sendJsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Set(headerContentType, contentTypeJson)
	w.WriteHeader(statusCode)

	if payload != nil {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(payload)
		if err != nil {
			logger.Error("Failed to encode JSON response!",
				zap.Error(err),
			)
			sendInternalServerErrorResponse(w)
			return err
		}
	}

	return nil
}

// sendContentResponse interprets status of the file get operation. The status determines
// what is the side effect and what reponse is given to the client.
func sendContentResponse(w http.ResponseWriter, r *http.Request, content *files.FileContent) {
	// Faithfully use the sync provided value.
	w.Header().Set(headerRetryAfterTag, fmt.Sprintf("%d", content.RetryAfterSeconds))
	// Status dependent actions are taken.
	switch(content.Status) {
	case file.StatusOk:
		// The file was found locally so it can be rendered. It is open so we close it.
		defer content.ServedFile.Close()
		defer metrics.MetricGetFileResponses.Inc()
		http.ServeContent(w, r, content.FileName, time.Time{}, content.ServedFile)
	case file.StatusNotFound:
		// The file was not found neither locally not in the source.
		sendNotFoundErrorResponse(w)
	case file.StatusPending:
		// The sync was requested to download the file in async mode.
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)		
	case file.StatusError:
	case file.StatusUnknown:
		// Something wrong happend during file transfer or in the rpc call to sync.
		sendInternalServerErrorResponse(w)		
	}
}
