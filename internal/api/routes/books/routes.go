package books

import (
	"net/http"
	"time"

	"github.com/akramboussanni/marchive/internal/middleware"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/go-chi/chi/v5"
)

type BookRouter struct {
	BookRepo            *repo.BookRepo
	DownloadJobRepo     *repo.DownloadJobRepo
	DownloadRequestRepo *repo.DownloadRequestRepo
	FavoriteRepo        *repo.FavoriteRepo
	RequestCreditsRepo  *repo.RequestCreditsRepo
	UserRepo            *repo.UserRepo
}

func NewBookRouter(repos *repo.Repos) http.Handler {
	br := &BookRouter{
		BookRepo:            repos.Book,
		DownloadJobRepo:     repos.DownloadJob,
		DownloadRequestRepo: repos.DownloadRequest,
		FavoriteRepo:        repos.Favorite,
		RequestCreditsRepo:  repos.RequestCredits,
		UserRepo:            repos.User,
	}
	r := chi.NewRouter()

	// Default 1MB limit for most routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.MaxBytesMiddleware(1 << 20))

		r.Group(func(r chi.Router) {
			middleware.AddRatelimit(r, 100, 1*time.Minute)
			middleware.AddOptionalAuth(r, repos.User, repos.Token)
			r.Get("/explore", br.HandleExplore)
			r.Get("/{hash}", br.HandleGetBookDetail)
			r.Post("/search", br.HandleSearch)
		})

		r.Group(func(r chi.Router) {
			middleware.AddRatelimit(r, 50, 1*time.Minute)
			r.Get("/job/{jobID}", br.HandleJobStatus)
		})

		r.Group(func(r chi.Router) {
			middleware.AddRatelimit(r, 30, 1*time.Minute)
			r.Get("/{hash}/download", br.HandleDownloadFile)
		})

		r.Group(func(r chi.Router) {
			middleware.AddRatelimit(r, 30, 1*time.Minute)
			middleware.AddAuth(r, repos.User, repos.Token)
			r.Get("/downloads", br.HandleUserDownloads)
			r.Get("/download-status", br.HandleDownloadStatus)
			r.Get("/favorites", br.HandleGetFavorites)
			r.Post("/favorite", br.HandleToggleFavorite)
		})

		r.Group(func(r chi.Router) {
			middleware.AddRatelimit(r, 15, 1*time.Minute)
			middleware.AddAuth(r, repos.User, repos.Token)
			r.Post("/download", br.HandleRequestDownload)
			r.Post("/ghost-mode", br.HandleUpdateGhostMode)
			r.Post("/delete", br.HandleDeleteBook)
			r.Post("/metadata", br.HandleUpdateBookMetadata)
			r.Post("/restore", br.HandleRestoreBooks) // Admin only
		})
	})


	// Upload routes with 500MB limit
	r.Group(func(r chi.Router) {
		r.Use(middleware.MaxBytesMiddleware(500 << 20)) // 500MB
		middleware.AddRatelimit(r, 5, 1*time.Minute)
		middleware.AddAuth(r, repos.User, repos.Token)
		r.Post("/upload", br.HandleUploadBook)
		r.Put("/{hash}/cover", br.HandleUpdateCover)
	})

	return r
}


