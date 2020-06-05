package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"gitlab.com/martinfleming/spa-server/internal/config"
	"gitlab.com/martinfleming/spa-server/internal/logging"
)

const (
	httpReadTimeout = 15 * time.Second
	httpWriteTimeout
)

var cfg *config.Configuration = &config.Config

// Server exposes an HTTP endpoint
type Server struct {
	server *http.Server
	router *mux.Router
}

// NewServer creates a new server ready to start listening for REST requests
func NewServer() *Server {
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  httpReadTimeout,
		Handler:      handlers.CombinedLoggingHandler(os.Stdout, http.DefaultServeMux),
		WriteTimeout: httpWriteTimeout,
	}
	server := &Server{httpServer, mux.NewRouter()}
	server.ConfigureRoutes()
	return server
}

// ConfigureRoutes declares how all the routing is handled
func (s *Server) ConfigureRoutes() {

	spa := spaHandler{StaticPath: "/var/www/html", IndexPath: "index.html"}

	// remove plain text response from default 404 handler
	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	s.router.PathPrefix("/").Handler(spa)
	http.Handle("/", s.router)
}

// Start the server listening
func (s *Server) Start() {
	logging.Info("Server starting; listening on port %s", cfg.Port)
	listenAndServe := func(s *Server) error {
		certfile := cfg.CertFile
		keyfile := cfg.KeyFile
		if certfile == "" || keyfile == "" {
			return s.server.ListenAndServe()
		}

		err := make(chan error)
		go func() {
			err <- http.ListenAndServe(":80", handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(redirectToTLS)))
		}()
		go func() {
			err <- s.server.ListenAndServeTLS(
				cfg.CertFile,
				cfg.KeyFile,
			)
		}()

		select {
		case msg := <-err:
			return msg
		}

	}
	if err := listenAndServe(s); err != http.ErrServerClosed {
		logging.Error("Error starting server: %s", err)
	}
}

// Stop the server listening
func (s *Server) Stop() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		logging.Error("Error stopping server: %s", err)
		return
	}
	logging.Info("Server stopped successfully; releasing port %s", cfg.Port)
}

// handleHTTPError sends an internal server error response if an error occurred
// func handleHTTPError(err error, w http.ResponseWriter) bool {
// 	if err != nil {
// 		logging.Error("Server error occurred: %s", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Header().Set("Content-Type", "application/json")
// 		errResponse := fmt.Sprintf("{\"error\": \"%s\"}", err.Error())
// 		_, err = w.Write([]byte(errResponse))
// 		if err != nil {
// 			logging.Fatal("Server error occurred: %s", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 		}
// 		return true
// 	}
// 	return false
// }
