package rest

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Route - used to route REST requests received by the service.
type Route struct {
	Name        string           // Name of the route
	Method      string           // REST method
	Path        string           // Resource path
	HandlerFunc http.HandlerFunc // Request handler function.
}

type routes []Route

// List of Routes and corresponding handler functions registered
// with the router.
var registeredRoutes = routes{
	// Health method.
	Route{
		Name:        "GetHealth",
		Method:      http.MethodGet,
		Path:        "/health",
		HandlerFunc: GetHealthHandler,
	},

	// Metrics method.
	Route{
		Name:        "GetMetrics",
		Method:      http.MethodGet,
		Path:        "/metrics",
		HandlerFunc: promhttp.Handler().(http.HandlerFunc),
	},

	///////////////////////////////////////////////////////////////////////////
	//                              API routes                               //
	///////////////////////////////////////////////////////////////////////////

	// Get information about the file corresponding to the specified file ID.
	Route{
		Name:        "GetFileInfo",
		Method:      http.MethodHead,
		Path:        "/api/v1/files/{id:[0-9]+}",
		HandlerFunc: GetFileInfoHandler,
	},

	// Get contents of the specified file ID.
	Route{
		Name:        "GetFileContent",
		Method:      http.MethodGet,
		Path:        "/api/v1/files/{id:[0-9]+}",
		HandlerFunc: GetFileContentHandler,
	},	
}
