package server

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.com/martinfleming/spa-server/internal/logging"
)

func redirectToTLS(w http.ResponseWriter, r *http.Request) {

	host := regexp.MustCompile(
		`(.*):[0-9]+$`,
	).ReplaceAllString(r.Host, `$1`)

	target := "https://" + host + r.URL.Path

	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}

	logging.Info("redirecting non TLS url to %s", target)

	http.Redirect(w, r, target, http.StatusTemporaryRedirect)

}

// @todo make private
type spaHandler struct {
	StaticPath string
	IndexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// @todo tidy this up stops directory indexing
	if strings.HasSuffix(r.URL.Path, "/") {
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}
