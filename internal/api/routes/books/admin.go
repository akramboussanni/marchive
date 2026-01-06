package books

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
)

func (br *BookRouter) HandleUpdateGhostMode(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	if user.Role != "admin" {
		api.WriteMessage(w, http.StatusForbidden, "error", "admin access required")
		return
	}

	req, err := api.DecodeJSON[UpdateGhostModeRequest](w, r)
	if err != nil {
		return
	}

	err = br.BookRepo.UpdateGhostMode(r.Context(), req.BookHash, req.IsGhost)
	if err != nil {
		applog.Error("Failed to update ghost mode:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "Ghost mode updated successfully")
}

func (br *BookRouter) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	if user.Role != "admin" {
		api.WriteMessage(w, http.StatusForbidden, "error", "admin access required")
		return
	}

	req, err := api.DecodeJSON[DeleteBookRequest](w, r)
	if err != nil {
		return
	}

	err = br.BookRepo.DeleteBook(r.Context(), req.BookHash)
	if err != nil {
		applog.Error("Failed to delete book:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "Book deleted successfully")
}

func (br *BookRouter) HandleUpdateBookMetadata(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	if user.Role != "admin" {
		api.WriteMessage(w, http.StatusForbidden, "error", "admin access required")
		return
	}

	req, err := api.DecodeJSON[UpdateBookMetadataRequest](w, r)
	if err != nil {
		return
	}

	err = br.BookRepo.UpdateBookMetadata(r.Context(), req.BookHash, req.Title, req.Authors, req.Publisher)
	if err != nil {
		applog.Error("Failed to update book metadata:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "Book metadata updated successfully")
}

func (br *BookRouter) HandleGetBookDetail(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		api.WriteMessage(w, http.StatusBadRequest, "error", "hash is required")
		return
	}

	user, hasUser := utils.UserFromContext(r.Context())
	var userID int64
	var isAdmin bool
	if hasUser {
		userID = user.ID
		isAdmin = user.Role == "admin"
	}

	var book *model.SavedBook
	var err error

	if hasUser {
		book, err = br.BookRepo.GetBookByHashForUser(r.Context(), hash, userID, isAdmin)
	} else {
		book, err = br.BookRepo.GetBookByHash(r.Context(), hash)
		// Non-authenticated users can't see ghost books
		if err == nil && book.IsGhost {
			api.WriteMessage(w, http.StatusNotFound, "error", "book not found")
			return
		}
	}

	if err != nil {
		applog.Error("Failed to get book:", err)
		api.WriteMessage(w, http.StatusNotFound, "error", "book not found")
		return
	}

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

	response := BookDetailResponse{
		Book: bookStats,
	}

	// If the book was requested by someone, fetch their info
	if book.RequestedBy != nil && isAdmin {
		requester, err := br.UserRepo.GetUserByID(r.Context(), *book.RequestedBy)
		if err == nil {
			response.RequestedBy = requester
		}
	}

	api.WriteJSON(w, http.StatusOK, response)
}
