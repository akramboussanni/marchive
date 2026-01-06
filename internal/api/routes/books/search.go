package books

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/anna"
	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
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

	books, err := anna.FindBook(req.Query)
	if err != nil {
		applog.Error("Failed to search books:", err)
		api.WriteMessage(w, http.StatusInternalServerError, "error", "search failed")
		return
	}

	total := len(books)

	start := req.Offset
	end := start + req.Limit

	if start >= total {
		books = []*anna.Book{}
	} else {
		if end > total {
			end = total
		}
		books = books[start:end]
	}

	availableBooks := make(map[string]bool)
	if len(books) > 0 {
		hashes := make([]string, len(books))
		for i, book := range books {
			hashes[i] = book.Hash
		}

		availableStatuses, err := br.BookRepo.GetBooksAvailabilityByHashes(r.Context(), hashes)
		if err != nil {
			applog.Error("Failed to check books availability:", err)
		} else {
			for _, status := range availableStatuses {
				availableBooks[status.Hash] = status.Status == model.BookStatusReady && status.FilePath != ""
			}
		}
	}

	booksWithStatus := make([]*BookWithStatus, len(books))
	for i, book := range books {
		status := "not_available"
		if availableBooks[book.Hash] {
			status = "available"
		}
		booksWithStatus[i] = &BookWithStatus{
			Book:   book,
			Status: status,
		}
	}

	response := SearchResponse{
		Books: booksWithStatus,
		Total: total,
		Query: req.Query,
		Pagination: Pagination{
			Limit:   req.Limit,
			Offset:  req.Offset,
			Total:   total,
			HasNext: req.Offset+req.Limit < total,
		},
	}

	api.WriteJSON(w, http.StatusOK, response)
}
