package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/RainThief/spa-server/internal/config"
	"github.com/RainThief/spa-server/internal/logging"
)

type spaHandler struct {
	site config.Site
}

// ServeHTTP calls HandlerFunc(w, r)
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		logging.Logger.Error("Cannot get request absolute path: %s", err)
		// if we failed to get the absolute path respond with a 400 bad request
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	path, directory := h.getFilePath(r, path)
	if directory {
		http.FileServer(http.Dir(h.site.StaticPath)).ServeHTTP(w, r)

		return
	}

	expiresMiddleware(http.ServeFile, h.site, path, w, r)
}

// return full path of file to serve for request
func (h spaHandler) getFilePath(r *http.Request, path string) (toServe string, isDirectory bool) {
	index := filepath.Join(h.site.StaticPath, h.site.IndexFile)
	toServe = filepath.Join(h.site.StaticPath, path)

	// if not valid path return index file
	pathStat, err := os.Stat(toServe)
	if err != nil {
		toServe = index
		return
	}

	// if path is directory
	if pathStat.IsDir() {
		// if is index location always serve index file
		// if in subfolder and directory browsing is disabled server index file
		if r.URL.Path == "/" || !cfg.AllowDirectoryIndex {
			toServe = index
			return
		}

		// else use file server so no file to return
		isDirectory = true

		return
	}

	toServe = filepath.Join(h.site.StaticPath, path)

	return
}
