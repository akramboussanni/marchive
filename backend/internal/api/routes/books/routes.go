package books

import (
	"net/http"
	"time"

	"github.com/akramboussanni/gocode/internal/middleware"
	"github.com/akramboussanni/gocode/internal/repo"
	"github.com/go-chi/chi/v5"
)

type BookRouter struct {
	BookRepo            *repo.BookRepo
	DownloadJobRepo     *repo.DownloadJobRepo
	DownloadRequestRepo *repo.DownloadRequestRepo
	SearchCacheRepo     *repo.SearchCacheRepo
	FavoriteRepo        *repo.FavoriteRepo
}

func NewBookRouter(repos *repo.Repos) http.Handler {
	br := &BookRouter{
		BookRepo:            repos.Book,
		DownloadJobRepo:     repos.DownloadJob,
		DownloadRequestRepo: repos.DownloadRequest,
		SearchCacheRepo:     repos.SearchCache,
		FavoriteRepo:        repos.Favorite,
	}
	r := chi.NewRouter()

	r.Use(middleware.MaxBytesMiddleware(1 << 20))

	r.Group(func(r chi.Router) {
		middleware.AddRatelimit(r, 30, 1*time.Minute)
		r.Get("/explore", br.HandleExplore)
		r.Post("/search", br.HandleSearch)
		r.Get("/job/{jobID}", br.HandleJobStatus)
		r.Get("/{hash}/download", br.HandleDownloadFile)
	})

	r.Group(func(r chi.Router) {
		middleware.AddRatelimit(r, 10, 1*time.Minute)
		middleware.AddAuth(r, repos.User, repos.Token)
		r.Post("/download", br.HandleRequestDownload)
		r.Post("/download/cached", br.HandleCachedDownloadRequest)
		r.Get("/downloads", br.HandleUserDownloads)
		r.Post("/favorite", br.HandleToggleFavorite)
		r.Get("/favorites", br.HandleGetFavorites)
	})

	return r
}
