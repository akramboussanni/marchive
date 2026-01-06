package books

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/anna"
	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
)

func (br *BookRouter) HandleSearch(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[SearchRequest](w, r)
	if err != nil {
		return
	}

	if req.Query == "" {
		api.WriteMessage(w, http.StatusBadRequest, "error", "query is required")
		return
	}

	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Default to "all" if not specified
	if req.SearchType == "" {
		req.SearchType = "all"
	}

	// Get user context for filtering
	user, hasUser := utils.UserFromContext(r.Context())
	var userID int64
	var isAdmin bool
	if hasUser {
		userID = user.ID
		isAdmin = user.Role == "admin"
	}

	downloadedBooks := []*BookWithStatus{}
	missingBooks := []*BookWithStatus{}
	var totalDownloaded, totalMissing int

	// Search downloaded books (from database)
	if req.SearchType == "all" || req.SearchType == "downloaded" {
		var dbBooks []model.SavedBook
		var dbTotal int

		if hasUser {
			dbBooks, err = br.BookRepo.SearchBooksForUser(r.Context(), userID, isAdmin, req.Query, req.Limit, req.Offset)
			if err != nil {
				applog.Error("Failed to search database books:", err)
			} else {
				dbTotal, err = br.BookRepo.CountSearchBooks(r.Context(), req.Query)
				if err != nil {
					applog.Error("Failed to count database books:", err)
				}
			}
		} else {
			dbBooks, err = br.BookRepo.SearchBooks(r.Context(), req.Query, req.Limit, req.Offset)
			if err != nil {
				applog.Error("Failed to search database books:", err)
			} else {
				dbTotal, err = br.BookRepo.CountSearchBooks(r.Context(), req.Query)
				if err != nil {
					applog.Error("Failed to count database books:", err)
				}
			}
		}

		// Convert database books to BookWithStatus
		for _, book := range dbBooks {
			annaBook := &anna.Book{
				Hash:      book.Hash,
				Title:     book.Title,
				Authors:   book.Authors,
				Publisher: book.Publisher,
				Language:  book.Language,
				Format:    book.Format,
				Size:      book.Size,
				CoverURL:  book.CoverURL,
				CoverData: book.CoverData,
			}
			status := "not_available"
			if book.Status == model.BookStatusReady && book.FilePath != "" {
				status = "available"
			}
			downloadedBooks = append(downloadedBooks, &BookWithStatus{
				Book:   annaBook,
				Status: status,
			})
		}
		totalDownloaded = dbTotal
	}

	// Search missing books (from Anna) if needed - only for authenticated users
	if hasUser && (req.SearchType == "all" || req.SearchType == "missing") {
		annaBooks, err := anna.FindBook(req.Query)
		if err != nil {
			applog.Error("Failed to search Anna books:", err)
		} else {
			// Get all hashes from Anna results
			hashes := make([]string, len(annaBooks))
			for i, book := range annaBooks {
				hashes[i] = book.Hash
			}

			// Check which books are already in database
			existingBooks := make(map[string]bool)
			if len(hashes) > 0 {
				availableStatuses, err := br.BookRepo.GetBooksAvailabilityByHashes(r.Context(), hashes)
				if err != nil {
					applog.Error("Failed to check books availability:", err)
				} else {
					for _, status := range availableStatuses {
						existingBooks[status.Hash] = true
					}
				}
			}

			// Filter out books that are already in database
			for _, book := range annaBooks {
				if !existingBooks[book.Hash] {
					missingBooks = append(missingBooks, &BookWithStatus{
						Book:   book,
						Status: "not_available",
					})
				}
			}
			totalMissing = len(missingBooks)

			// Apply pagination to missing books
			start := req.Offset
			end := start + req.Limit
			if start >= len(missingBooks) {
				missingBooks = []*BookWithStatus{}
			} else {
				if end > len(missingBooks) {
					end = len(missingBooks)
				}
				missingBooks = missingBooks[start:end]
			}
		}
	}

	total := totalDownloaded + totalMissing

	response := SearchResponse{
		DownloadedBooks: downloadedBooks,
		MissingBooks:    missingBooks,
		Total:           total,
		Query:           req.Query,
		SearchType:      req.SearchType,
		Pagination: Pagination{
			Limit:   req.Limit,
			Offset:  req.Offset,
			Total:   total,
			HasNext: req.Offset+req.Limit < total,
		},
	}

	api.WriteJSON(w, http.StatusOK, response)
}
