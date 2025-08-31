package admin

import (
	"net/http"
	"time"

	"github.com/akramboussanni/marchive/internal/middleware"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/go-chi/chi/v5"
)

type AdminRouter struct {
	UserRepo            *repo.UserRepo
	TokenRepo           *repo.TokenRepo
	BookRepo            *repo.BookRepo
	DownloadRequestRepo *repo.DownloadRequestRepo
	RequestCreditsRepo  *repo.RequestCreditsRepo
	RedemptionCodeRepo  *repo.RedemptionCodeRepo
}

func NewAdminRouter(repos *repo.Repos) http.Handler {
	ar := &AdminRouter{
		UserRepo:            repos.User,
		TokenRepo:           repos.Token,
		BookRepo:            repos.Book,
		DownloadRequestRepo: repos.DownloadRequest,
		RequestCreditsRepo:  repos.RequestCredits,
		RedemptionCodeRepo:  repos.RedemptionCode,
	}
	r := chi.NewRouter()

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

		// Redemption codes management
		r.Post("/redemption-codes", ar.HandleCreateRedemptionCode)
		r.Get("/redemption-codes", ar.HandleListRedemptionCodes)
		r.Post("/redemption-codes/{codeID}/revoke", ar.HandleRevokeRedemptionCode)
		r.Delete("/redemption-codes/{codeID}", ar.HandleDeleteRedemptionCode)

		// User management
		r.Post("/users/give-everyone-invite", ar.HandleGiveEveryoneInvite)
	})

	return r
}
