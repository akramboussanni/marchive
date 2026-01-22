package routes

import (
	"net/http"

	"github.com/akramboussanni/marchive/config"
	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/api/routes/admin"
	"github.com/akramboussanni/marchive/internal/api/routes/auth"
	"github.com/akramboussanni/marchive/internal/api/routes/books"
	"github.com/akramboussanni/marchive/internal/api/routes/invites"
	"github.com/akramboussanni/marchive/internal/middleware"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/services"

	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(repos *repo.Repos) http.Handler {
	r := chi.NewRouter()

	if config.App.TrustIpHeaders {
		r.Use(chimiddleware.RealIP)
	}
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORSHeaders)

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Initialize services
	userService := services.NewUserService(repos.User)

	api.AddSwaggerRoutes(r)
	r.Mount("/api/auth", auth.NewAuthRouter(repos.User, repos.Token, repos.Lockout, repos.RequestCredits))
	r.Mount("/api/books", books.NewBookRouter(repos))
	r.Mount("/api/admin", admin.NewAdminRouter(repos, userService))
	r.Mount("/api/invites", invites.NewInviteRouter(repos.Invite, repos.User, repos.Token))

	// Public settings endpoint (for frontend to check anonymous access)
	r.Get("/api/settings/public", func(w http.ResponseWriter, r *http.Request) {
		anonymousAccessEnabled := repos.Settings.IsAnonymousAccessEnabled(r.Context())
		api.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"status": "success",
			"settings": map[string]interface{}{
				model.SettingAnonymousAccessEnabled: anonymousAccessEnabled,
			},
		})
	})

	setupStaticRoutes(r, repos)

	return r
}
