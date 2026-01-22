package admin

import (
	"net/http"
	"time"

	"github.com/akramboussanni/marchive/internal/middleware"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/akramboussanni/marchive/internal/services"
	"github.com/go-chi/chi/v5"
)

type AdminRouter struct {
	UserRepo            *repo.UserRepo
	TokenRepo           *repo.TokenRepo
	BookRepo            *repo.BookRepo
	DownloadRequestRepo *repo.DownloadRequestRepo
	RequestCreditsRepo  *repo.RequestCreditsRepo
	SettingsRepo        *repo.SettingsRepo
	UserService         *services.UserService
}

func NewAdminRouter(repos *repo.Repos, userService *services.UserService) http.Handler {
	ar := &AdminRouter{
		UserRepo:            repos.User,
		TokenRepo:           repos.Token,
		BookRepo:            repos.Book,
		DownloadRequestRepo: repos.DownloadRequest,
		RequestCreditsRepo:  repos.RequestCredits,
		SettingsRepo:        repos.Settings,
		UserService:         userService,
	}
	r := chi.NewRouter()

	// Create settings handler
	settingsHandler := NewSettingsHandler(repos.Settings)

	r.Use(middleware.MaxBytesMiddleware(1 << 20))

	r.Group(func(r chi.Router) {
		middleware.AddRatelimit(r, 60, 1*time.Minute)
		middleware.AddAuth(r, repos.User, repos.Token)
		r.Use(middleware.AdminOnly)

		r.Get("/stats", ar.HandleSystemStats)
		r.Post("/users/search", ar.HandleSearchUsers)
		r.Get("/users", ar.HandleListUsers)
		r.Post("/users", ar.HandleCreateUser)
		r.Get("/users/{userID}", ar.HandleGetUser)
		r.Put("/users/{userID}", ar.HandleUpdateUser)
		r.Delete("/users/{userID}", ar.HandleDeleteUser)
		r.Post("/users/{userID}/password", ar.HandleChangeUserPassword)
		r.Post("/users/{userID}/invalidate-sessions", ar.HandleInvalidateUserSessions)

		// Request credits management
		r.Post("/users/credits/grant", ar.HandleGrantRequestCredits)

		// Daily download limit management
		r.Post("/users/daily-limit", ar.HandleSetDailyLimit)

		// Settings management
		r.Get("/settings", settingsHandler.HandleGetSettings)
		r.Post("/settings", settingsHandler.HandleUpdateSetting)
	})

	return r
}
