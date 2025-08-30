package books

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/akramboussanni/marchive/internal/anna"
	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
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

	// Get user for caching (if authenticated)
	var searchCache *model.SearchCache
	user, isAuthenticated := utils.UserFromContext(r.Context())

	books, err := anna.FindBook(req.Query)
	if err != nil {
		applog.Error("Failed to search books:", err)
		api.WriteMessage(w, http.StatusInternalServerError, "error", "search failed")
		return
	}

	total := len(books)

	// Store full results in cache if user is authenticated
	if isAuthenticated && len(books) > 0 {
		searchCache, err = br.SearchCacheRepo.StoreSearchResults(r.Context(), user.ID, req.Query, books, total)
		if err != nil {
			applog.Error("Failed to cache search results:", err)
			// Continue without cache - not a critical error
		}
	}

	// Paginate results for response
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

	// Check which books are already available in the database
	availableBooks := make(map[string]bool)
	if len(books) > 0 {
		hashes := make([]string, len(books))
		for i, book := range books {
			hashes[i] = book.Hash
		}

		// Get availability status for all books in this page
		availableStatuses, err := br.BookRepo.GetBooksAvailabilityByHashes(r.Context(), hashes)
		if err != nil {
			applog.Error("Failed to check books availability:", err)
			// Continue without availability info - not a critical error
		} else {
			for _, status := range availableStatuses {
				availableBooks[status.Hash] = status.Status == model.BookStatusReady && status.FilePath != ""
			}
		}
	}

	// Convert anna.Book to BookWithStatus with availability info
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

	// Add cache info if available
	if searchCache != nil {
		response.SearchID = searchCache.ID
		response.ExpiresAt = searchCache.ExpiresAt
	}

	api.WriteJSON(w, http.StatusOK, response)
}

func (br *BookRouter) HandleRequestDownload(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	req, err := api.DecodeJSON[DownloadRequest](w, r)
	if err != nil {
		return
	}

	// Check if user already requested this book
	hasRequested, err := br.DownloadRequestRepo.HasUserRequestedBook(r.Context(), user.ID, req.Hash)
	if err != nil {
		applog.Error("Failed to check if book already requested:", err)
		api.WriteInternalError(w)
		return
	}

	if hasRequested {
		api.WriteMessage(w, http.StatusConflict, "error", "you have already requested this book")
		return
	}

	canDownload, err := br.DownloadRequestRepo.CheckAndCreateDownload(r.Context(), user.ID, req.Hash, req.Title)
	if err != nil {
		applog.Error("Failed to check download limit:", err)
		api.WriteInternalError(w)
		return
	}

	if !canDownload {
		api.WriteMessage(w, http.StatusTooManyRequests, "error", "daily download limit reached")
		return
	}

	existingBook, err := br.BookRepo.GetBookByHash(r.Context(), req.Hash)
	if err == nil && existingBook.Status == model.BookStatusReady {
		response := DownloadResponse{
			JobID:   0,
			Status:  "ready",
			Message: "Book is already available for download",
		}
		api.WriteJSON(w, http.StatusOK, response)
		return
	}

	job, err := br.DownloadJobRepo.CreateJob(r.Context(), user.ID, req.Hash)
	if err != nil {
		applog.Error("Failed to create download job:", err)
		api.WriteInternalError(w)
		return
	}

	response := DownloadResponse{
		JobID:   job.ID,
		Status:  job.Status,
		Message: "Download job created successfully",
	}

	api.WriteJSON(w, http.StatusCreated, response)
}

func (br *BookRouter) HandleCachedDownloadRequest(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	req, err := api.DecodeJSON[CachedDownloadRequest](w, r)
	if err != nil {
		return
	}

	// Retrieve book data from search cache
	book, err := br.SearchCacheRepo.GetSearchResult(r.Context(), user.ID, req.SearchID, req.Index)
	if err != nil {
		applog.Error("Failed to get search result from cache:", err)
		api.WriteMessage(w, http.StatusBadRequest, "error", "search result not found or expired")
		return
	}

	// Check if user already requested this book
	hasRequested, err := br.DownloadRequestRepo.HasUserRequestedBook(r.Context(), user.ID, book.Hash)
	if err != nil {
		applog.Error("Failed to check if book already requested:", err)
		api.WriteInternalError(w)
		return
	}

	if hasRequested {
		api.WriteMessage(w, http.StatusConflict, "error", "you have already requested this book")
		return
	}

	// Check download limits
	canDownload, err := br.DownloadRequestRepo.CheckAndCreateDownload(r.Context(), user.ID, book.Hash, book.Title)
	if err != nil {
		applog.Error("Failed to check download limit:", err)
		api.WriteInternalError(w)
		return
	}

	if !canDownload {
		api.WriteMessage(w, http.StatusTooManyRequests, "error", "daily download limit reached")
		return
	}

	// Store book metadata if it doesn't exist
	existingBook, err := br.BookRepo.GetBookByHash(r.Context(), book.Hash)
	if err != nil {
		// Book doesn't exist, create it with metadata from cache
		savedBook := &model.SavedBook{
			Hash:      book.Hash,
			Title:     book.Title,
			Authors:   book.Authors,
			Publisher: book.Publisher,
			Language:  book.Language,
			Format:    book.Format,
			Size:      book.Size,
			CoverURL:  book.CoverURL,
			CoverData: book.CoverData,
			Status:    model.BookStatusProcessing,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		err = br.BookRepo.CreateBook(r.Context(), savedBook)
		if err != nil {
			applog.Error("Failed to create book record:", err)
			api.WriteInternalError(w)
			return
		}
	} else if existingBook.Status == model.BookStatusReady {
		response := DownloadResponse{
			JobID:   0,
			Status:  "ready",
			Message: "Book is already available for download",
		}
		api.WriteJSON(w, http.StatusOK, response)
		return
	}

	// Create download job
	job, err := br.DownloadJobRepo.CreateJob(r.Context(), user.ID, book.Hash)
	if err != nil {
		applog.Error("Failed to create download job:", err)
		api.WriteInternalError(w)
		return
	}

	response := DownloadResponse{
		JobID:   job.ID,
		Status:  job.Status,
		Message: "Download job created successfully",
	}

	api.WriteJSON(w, http.StatusCreated, response)
}

func (br *BookRouter) HandleJobStatus(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "jobID")
	jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid job ID")
		return
	}

	job, err := br.DownloadJobRepo.GetJobByID(r.Context(), jobID)
	if err != nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "job not found")
		return
	}

	book, err := br.BookRepo.GetBookByHash(r.Context(), job.BookHash)
	available := err == nil && book.Status == model.BookStatusReady

	response := JobStatusResponse{
		JobID:     job.ID,
		Status:    job.Status,
		Progress:  job.Progress,
		ErrorMsg:  job.ErrorMsg,
		BookHash:  job.BookHash,
		Available: available,
	}

	api.WriteJSON(w, http.StatusOK, response)
}

