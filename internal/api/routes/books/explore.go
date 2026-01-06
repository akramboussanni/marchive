package books

import (
	"net/http"
	"strconv"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
)

func (br *BookRouter) HandleExplore(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	books, err := br.BookRepo.GetBooksWithDownloadCount(r.Context(), limit, offset)
	if err != nil {
		applog.Error("Failed to get books:", err)
		api.WriteInternalError(w)
		return
	}

	total, err := br.BookRepo.CountBooks(r.Context())
	if err != nil {
		applog.Error("Failed to count books:", err)
		api.WriteInternalError(w)
		return
	}

	response := BookListResponse{
		Books: make([]BookWithStats, 0, len(books)),
		Pagination: Pagination{
			Limit:   limit,
			Offset:  offset,
			Total:   total,
			HasNext: offset+limit < total,
		},
	}

	for _, book := range books {
		bookStats := BookWithStats{
			Hash:          book["hash"].(string),
			Title:         book["title"].(string),
			Authors:       book["authors"].(string),
			Publisher:     book["publisher"].(string),
			Language:      book["language"].(string),
			Format:        book["format"].(string),
			Size:          book["size"].(string),
			CoverURL:      book["cover_url"].(string),
			CoverData:     book["cover_data"].(string),
			Status:        book["status"].(string),
			DownloadCount: book["download_count"].(int),
			CreatedAt:     book["created_at"].(int64),
		}
		response.Books = append(response.Books, bookStats)
	}

	api.WriteJSON(w, http.StatusOK, response)
}
