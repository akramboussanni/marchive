package books

import (
	"net/http"
	"strconv"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
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

	// Get user from context if authenticated
	user, hasUser := utils.UserFromContext(r.Context())
	var userID int64
	var isAdmin bool
	if hasUser {
		userID = user.ID
		isAdmin = user.Role == "admin"
	}

	var books []model.SavedBook
	var total int
	var err error

	if hasUser {
		books, err = br.BookRepo.GetBooksForUser(r.Context(), userID, isAdmin, limit, offset)
		if err != nil {
			applog.Error("Failed to get books:", err)
			api.WriteInternalError(w)
			return
		}

		total, err = br.BookRepo.CountBooksForUser(r.Context(), userID, isAdmin)
	} else {
		books, err = br.BookRepo.GetBooks(r.Context(), limit, offset)
		if err != nil {
			applog.Error("Failed to get books:", err)
			api.WriteInternalError(w)
			return
		}

		total, err = br.BookRepo.CountBooks(r.Context())
	}

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
			Hash:          book.Hash,
			Title:         book.Title,
			Authors:       book.Authors,
			Publisher:     book.Publisher,
			Language:      book.Language,
			Format:        book.Format,
			Size:          book.Size,
			CoverURL:      book.CoverURL,
			CoverData:     book.CoverData,
			Status:        book.Status,
			DownloadCount: book.DownloadCount,
			IsGhost:       book.IsGhost,
			RequestedBy:   book.RequestedBy,
			CreatedAt:     book.CreatedAt,
		}
		response.Books = append(response.Books, bookStats)
	}

	api.WriteJSON(w, http.StatusOK, response)
}
