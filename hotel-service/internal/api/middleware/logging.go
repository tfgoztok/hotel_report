package middleware

import (
	"net/http"
	"time"

	"github.com/tfgoztok/hotel-service/pkg/logger"
)

// Logging is a middleware function that logs HTTP requests.
func Logging(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Return an http.HandlerFunc to handle the request
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()  // Record the start time of the request
			next.ServeHTTP(w, r) // Call the next handler in the chain
			// Log the request details including method, path, and duration
			l.Info("Request processed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
	}
}
