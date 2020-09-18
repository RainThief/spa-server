package server

import (
	"net/http"
)

// HealthCheckHandler responds with 204 status
type HealthCheckHandler struct{}

// ServeHTTP calls HandlerFunc(w, r)
func (h HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
