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

	// Serve assets directory (Vite outputs JS/CSS here) - MUST be before catch-all
	assetsDir := filepath.Join(frontendDir, "assets")
	r.Get("/assets/*", func(w http.ResponseWriter, req *http.Request) {
		// Strip /assets/ prefix and serve from assets directory
		path := strings.TrimPrefix(req.URL.Path, "/assets/")
		fullPath := filepath.Join(assetsDir, path)
		http.ServeFile(w, req, fullPath)
	})

	// Catch-all: serve index.html for SPA routing
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
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
