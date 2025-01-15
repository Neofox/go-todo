package middleware

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type CacheConfig struct {
	MaxAge    time.Duration
	Immutable bool
	Public    bool
}

var DefaultCacheConfig = CacheConfig{
	MaxAge:    365 * 24 * time.Hour, // 1 year
	Immutable: true,
	Public:    true,
}

func Cache(cfg CacheConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if os.Getenv("APP_ENV") == "development" {
				w.Header().Set("Cache-Control", "no-store")
				next.ServeHTTP(w, r)
				return
			}

			bodyContent, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				slog.Error("Error reading body", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			etag := fmt.Sprintf(`W/"%x"`, sha256.Sum256(bodyContent))
			w.Header().Set("ETag", etag)

			if match := r.Header.Get("If-None-Match"); match == etag {
				w.WriteHeader(http.StatusNotModified)
			}

			// set cache headers
			var cacheControl []string
			if cfg.Public {
				cacheControl = append(cacheControl, "public")
			}
			cacheControl = append(cacheControl, "max-age="+fmt.Sprint(int(cfg.MaxAge.Seconds())))
			if cfg.Immutable {
				cacheControl = append(cacheControl, "immutable")
			}
			w.Header().Set("Cache-Control", strings.Join(cacheControl, ", "))
			w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
			w.Header().Set("Vary", "Accept-Encoding")

			next.ServeHTTP(w, r)
		})
	}
}
