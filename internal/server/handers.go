package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type spaHandler struct {
	staticPath string
	indexFile  string
}

// ServeHTTP calls HandlerFunc(w, r)
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// @todo log error
		// if we failed to get the absolute path respond with a 400 bad request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check whether a file exists at the given path
	if err = h.checkFile(filepath.Join(h.staticPath, path), w, r); err == nil {
		// if directory indexing is disallowed and the filepath is dir, server spa index
		if !cfg.AllowDirectoryIndex && strings.HasSuffix(r.URL.Path, "/") {
			http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexFile))
			return
		}
		// otherwise, use http.FileServer to serve the static dir
		http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
	}
}

func (h spaHandler) checkFile(path string, w http.ResponseWriter, r *http.Request) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// file does not exist, serve IndexPath
			http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexFile))

		}
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		// @todo log error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
