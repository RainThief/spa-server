package server

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/RainThief/spa-server/internal/config"
)

func expiresMiddleware(
	next func(w http.ResponseWriter, r *http.Request, name string),
	siteConf config.Site,
	path string,
	w http.ResponseWriter,
	r *http.Request,
) {
	// no expires config set, nothing for middleware to do
	if len(siteConf.ExpiresParsed) == 0 && len(cfg.ExpiresParsed) == 0 {
		next(w, r, path)
		return
	}

	for _, config := range [2]config.ExpiresParsed{
		siteConf.ExpiresParsed,
		cfg.ExpiresParsed,
	} {
		if duration, ok := config[filepath.Ext(path)[1:]]; ok {
			w.Header().Set("cache-control", "public")
			w.Header().Set("expires", time.Now().Add(duration).Format(http.TimeFormat))

			break
		}
	}

	next(w, r, path)
}
