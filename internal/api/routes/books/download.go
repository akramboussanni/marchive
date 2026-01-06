package books

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
)

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

	if !canDownload && user.Role != "admin" {
		api.WriteMessage(w, http.StatusTooManyRequests, "error", "daily download limit reached")
		return
	}

	existingBook, err := br.BookRepo.GetBookByHash(r.Context(), req.Hash)
	if err == nil {
		// Book exists
		if existingBook.Status == model.BookStatusReady {
			// If book is ready and was in ghost mode, but user wants it public, disable ghost mode
			if existingBook.IsGhost && !req.IsGhost {
				err = br.BookRepo.UpdateGhostMode(r.Context(), req.Hash, false)
				if err != nil {
					applog.Error("Failed to update ghost mode:", err)
				}
			}

			err = br.BookRepo.IncrementDownloadCount(r.Context(), req.Hash)
			if err != nil {
				applog.Error("Failed to increment download count:", err)
			}

			response := DownloadResponse{
				JobID:   0,
				Status:  "ready",
				Message: "Book is already available for download",
			}
			api.WriteJSON(w, http.StatusOK, response)
			return
		} else if existingBook.IsGhost && !req.IsGhost {
			// Book exists but is processing/error in ghost mode, user wants it public
			err = br.BookRepo.UpdateGhostMode(r.Context(), req.Hash, false)
			if err != nil {
				applog.Error("Failed to update ghost mode:", err)
			}
		}

		// Set requested_by if it's not set yet
		if existingBook.RequestedBy == nil {
			requestedBy := user.ID
			err = br.BookRepo.UpdateRequestedBy(r.Context(), req.Hash, &requestedBy)
			if err != nil {
				applog.Error("Failed to update requested_by:", err)
			}
		}
	}

	if err != nil {
		requestedBy := user.ID
		newBook := &model.SavedBook{
			Hash:        req.Hash,
			Title:       req.Title,
			Authors:     req.Authors,
			Publisher:   req.Publisher,
			Language:    req.Language,
			Format:      req.Format,
			Size:        req.Size,
			CoverURL:    req.CoverURL,
			CoverData:   req.CoverData,
			Status:      model.BookStatusProcessing,
			IsGhost:     req.IsGhost,
			RequestedBy: &requestedBy,
			CreatedAt:   utils.GenerateSnowflakeID(),
			UpdatedAt:   utils.GenerateSnowflakeID(),
		}

		err = br.BookRepo.CreateBook(r.Context(), newBook)
		if err != nil {
			applog.Error("Failed to create book record:", err)
			api.WriteInternalError(w)
			return
		}
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

	if _, err := os.Stat(book.FilePath); os.IsNotExist(err) {
		applog.Error("Book file not found on disk:", book.FilePath)
		api.WriteMessage(w, http.StatusNotFound, "error", "book file not available")
		return
	}

	// Increment download count
	err = br.BookRepo.IncrementDownloadCount(r.Context(), hash)
	if err != nil {
		applog.Error("Failed to increment download count:", err)
	}

	filename := filepath.Base(book.FilePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	applog.Info("User download", "book_hash", hash, "filename", filename)

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

	total, err := br.DownloadJobRepo.CountUserJobs(r.Context(), user.ID)
	if err != nil {
		applog.Error("Failed to count user jobs:", err)
		api.WriteInternalError(w)
		return
	}

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

func (br *BookRouter) HandleDownloadStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	hash := chi.URLParam(r, "hash")
	if hash == "" {
		api.WriteMessage(w, http.StatusBadRequest, "error", "hash is required")
		return
	}

	hasRequested, err := br.DownloadRequestRepo.HasUserRequestedBook(r.Context(), user.ID, hash)
	if err != nil {
		applog.Error("Failed to check download status:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteJSON(w, http.StatusOK, map[string]bool{
		"requested": hasRequested,
	})
}
