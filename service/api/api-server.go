package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-example/shared/api"
	"github.com/nodeset-org/hyperdrive-example/shared/config"
)

// Simple HTTP server manager for API request handling
type ApiServer struct {
	cfgMgr *config.ConfigManager
	logger *slog.Logger
	port   uint16
	router *mux.Router
	server *http.Server
}

// Creates a new API server
func NewApiServer(ip string, port uint16, cfgMgr *config.ConfigManager, logger *slog.Logger, wg *sync.WaitGroup) (*ApiServer, error) {
	router := mux.NewRouter()
	httpServer := &http.Server{
		Handler: router,
	}

	// Create the socket
	socket, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, fmt.Errorf("error creating socket: %w", err)
	}

	// Get the port if random
	if port == 0 {
		port = uint16(socket.Addr().(*net.TCPAddr).Port)
	}

	// Start listening
	wg.Add(1)
	go func() {
		err := httpServer.Serve(socket)
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error while listening for HTTP requests",
				slog.String("error", err.Error()),
			)
		}
		wg.Done()
	}()

	// Create the server
	server := &ApiServer{
		cfgMgr: cfgMgr,
		logger: logger,
		port:   port,
		router: router,
		server: httpServer,
	}

	// Register API routes
	apiRouter := router.PathPrefix("/" + api.ApiRoute).Subrouter()
	apiRouter.HandleFunc("/param", server.HandleParam)

	return server, nil
}

// Stops the server
func (s *ApiServer) Stop() error {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("error stopping listener: %w", err)
	}
	return nil
}

// Get the port the server is running on - useful if the port was automatically assigned
func (s *ApiServer) GetPort() uint16 {
	return s.port
}
