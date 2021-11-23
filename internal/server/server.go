package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/RainThief/spa-server/internal/config"
	"github.com/RainThief/spa-server/internal/logging"
	httphandlers "github.com/RainThief/spa-server/pkg/httpHandlers"
)

var logger = logging.Logger

var cfg *config.Configuration = &config.Config

const (
	httpReadTimeout = 15 * time.Second
	httpWriteTimeout
	healthCheckDefaultPort = 8080
)

// Servers collates TLS and non TLS servers with routing and sites configuration
type Servers struct {
	server            *http.Server
	router            *mux.Router
	sites             []config.Site
	tlsServer         *http.Server
	tlsRouter         *mux.Router
	tlsSites          []config.Site
	certificates      []tls.Certificate
	healthCheckServer *http.Server
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
		&http.Server{},
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

		// set site non secure
		if httpErr == nil {
			s.sites = append(s.sites, config.Site(spaConfig))
			logMsg := "Setting HTTP site %s"
			if spaConfig.Redirect {
				logMsg = "Setting HTTP site %s for TLS redirect"
			}
			logger.Info(logMsg, spaConfig.HostName)
		}

		// if available set site secure (TLS)
		if httpsErr == nil {
			logMsg := "No valid certificate information for site %s, setting HTTP only"
			if config.IsTLSsite(spaConfig) {
				logMsg = "Setting TLS site %s"
				s.tlsSites = append(s.tlsSites, config.Site(spaConfig))
			}
			logger.Info(logMsg, spaConfig.HostName)
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

	getHandler := func(site config.Site) http.Handler {
		return compressMiddleware(
			spaHandler{site},
			site,
		)
	}

	for _, site := range s.sites {
		var handler = getHandler(site)
		if site.Redirect {
			handler = httphandlers.RedirectNonTLSHandler{}
		}
		s.router.Host(site.HostName).PathPrefix("/").Handler(handler)
	}

	for _, site := range s.tlsSites {
		// @todo test with no cert or invalid cert
		cert, err := tls.LoadX509KeyPair(site.CertFile, site.KeyFile)
		if err == nil {
			s.certificates = append(s.certificates, cert)
			s.tlsRouter.Host(site.HostName).PathPrefix("/").Handler(getHandler(site))
		}
	}
}

func (s *Servers) configureServers() {
	if len(s.sites) > 0 {
		s.server = configureServer(cfg.Port, handlers.CombinedLoggingHandler(os.Stdout, s.router))
	}
	if len(s.tlsSites) > 0 {
		s.tlsServer = configureServer(cfg.TLSPort, handlers.CombinedLoggingHandler(os.Stdout, s.tlsRouter))
		s.tlsServer.TLSConfig = &tls.Config{Certificates: s.certificates}
	}
	if cfg.HealthCheckPort == 0 {
		cfg.HealthCheckPort = healthCheckDefaultPort
	}
	s.healthCheckServer = configureServer(strconv.Itoa(cfg.HealthCheckPort), healthCheckHandler{})
}

// Start the server listening
func (s *Servers) Start(signalChan chan<- os.Signal) {
	listenAndServe := func(s *Servers) error {
		err := make(chan error)
		if !cfg.DisableHealthCheck {
			go func() {
				logger.Info("Healthcheck server starting; listening on port %v", cfg.HealthCheckPort)
				err <- s.healthCheckServer.ListenAndServe()
			}()
		}
		// @todo test if both servers have not been started (0 sites each )
		go func() {
			logger.Info("HTTP server starting; listening on port %s", cfg.Port)
			err <- s.server.ListenAndServe()
		}()
		go func() {
			// @todo check that certs are valid?
			logger.Info("HTTPS server starting; listening on port %s", cfg.TLSPort)
			err <- s.tlsServer.ListenAndServeTLS("", "")
		}()
		return <-err
	}
	if err := listenAndServe(s); err != http.ErrServerClosed {
		logger.Error("Error starting server: %s", err)
		signalChan <- os.Kill
	}
}

// Stop the server listening; graceful shutdown
func (s *Servers) Stop() {
	var wg sync.WaitGroup
	servers := []*http.Server{
		s.server,
		s.tlsServer,
	}
	if !cfg.DisableHealthCheck {
		servers = append(servers, s.healthCheckServer)
	}
	wg.Add(len(servers))
	for _, server := range servers {
		go func(server *http.Server) {
			_ = shutdownServer(server)
			wg.Done()
		}(server)
	}
	wg.Wait()
}

func shutdownServer(server *http.Server) error {
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error("Error stopping server: %s", err)
		return err
	}
	logger.Info("Server stopped successfully; releasing binding %s", server.Addr)

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
		return logger.LogAndRaiseError("Can not serve %s, invalid port declared %s", serverType, port)
	}
	return nil
}

func configureServer(port string, handler http.Handler) *http.Server {
	return &http.Server{
		ReadTimeout:  httpReadTimeout,
		Handler:      handler,
		WriteTimeout: httpWriteTimeout,
		Addr:         ":" + port,
	}
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
