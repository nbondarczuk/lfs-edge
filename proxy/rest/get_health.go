package rest

import (
	"net/http"
)

// GetHealthHandler responds with system health feedback for K8S
func GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	sendJsonResponse(w, http.StatusOK, nil)
}
