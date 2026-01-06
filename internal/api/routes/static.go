//go:build !debug
// +build !debug

package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func setupStaticRoutes(r chi.Router) {
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = "./frontend/build"
	}

	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		return
	}

	// Serve static assets (JS, CSS, etc.)
	assetsPath := filepath.Join(frontendDir, "assets")
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsPath))))
	
	// Serve favicon
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(frontendDir, "favicon.ico"))
	})

	// Catch-all: serve index.html for any route not matched above (SPA routing)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// Don't serve index.html for API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		indexPath := filepath.Join(frontendDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, indexPath)
	})
}