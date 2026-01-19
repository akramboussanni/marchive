package books

import (
	"net/http"
	"os"

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

	// Get book info first to delete the file
	book, err := br.BookRepo.GetBookByHash(r.Context(), req.BookHash)
	if err != nil {
		applog.Error("Failed to get book:", err)
		api.WriteInternalError(w)
		return
	}

	// Delete the physical file if it exists
	if book.FilePath != "" {
		if err := os.Remove(book.FilePath); err != nil && !os.IsNotExist(err) {
			applog.Error("Failed to delete book file:", err)
			// Continue anyway to remove from database
		}
	}

	// Delete from database
	err = br.BookRepo.DeleteBook(r.Context(), req.BookHash)
	if err != nil {
		applog.Error("Failed to delete book from database:", err)
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

func (br *BookRouter) HandleUploadBook(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	// Parse multipart form (max 500MB)
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		applog.Error("Failed to parse multipart form:", err)
		api.WriteMessage(w, http.StatusBadRequest, "error", "failed to parse form data")
		return
	}

	// Get book file
	bookFile, bookHeader, err := r.FormFile("book")
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "book file is required")
		return
	}
	defer bookFile.Close()

	// Validate book file
	if err := utils.ValidateBookFile(bookHeader); err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Get metadata
	title := r.FormValue("title")
	if title == "" {
		api.WriteMessage(w, http.StatusBadRequest, "error", "title is required")
		return
	}
	authors := r.FormValue("authors")
	publisher := r.FormValue("publisher")
	language := r.FormValue("language")

	// Get cover image (optional)
	var coverPath string
	var coverData string
	coverFile, coverHeader, err := r.FormFile("cover")
	if err == nil {
		defer coverFile.Close()
		
		// Validate cover image
		if err := utils.ValidateImageFile(coverHeader); err != nil {
			api.WriteMessage(w, http.StatusBadRequest, "error", err.Error())
			return
		}

		// Save cover image
		coverPath, err = utils.SaveUploadedFile(coverFile, coverHeader, "downloads/covers")
		if err != nil {
			applog.Error("Failed to save cover image:", err)
			api.WriteInternalError(w)
			return
		}
		coverData = coverPath // Store path as cover data for now
	}

	// Save book file
	bookPath, err := utils.SaveUploadedFile(bookFile, bookHeader, "downloads/uploads")
	if err != nil {
		applog.Error("Failed to save book file:", err)
		// Clean up cover if it was saved
		if coverPath != "" {
			utils.DeleteFile(coverPath)
		}
		api.WriteInternalError(w)
		return
	}

	// Get file size and format
	fileSize, _ := utils.GetFileSize(bookPath)
	format := utils.GetFileExtension(bookHeader.Filename)

	// Generate unique hash for the book
	hash := utils.GenerateHash(bookPath)

	// Create book record
	book := &model.SavedBook{
		Hash:             hash,
		Title:            title,
		Authors:          authors,
		Publisher:        publisher,
		Language:         language,
		Format:           format,
		Size:             utils.FormatFileSize(fileSize),
		CoverURL:         "",
		CoverData:        coverData,
		FilePath:         bookPath,
		Status:           model.BookStatusReady,
		IsUploaded:       true,
		UploadedBy:       &user.ID,
		OriginalFilename: bookHeader.Filename,
	}

	// Save to database
	if err := br.BookRepo.CreateUploadedBook(r.Context(), book); err != nil {
		applog.Error("Failed to create book record:", err)
		// Clean up files
		utils.DeleteFile(bookPath)
		if coverPath != "" {
			utils.DeleteFile(coverPath)
		}
		api.WriteInternalError(w)
		return
	}

	api.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Book uploaded successfully",
		"book": BookWithStats{
			Hash:             book.Hash,
			Title:            book.Title,
			Authors:          book.Authors,
			Publisher:        book.Publisher,
			Language:         book.Language,
			Format:           book.Format,
			Size:             book.Size,
			CoverURL:         book.CoverURL,
			CoverData:        book.CoverData,
			Status:           book.Status,
			IsUploaded:       book.IsUploaded,
			UploadedBy:       book.UploadedBy,
			OriginalFilename: book.OriginalFilename,
			CreatedAt:        book.CreatedAt,
		},
	})
}

