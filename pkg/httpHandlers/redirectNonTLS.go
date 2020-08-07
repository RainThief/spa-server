package httphandlers

import (
	"net/http"
	"regexp"

	"gitlab.com/martinfleming/spa-server/internal/logging"
)

// RedirectNonTLSHandler redirects non TLS URLs to TLS with same host
type RedirectNonTLSHandler struct{}

// ServeHTTP calls HandlerFunc(w, r)
func (h RedirectNonTLSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	host := regexp.MustCompile(
		`(.*):[0-9]+$`,
	).ReplaceAllString(r.Host, `$1`)

	target := "https://" + host + r.URL.Path

	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}

	logging.Debug("redirecting non TLS url to %s", target)

	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}
