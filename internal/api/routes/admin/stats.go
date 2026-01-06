package admin

import (
	"net/http"
	"time"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
)

func (ar *AdminRouter) HandleSystemStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalUsers, err := ar.UserRepo.CountUsers(ctx)
	if err != nil {
		applog.Error("Failed to count users:", err)
		api.WriteInternalError(w)
		return
	}

	totalBooks, err := ar.BookRepo.CountBooks(ctx)
	if err != nil {
		applog.Error("Failed to count books:", err)
		api.WriteInternalError(w)
		return
	}

	totalDownloads, err := ar.DownloadRequestRepo.CountAllDownloads(ctx)
	if err != nil {
		applog.Error("Failed to count downloads:", err)
		api.WriteInternalError(w)
		return
	}

	activeUsers, err := ar.DownloadRequestRepo.CountActiveUsers(ctx, 24*time.Hour)
	if err != nil {
		applog.Error("Failed to count active users:", err)
		api.WriteInternalError(w)
		return
	}

	recentDownloads, err := ar.DownloadRequestRepo.GetRecentDownloads(ctx, 10)
	if err != nil {
		applog.Error("Failed to get recent downloads:", err)
		api.WriteInternalError(w)
		return
	}

	topBooks, err := ar.DownloadRequestRepo.GetTopDownloadedBooks(ctx, 10)
	if err != nil {
		applog.Error("Failed to get top books:", err)
		api.WriteInternalError(w)
		return
	}

	topBooksStats := make([]BookDownloadStats, 0, len(topBooks))
	for _, book := range topBooks {
		topBooksStats = append(topBooksStats, BookDownloadStats{
			Hash:          book["hash"].(string),
			Title:         book["title"].(string),
			Authors:       book["authors"].(string),
			DownloadCount: book["download_count"].(int),
		})
	}

	response := SystemStatsResponse{
		TotalUsers:      totalUsers,
		TotalBooks:      totalBooks,
		TotalDownloads:  totalDownloads,
		ActiveUsers:     activeUsers,
		RecentDownloads: api.EmptyIfNil(recentDownloads),
		TopBooks:        api.EmptyIfNil(topBooksStats),
	}

	api.WriteJSON(w, http.StatusOK, response)
}
