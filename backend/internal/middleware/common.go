package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cmdb-lite/backend/internal/config"
)

// CORS adds Cross-Origin Resource Sharing headers to the response
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set allowed origins
			origin := r.Header.Get("Origin")
			if origin != "" {
				// Check if the origin is in the allowed list
				for _, allowedOrigin := range cfg.CORSAllowedOrigins {
					if allowedOrigin == "*" || allowedOrigin == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			// Set allowed methods
			methods := strings.Join(cfg.CORSAllowedMethods, ", ")
			w.Header().Set("Access-Control-Allow-Methods", methods)

			// Set allowed headers
			headers := strings.Join(cfg.CORSAllowedHeaders, ", ")
			w.Header().Set("Access-Control-Allow-Headers", headers)

			// Set credentials policy
			if cfg.CORSAllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Set max age for preflight requests
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

			// Stop here if its a Preflighted OPTIONS request
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Logging logs the HTTP request
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the request
		log.Printf(
			"Started %s %s",
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(w, r)

		// Log the completion time
		log.Printf(
			"Completed %s %s in %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

// Recovery recovers from panics and logs the error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}