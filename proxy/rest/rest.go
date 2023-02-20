package rest

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.azc.ext.hp.com/Krypton/lfs-edge/proxy/config"
)

var (
	logger               *zap.Logger
	debugLogRestRequests bool
)

const (
	// HTTP server timeouts for the REST endpoint.
	readTimeout  = (time.Second * 5)
	writeTimeout = (time.Second * 5)
)

// Represents the FS REST service.
type fsRestService struct {
	// Signal handling to support SIGTERM and SIGINT for the service.
	errChannel  chan error
	stopChannel chan os.Signal

	// Prometheus metrics reporting.
	metricRegistry *prometheus.Registry

	// Request router
	router *mux.Router

	// HTTP port on which the REST server is available.
	port int
}

// Creates a new instance of the FS REST service and initalizes the request
// router for the FS REST endpoint.
func newFsRestService() *fsRestService {
	s := &fsRestService{}

	// Initial signal handling.
	s.errChannel = make(chan error)
	s.stopChannel = make(chan os.Signal, 1)
	signal.Notify(s.stopChannel, syscall.SIGINT, syscall.SIGTERM)

	// Initialize the prometheus metric reporting registry.
	s.metricRegistry = prometheus.NewRegistry()

	s.router = initRequestRouter()
	return s
}

// Starts the HTTP REST server for the FS service and starts serving requests
// at the REST endpoint.
func (s *fsRestService) startServing() {
	// Start the HTTP REST server. http.ListenAndServe() always returns
	// a non-nil error
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", s.port),
		Handler:        s.router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	logger.Error("Received a fatal error from http.ListenAndServe",
		zap.Error(err),
	)

	// Signal the error channel so we can shutdown the service.
	s.errChannel <- err
}

// Waits for the FS REST server to be terminated - either in response to a
// system event received on the stop channel or a fatal error signal received
// on the error channel.
func (s *fsRestService) awaitTermination() {
	select {
	case err := <-s.errChannel:
		logger.Error("Shutting down due to a fatal error.",
			zap.Error(err),
		)
	case sig := <-s.stopChannel:
		logger.Info("Received an OS signal to shut down!",
			zap.String("signal", sig.String()),
		)
	}
}

// Init initializes the FS REST server and starts serving REST requests at the
// FS's REST endpoint.
func Init(usedLogger *zap.Logger, serverConfig *config.Server) {
	logger = usedLogger
	debugLogRestRequests = serverConfig.DebugRestRequests

	s := newFsRestService()
	s.port = serverConfig.Port

	// Initialize the REST server and listen for REST requests on a separate
	// goroutine. Report fatal errors via the error channel.
	go s.startServing()
	logger.Info("Started REST service",
		zap.Int("port", s.port),
	)

	// Wait for the REST server to be terminated either in response to a system
	// event (like service shutdown) or a fatal error.
	s.awaitTermination()
}

func InitForTest(usedLogger *zap.Logger) {
	logger = usedLogger
}
