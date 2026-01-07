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

	// Create a file server for the frontend directory
	fs := http.FileServer(http.Dir(frontendDir))

	// Serve static files (assets, favicon, etc.)
	r.Get("/assets/*", func(w http.ResponseWriter, req *http.Request) {
		fs.ServeHTTP(w, req)
	})

	// Serve favicon and other root-level static files
	r.Get("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		fs.ServeHTTP(w, req)
	})

	// Catch-all: serve index.html for SPA routing
	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		// Skip API routes
		if strings.HasPrefix(req.URL.Path, "/api/") {
			http.NotFound(w, req)
			return
		}

		// Check if the requested file exists
		requestedPath := filepath.Join(frontendDir, req.URL.Path)
		if info, err := os.Stat(requestedPath); err == nil && !info.IsDir() {
			// File exists, serve it directly
			fs.ServeHTTP(w, req)
			return
		}

		// Otherwise serve index.html for SPA routing
		indexPath := filepath.Join(frontendDir, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			http.NotFound(w, req)
			return
		}

		http.ServeFile(w, req, indexPath)
	})
}
