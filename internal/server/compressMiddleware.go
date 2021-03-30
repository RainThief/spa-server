package server

import (
	"compress/gzip"
	"errors"
	"net/http"

	"github.com/RainThief/spa-server/internal/config"
	"github.com/gorilla/handlers"
)

const compressDefaultLevel = gzip.DefaultCompression

func compressMiddleware(next http.Handler, siteConfig config.Site) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !siteConfig.Compress && !cfg.Compress {
			next.ServeHTTP(w, r)
			return
		}

		checkLevel := func(level int) (int, error) {
			if level != 0 && level < 10 {
				return level, nil
			}
			return 0, errors.New("Invalid level number, must be > 0 & < 10")
		}

		var compressLevel int
		for _, level := range [3]int{siteConfig.CompressLevel, cfg.CompressLevel, compressDefaultLevel} {
			if validLevel, err := checkLevel(level); err == nil {
				compressLevel = validLevel
				break
			}
		}

		handlers.CompressHandlerLevel(next, compressLevel).ServeHTTP(w, r)
	})
}
