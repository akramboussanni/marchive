package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/api/routes/admin"
	"github.com/akramboussanni/marchive/internal/api/routes/auth"
	"github.com/akramboussanni/marchive/internal/api/routes/books"
	"github.com/akramboussanni/marchive/internal/middleware"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(repos *repo.Repos) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORSHeaders)

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	api.AddSwaggerRoutes(r)
	r.Mount("/api/auth", auth.NewAuthRouter(repos.User, repos.Token, repos.Lockout))
	r.Mount("/api/books", books.NewBookRouter(repos))
	r.Mount("/api/admin", admin.NewAdminRouter(repos))

	// Get frontend directory from environment or use default
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = "./frontend/build"
	}

	// Validate directory exists
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		// Log warning but continue - this allows the API to work without frontend
		// log.Printf("Warning: Frontend directory %s does not exist", frontendDir)
	}

	// Serve static assets
	r.Handle("/_app/*", http.StripPrefix("/_app/", http.FileServer(http.Dir(filepath.Join(frontendDir, "_app")))))
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(frontendDir, "static")))))

	// Serve index.html for SPA routing
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Skip API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Skip static asset routes
		if strings.HasPrefix(r.URL.Path, "/_app/") || strings.HasPrefix(r.URL.Path, "/static/") {
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

	return r
}