func (br *BookRouter) HandleUpdateCover(w http.ResponseWriter, r *http.Request) {
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

	// Get the book to verify ownership
	book, err := br.BookRepo.GetBookByHash(r.Context(), hash)
	if err != nil {
		applog.Error("Failed to get book:", err)
		api.WriteMessage(w, http.StatusNotFound, "error", "book not found")
		return
	}

	// Check if user owns the book (or is admin)
	if book.UploadedBy == nil || (*book.UploadedBy != user.ID && user.Role != "admin") {
		api.WriteMessage(w, http.StatusForbidden, "error", "you don't have permission to edit this book")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		applog.Error("Failed to parse multipart form:", err)
		api.WriteMessage(w, http.StatusBadRequest, "error", "failed to parse form data")
		return
	}

	// Get cover image
	coverFile, coverHeader, err := r.FormFile("cover")
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "cover image is required")
		return
	}
	defer coverFile.Close()

	// Validate cover image
	if err := utils.ValidateImageFile(coverHeader); err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	// Delete old cover if it exists
	if book.CoverData != "" {
		utils.DeleteFile(book.CoverData)
	}

	// Save new cover image
	coverPath, err := utils.SaveUploadedFile(coverFile, coverHeader, "downloads/covers")
	if err != nil {
		applog.Error("Failed to save cover image:", err)
		api.WriteInternalError(w)
		return
	}

	// Update database
	if err := br.BookRepo.UpdateBookCover(r.Context(), hash, "", coverPath); err != nil {
		applog.Error("Failed to update book cover:", err)
		utils.DeleteFile(coverPath)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "Cover updated successfully")
}

// HandleRestoreBooks scans the downloads directory and adds any books not in the database
func (br *BookRouter) HandleRestoreBooks(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	if user.Role != "admin" {
		api.WriteMessage(w, http.StatusForbidden, "error", "admin access required")
		return
	}

	// Allowed book formats
	allowedFormats := map[string]bool{
		".pdf": true, ".epub": true, ".mobi": true, ".azw3": true,
		".djvu": true, ".fb2": true, ".txt": true,
	}

	// Directories to scan
	downloadsDirs := []string{"downloads", "downloads/uploads"}
	
	var restored int
	var skipped int
	var errors []string

	for _, dir := range downloadsDirs {
		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Read directory entries
		entries, err := os.ReadDir(dir)
		if err != nil {
			applog.Error("Failed to read directory:", err)
			errors = append(errors, "Failed to read "+dir)
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filename := entry.Name()
			ext := utils.GetFileExtensionWithDot(filename)

			// Skip non-book files
			if !allowedFormats[ext] {
				continue
			}

			// Get full path
			filePath, err := utils.GetAbsolutePath(dir, filename)
			if err != nil {
				continue
			}

			// Check if this file is already in the database by file path
			exists, err := br.BookRepo.BookExistsByFilePath(r.Context(), filePath)
			if err != nil {
				applog.Error("Failed to check book existence:", err)
				continue
			}

			if exists {
				skipped++
				continue
			}

			// Get file info
			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Guess metadata from filename
			title := utils.GuessBookTitle(filename)
			format := utils.GetFileExtension(filename)

			// Create book record
			book := &model.SavedBook{
				Hash:             utils.GenerateHash(filePath),
				Title:            title,
				Authors:          "",
				Publisher:        "",
				Language:         "",
				Format:           format,
				Size:             utils.FormatFileSize(info.Size()),
				CoverURL:         "",
				CoverData:        "",
				FilePath:         filePath,
				Status:           model.BookStatusReady,
				IsUploaded:       true,
				UploadedBy:       &user.ID,
				OriginalFilename: filename,
			}

			if err := br.BookRepo.CreateUploadedBook(r.Context(), book); err != nil {
				applog.Error("Failed to create book record:", err)
				errors = append(errors, "Failed to add: "+filename)
				continue
			}

			restored++
		}
	}

	response := map[string]interface{}{
		"restored": restored,
		"skipped":  skipped,
		"errors":   errors,
		"message":  "Restore completed",
	}

	api.WriteJSON(w, http.StatusOK, response)
}
