package middleware

import (
	"net/http"
	"strings"

	"github.com/akramboussanni/marchive/config"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		if config.App.TLSEnabled || r.TLS != nil {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'wasm-unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https: blob:; font-src 'self' data:;")
		w.Header().Set("Server", "")
		next.ServeHTTP(w, r)
	})
}

func CORSHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Since frontend is now served by backend, we only need minimal CORS for API calls
		// from potential external clients (if any)
		origin := r.Header.Get("Origin")

		if origin != "" {
			// Check if the origin matches our configured domain
			configuredDomain := config.App.Domain

			// Remove port for comparison
			domainWithoutPort := configuredDomain
			if colonIndex := strings.LastIndex(configuredDomain, ":"); colonIndex > 0 {
				domainWithoutPort = configuredDomain[:colonIndex]
			}

			// Check if origin contains our domain (allows for subdomains and ports)
			if strings.Contains(origin, domainWithoutPort) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