func (br *BookRouter) HandleDownloadFile(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		api.WriteMessage(w, http.StatusBadRequest, "error", "hash is required")
		return
	}

	// Get book information
	book, err := br.BookRepo.GetBookByHash(r.Context(), hash)
	if err != nil {
		applog.Error("Failed to get book:", err)
		api.WriteMessage(w, http.StatusNotFound, "error", "book not found")
		return
	}

	if book.Status != model.BookStatusReady || book.FilePath == "" {
		api.WriteMessage(w, http.StatusNotFound, "error", "book file not available")
		return
	}

	// Check if file actually exists on disk
	if _, err := os.Stat(book.FilePath); os.IsNotExist(err) {
		applog.Error("Book file not found on disk:", book.FilePath)
		api.WriteMessage(w, http.StatusNotFound, "error", "book file not available")
		return
	}

	// Set proper headers for file download
	filename := filepath.Base(book.FilePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Log the download for analytics
	applog.Info("User download", "book_hash", hash, "filename", filename)

	// Serve the file
	http.ServeFile(w, r, book.FilePath)
}

func (br *BookRouter) HandleUserDownloads(w http.ResponseWriter, r *http.Request) {
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

	jobs, err := br.DownloadJobRepo.GetUserJobs(r.Context(), user.ID, limit, offset)
	if err != nil {
		applog.Error("Failed to get user jobs:", err)
		api.WriteInternalError(w)
		return
	}

	// Get total count for pagination
	total, err := br.DownloadJobRepo.CountUserJobs(r.Context(), user.ID)
	if err != nil {
		applog.Error("Failed to count user jobs:", err)
		api.WriteInternalError(w)
		return
	}

	// Prevent caching to avoid 304 responses
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	response := DownloadsResponse{
		Jobs: api.EmptyIfNil(jobs),
		Pagination: Pagination{
			Limit:   limit,
			Offset:  offset,
			Total:   total,
			HasNext: offset+limit < total,
		},
	}

	api.WriteJSON(w, http.StatusOK, response)
}

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

	// Check if book exists
	_, err = br.BookRepo.GetBookByHash(r.Context(), req.BookHash)
	if err != nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "book not found")
		return
	}

	// Check if already favorited
	isFavorited, err := br.FavoriteRepo.IsFavorited(r.Context(), user.ID, req.BookHash)
	if err != nil {
		applog.Error("Failed to check favorite status:", err)
		api.WriteInternalError(w)
		return
	}

	if isFavorited {
		// Remove from favorites
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
		// Add to favorites
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

	// Get favorite book hashes
	favorites, err := br.FavoriteRepo.GetUserFavorites(r.Context(), user.ID, limit, offset)
	if err != nil {
		applog.Error("Failed to get user favorites:", err)
		api.WriteInternalError(w)
		return
	}

	// Get book details for each favorite
	var books []BookWithStats
	for _, fav := range favorites {
		book, err := br.BookRepo.GetBookByHash(r.Context(), fav.BookHash)
		if err != nil {
			continue // Skip if book not found
		}

		// Get download count
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

func (br *BookRouter) HandleDownloadStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	// Create download helper to get status
	helper := repo.NewDownloadRequestHelper(&repo.Repos{
		DownloadRequest: br.DownloadRequestRepo,
		RequestCredits:  br.RequestCreditsRepo,
	})

	// Get download status using the helper
	status, err := helper.GetDownloadStatus(r.Context(), user.ID)
	if err != nil {
		applog.Error("Failed to get download status:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteJSON(w, http.StatusOK, status)
}
