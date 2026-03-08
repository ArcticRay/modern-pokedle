package middleware

import (
	"net/http"
	"time"

	"github.com/ArcticRay/modern-pokedle/internal/observability"
)

func RequestLogger(logger *observability.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Info("request", map[string]any{
				"method":   r.Method,
				"path":     r.URL.Path,
				"duration": time.Since(start).String(),
			})
		})
	}
}
