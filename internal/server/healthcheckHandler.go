package server

import (
	"net/http"
)

// HealthCheckHandler responds with 204 status
type healthCheckHandler struct{}

// ServeHTTP calls HandlerFunc(w, r)
func (h healthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
