package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"gitlab.com/martinfleming/spa-server/internal/config"
	"gitlab.com/martinfleming/spa-server/internal/logging"
	httphandlers "gitlab.com/martinfleming/spa-server/pkg/httpHandlers"
)

var cfg *config.Configuration = &config.Config

const (
	httpReadTimeout = 15 * time.Second
	httpWriteTimeout
)

// Servers collates TLS and non TLS servers with routing and sites configuration
type Servers struct {
	server       *http.Server
	router       *mux.Router
	sites        []config.Site
	tlsServer    *http.Server
	tlsRouter    *mux.Router
	tlsSites     []config.Site
	certificates []tls.Certificate
}

// NewServer creates a new server ready to start listening for REST requests
func NewServer() *Servers {
	server := &Servers{
		&http.Server{},
		mux.NewRouter(),
		[]config.Site{},
		&http.Server{},
		mux.NewRouter(),
		[]config.Site{},
		[]tls.Certificate{},
	}
	// @todo return errors and test
	server.processSites()
	server.configureRoutes()
	server.configureServers()
	return server
}

// sorts each configured site into TLS and NonTLS groups
// TLS sites that redirect from NonTLS are also added to NonTLS group
func (s *Servers) processSites() {
	httpsErr := checkPort("HTTPS")
	httpErr := checkPort("HTTP")

	for _, spaConfig := range cfg.SitesAvailable {
		if httpErr == nil && !config.IsTLSsite(spaConfig) {
			s.sites = append(s.sites, config.Site(spaConfig))
			logging.Debug("No valid certificate information for site %s, setting HTTP only", spaConfig.HostName)
			continue
		}

		if httpsErr == nil {
			s.tlsSites = append(s.tlsSites, config.Site(spaConfig))
			if spaConfig.Redirect {
				s.sites = append(s.sites, config.Site(spaConfig))
				logging.Debug("Setting TLS site %s for non TLS redirect", spaConfig.HostName)
			}
		}
	}
}

// ConfigureRoutes declares how all the routing is handled
func (s *Servers) configureRoutes() {
	// remove plain text response from default 404 handler
	// @todo test http router also uses this hander
	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	for _, site := range s.sites {
		var spa http.Handler = spaHandler{staticPath: site.StaticPath, indexFile: site.IndexFile}
		if site.Redirect {
			spa = httphandlers.RedirectNonTLSHandler{}
		}
		s.router.Host(site.HostName).PathPrefix("/").Handler(compress(spa, site))
	}

	for _, site := range s.tlsSites {
		// @todo test with no cert or invalid cert
		cert, err := tls.LoadX509KeyPair(site.CertFile, site.KeyFile)
		if err == nil {
			s.certificates = append(s.certificates, cert)
			s.tlsRouter.Host(site.HostName).PathPrefix("/").Handler(
				compress(spaHandler{staticPath: site.StaticPath, indexFile: site.IndexFile}, site),
			)
		}
	}
}

func (s *Servers) configureServers() {
	if len(s.sites) > 0 {
		s.server = &http.Server{
			ReadTimeout:  httpReadTimeout,
			Handler:      handlers.CombinedLoggingHandler(os.Stdout, s.router),
			WriteTimeout: httpWriteTimeout,
			Addr:         ":" + cfg.Port,
		}
	}
	if len(s.tlsSites) > 0 {
		s.tlsServer = &http.Server{
			ReadTimeout:  httpReadTimeout,
			Handler:      handlers.CombinedLoggingHandler(os.Stdout, s.tlsRouter),
			WriteTimeout: httpWriteTimeout,
			Addr:         ":" + cfg.TLSPort,
			TLSConfig:    &tls.Config{Certificates: s.certificates},
		}
	}
}

// Start the server listening
func (s *Servers) Start() {
	// @todo test if both servers have not been started (0 sites each )
	listenAndServe := func(s *Servers) error {
		err := make(chan error)
		go func() {
			logging.Info("Server starting; listening on port %s", cfg.Port)
			err <- s.server.ListenAndServe()
		}()
		go func() {
			// @todo this fails as we are not geetn certs from here
			// @todo check that certs are valid?
			logging.Info("Server starting; listening on port %s", cfg.TLSPort)
			err <- s.tlsServer.ListenAndServeTLS("", "")
		}()

		return <-err
	}
	if err := listenAndServe(s); err != http.ErrServerClosed {
		logging.Error("Error starting server: %s", err)
	}
}

// Stop the server listening
func (s *Servers) Stop() {
	_ = shutdownServer(s.server)
	_ = shutdownServer(s.tlsServer)
}

func shutdownServer(server *http.Server) error {
	if err := server.Shutdown(context.Background()); err != nil {
		logging.Error("Error stopping server: %s", err)
		return err
	}
	logging.Info("Server stopped successfully; releasing binding %s", server.Addr)

	return nil
}

func checkPort(serverType string) error {
	var port string
	switch serverType {
	case "HTTP":
		port = cfg.Port
	case "HTTPS":
		port = cfg.TLSPort
	}
	if !regexp.MustCompile(`^[0-9]{1,5}$`).MatchString(port) {
		return logging.LogAndRaiseError("Can not serve %s, invalid port declared %s", serverType, port)
	}
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
