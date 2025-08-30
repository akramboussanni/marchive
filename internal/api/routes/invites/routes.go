package invites

import (
	"net/http"
	"time"

	"github.com/akramboussanni/marchive/internal/middleware"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/go-chi/chi/v5"
)

func NewInviteRoutes(inviteRepo *repo.InviteRepo, userRepo *repo.UserRepo) http.Handler {
	ir := NewInviteRouter(inviteRepo, userRepo)
	r := chi.NewRouter()

	r.Use(middleware.MaxBytesMiddleware(1 << 20))

	// 15/min + auth for invite management
	r.Group(func(r chi.Router) {
		middleware.AddRatelimit(r, 15, 1*time.Minute)
		middleware.AddAuth(r, userRepo, nil)
		r.Post("/", ir.HandleCreateInvite)
		r.Get("/", ir.HandleListInvites)
		r.Post("/{token}/revoke", ir.HandleRevokeInvite)
	})

	// 5/min for invite usage (registration)
	r.Group(func(r chi.Router) {
		middleware.AddRatelimit(r, 5, 1*time.Minute)
		r.Post("/use", ir.HandleUseInvite)
	})

	return r
}
