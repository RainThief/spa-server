package server

import (
	"context"
	"net/http"
	"os"
	"regexp"
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
	server         *http.Server
	router         *mux.Router
	TLS            bool
	redirectServer *http.Server
}

// NewServer creates a new server ready to start listening for REST requests
func NewServer() *Server {
	httpServer := &http.Server{
		ReadTimeout:  httpReadTimeout,
		Handler:      handlers.CombinedLoggingHandler(os.Stdout, http.DefaultServeMux),
		WriteTimeout: httpWriteTimeout,
	}
	TLS := isTLS()
	server := &Server{httpServer, mux.NewRouter(), TLS, &http.Server{}}
	server.server.Addr = ":" + cfg.Port
	if TLS {
		server.server.Addr = ":" + cfg.TLSPort
		server.redirectServer = &http.Server{
			ReadTimeout:  httpReadTimeout,
			Handler:      handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(redirectToTLS)),
			WriteTimeout: httpWriteTimeout,
			Addr:         ":" + cfg.Port,
		}
	}
	server.configureRoutes()
	return server
}

func isTLS() bool {
	certfile := cfg.CertFile
	keyfile := cfg.KeyFile
	if certfile == "" || keyfile == "" {
		logging.Info("Can not use TLS, no certificates provided")
		return false
	}

	if cfg.TLSPort == "" {
		logging.Info("Can not use TLS, no TLS port declared")
		return false
	}

	validPort := regexp.MustCompile(`^[0-9]{1,5}$`).MatchString(cfg.TLSPort)
	if !validPort {
		logging.Error("Can not use TLS, invalid port declared %s", cfg.TLSPort)
	}

	return validPort
}

// ConfigureRoutes declares how all the routing is handled
func (s *Server) configureRoutes() {

	// remove plain text response from default 404 handler
	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	for _, spaConfig := range cfg.SpaDirs {
		spa := spaHandler{StaticPath: spaConfig.StaticPath, IndexFile: spaConfig.IndexFile}
		s.router.Host(spaConfig.HostName).PathPrefix("/").Handler(spa)
	}

	http.Handle("/", s.router)
}

// Start the server listening
func (s *Server) Start() {
	listenAndServe := func(s *Server) error {

		if !s.TLS {
			logging.Info("Server starting; listening on port %s", cfg.Port)
			return s.server.ListenAndServe()
		}

		err := make(chan error)
		go func() {
			logging.Info("Server starting; listening on port %s", cfg.Port)
			err <- s.redirectServer.ListenAndServe()
		}()
		go func() {
			logging.Info("Server starting; listening on port %s", cfg.TLSPort)
			err <- s.server.ListenAndServeTLS(
				cfg.CertFile,
				cfg.KeyFile,
			)
		}()

		return <-err

	}
	if err := listenAndServe(s); err != http.ErrServerClosed {
		logging.Error("Error starting server: %s", err)
	}
}

// Stop the server listening
func (s *Server) Stop() {

	_ = shutdownServer(s.server)

	if s.TLS {
		_ = shutdownServer(s.redirectServer)
	}
}

func shutdownServer(server *http.Server) error {

	if err := server.Shutdown(context.Background()); err != nil {
		logging.Error("Error stopping server: %s", err)
		return err
	}
	logging.Info("Server stopped successfully; releasing binding %s", server.Addr)

	return nil
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
