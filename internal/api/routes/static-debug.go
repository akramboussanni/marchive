//go:build debug
// +build debug

package routes

import (
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/go-chi/chi/v5"
)

// setupStaticRoutes is a no-op in debug mode - frontend runs separately
func setupStaticRoutes(r chi.Router, repos *repo.Repos) {
	// In debug mode, static files are served by the frontend dev server on localhost:5173
}
