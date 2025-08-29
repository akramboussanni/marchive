package routes

import (
	"net/http"

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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("github.com/akramboussanni/marchive"))
	})

	api.AddSwaggerRoutes(r)

	r.Mount("/auth", auth.NewAuthRouter(repos.User, repos.Token, repos.Lockout))
	r.Mount("/books", books.NewBookRouter(repos))
	r.Mount("/admin", admin.NewAdminRouter(repos))

	return r
}
