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

	logging.Debug("redirecting non TLS url to %s", target)

	http.Redirect(w, r, target, http.StatusTemporaryRedirect)

}

type spaHandler struct {
	StaticPath string
	IndexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// @todo log error
		// if we failed to get the absolute path respond with a 400 bad request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check whether a file exists at the given path
	if err = h.checkFile(filepath.Join(h.StaticPath, path), w, r); err == nil {
		// if directory indexing is disallowed and the filepath is dir, server spa index
		if cfg.AllowDirectoryIndex == false && strings.HasSuffix(r.URL.Path, "/") {
			http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
			return
		}
		// otherwise, use http.FileServer to serve the static dir
		http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
	}
}

func (h spaHandler) checkFile(path string, w http.ResponseWriter, r *http.Request) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// file does not exist, serve IndexPath
			http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))

		}
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		// @todo log error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
