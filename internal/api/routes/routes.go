package routes

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/api/routes/admin"
	"github.com/akramboussanni/marchive/internal/api/routes/auth"
	"github.com/akramboussanni/marchive/internal/api/routes/books"
	"github.com/akramboussanni/marchive/internal/api/routes/invites"
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
	r.Mount("/api/auth", auth.NewAuthRouter(repos.User, repos.Token, repos.Lockout, repos.RequestCredits))
	r.Mount("/api/books", books.NewBookRouter(repos))
	r.Mount("/api/admin", admin.NewAdminRouter(repos))
	r.Mount("/api/invites", invites.NewInviteRoutes(repos.Invite, repos.User, repos.Token))

	setupStaticRoutes(r)

	return r
}
