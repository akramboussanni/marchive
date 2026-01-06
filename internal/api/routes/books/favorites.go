package books

import (
	"net/http"
	"strconv"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/utils"
)

func (br *BookRouter) HandleToggleFavorite(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	req, err := api.DecodeJSON[ToggleFavoriteRequest](w, r)
	if err != nil {
		return
	}

	_, err = br.BookRepo.GetBookByHash(r.Context(), req.BookHash)
	if err != nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "book not found")
		return
	}

	isFavorited, err := br.FavoriteRepo.IsFavorited(r.Context(), user.ID, req.BookHash)
	if err != nil {
		applog.Error("Failed to check favorite status:", err)
		api.WriteInternalError(w)
		return
	}

	if isFavorited {
		err = br.FavoriteRepo.RemoveFavorite(r.Context(), user.ID, req.BookHash)
		if err != nil {
			applog.Error("Failed to remove favorite:", err)
			api.WriteInternalError(w)
			return
		}

		response := ToggleFavoriteResponse{
			IsFavorited: false,
			Message:     "Book removed from favorites",
		}
		api.WriteJSON(w, http.StatusOK, response)
	} else {
		err = br.FavoriteRepo.AddFavorite(r.Context(), user.ID, req.BookHash)
		if err != nil {
			applog.Error("Failed to add favorite:", err)
			api.WriteInternalError(w)
			return
		}

		response := ToggleFavoriteResponse{
			IsFavorited: true,
			Message:     "Book added to favorites",
		}
		api.WriteJSON(w, http.StatusOK, response)
	}
}

func (br *BookRouter) HandleGetFavorites(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

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

	favorites, err := br.FavoriteRepo.GetUserFavorites(r.Context(), user.ID, limit, offset)
	if err != nil {
		applog.Error("Failed to get user favorites:", err)
		api.WriteInternalError(w)
		return
	}

	var books []BookWithStats
	for _, fav := range favorites {
		book, err := br.BookRepo.GetBookByHash(r.Context(), fav.BookHash)
		if err != nil {
			continue
		}

		downloads, err := br.DownloadRequestRepo.GetDownloadsByMD5(r.Context(), fav.BookHash)
		downloadCount := 0
		if err == nil {
			downloadCount = len(downloads)
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
			DownloadCount: downloadCount,
			CreatedAt:     book.CreatedAt,
		}
		books = append(books, bookStats)
	}

	total, err := br.FavoriteRepo.CountUserFavorites(r.Context(), user.ID)
	if err != nil {
		applog.Error("Failed to count user favorites:", err)
		api.WriteInternalError(w)
		return
	}

	response := FavoritesResponse{
		Books: books,
		Pagination: Pagination{
			Limit:   limit,
			Offset:  offset,
			Total:   total,
			HasNext: offset+limit < total,
		},
	}

	api.WriteJSON(w, http.StatusOK, response)
}
