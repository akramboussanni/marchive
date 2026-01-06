package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/akramboussanni/marchive/internal/anna"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
)

type DownloadService struct {
	repos       *repo.Repos
	downloadDir string
	secretKey   string
}

func NewDownloadService(repos *repo.Repos, downloadDir, secretKey string) *DownloadService {
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		log.Printf("Failed to create download directory: %v", err)
	}

	return &DownloadService{
		repos:       repos,
		downloadDir: downloadDir,
		secretKey:   secretKey,
	}
}

func (ds *DownloadService) ProcessPendingDownloads(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Download service shutting down...")
			return
		default:
			// Clean up failed books older than 24 hours
			ds.cleanupFailedBooks(ctx)

			jobs, err := ds.repos.DownloadJob.GetPendingJobs(ctx, 5)
			if err != nil {
				log.Printf("Failed to get pending jobs: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			if len(jobs) == 0 {
				time.Sleep(5 * time.Second)
				continue
			}

			for _, job := range jobs {
				go ds.processJob(ctx, &job)
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func (ds *DownloadService) processJob(ctx context.Context, job *model.DownloadJob) {
	log.Printf("Processing download job %d for book %s", job.ID, job.BookHash)

	err := ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusDownloading, 10, "")
	if err != nil {
		log.Printf("Failed to update job status: %v", err)
		return
	}

	book, err := ds.repos.Book.GetBookByHash(ctx, job.BookHash)
	if err != nil {
		// Book doesn't exist, need to download it
		err = ds.processNewBook(ctx, job)
		if err != nil {
			log.Printf("Failed to process new book: %v", err)
			ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusFailed, 0, err.Error())
			return
		}
		// After successful download, mark job as completed
		ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusCompleted, 100, "")
		log.Printf("Download job %d completed successfully", job.ID)
		return
	}

	// Book already exists
	if book.Status == model.BookStatusReady && book.FilePath != "" {
		// Check if file actually exists on disk
		if _, err := os.Stat(book.FilePath); err == nil {
			// File exists, mark job as completed
			ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusCompleted, 100, "")
			ds.repos.DownloadJob.UpdateJobFilePath(ctx, job.ID, book.FilePath)
			log.Printf("Book %s already available, job %d completed", job.BookHash, job.ID)
			return
		} else {
			// File doesn't exist, need to re-download
			log.Printf("Book %s marked as ready but file missing, re-downloading", job.BookHash)
			err = ds.processNewBook(ctx, job)
			if err != nil {
				log.Printf("[RE-DOWNLOAD FAILED] Job %d, Book %s: %v", job.ID, job.BookHash, err)
				ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusFailed, 0, err.Error())
				ds.repos.Book.UpdateBookStatus(ctx, job.BookHash, model.BookStatusError, err.Error())
				return
			}
			ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusCompleted, 100, "")
			log.Printf("Re-download job %d completed successfully", job.ID)
			return
		}
	}

	// Book exists but not ready, need to download it
	err = ds.processNewBook(ctx, job)
	if err != nil {
		log.Printf("[DOWNLOAD FAILED] Job %d, Existing Book %s: %v", job.ID, job.BookHash, err)
		ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusFailed, 0, err.Error())
		ds.repos.Book.UpdateBookStatus(ctx, job.BookHash, model.BookStatusError, err.Error())
		return
	}

	ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusCompleted, 100, "")
	log.Printf("Download job %d completed successfully", job.ID)
}

func (ds *DownloadService) processNewBook(ctx context.Context, job *model.DownloadJob) error {

	existingBook, err := ds.repos.Book.GetBookByHash(ctx, job.BookHash)
	var bookMetadata *anna.Book

	if err != nil {

		log.Printf("Book %s not found, fetching metadata from Anna's archive", job.BookHash)

		bookMetadata, err = anna.GetBookMetadata(job.BookHash)
		if err != nil {
			log.Printf("Failed to get metadata from Anna's archive: %v", err)

			bookMetadata = &anna.Book{
				Hash:  job.BookHash,
				Title: "Unknown Title",
			}
		}

		book := &model.SavedBook{
			Hash:        job.BookHash,
			Title:       bookMetadata.Title,
			Authors:     bookMetadata.Authors,
			Publisher:   bookMetadata.Publisher,
			Language:    bookMetadata.Language,
			Format:      bookMetadata.Format,
			Size:        bookMetadata.Size,
			CoverURL:    bookMetadata.CoverURL,
			CoverData:   bookMetadata.CoverData,
			Status:      model.BookStatusProcessing,
			RequestedBy: &job.UserID,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		}

		err = ds.repos.Book.CreateBook(ctx, book)
		if err != nil {
			return fmt.Errorf("failed to create book record: %w", err)
		}
		existingBook = book
	} else {
		bookMetadata = &anna.Book{
			Hash:      existingBook.Hash,
			Title:     existingBook.Title,
			Authors:   existingBook.Authors,
			Publisher: existingBook.Publisher,
			Language:  existingBook.Language,
			Format:    existingBook.Format,
			Size:      existingBook.Size,
			CoverURL:  existingBook.CoverURL,
			CoverData: existingBook.CoverData,
		}
	}

	ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusDownloading, 30, "")

	annaBook := &anna.Book{
		Hash:   job.BookHash,
		Title:  bookMetadata.Title,
		Format: bookMetadata.Format,
	}

	ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusDownloading, 50, "")

	err = annaBook.Download(ds.secretKey, ds.downloadDir)
	if err != nil {
		log.Printf("[ANNA DOWNLOAD ERROR] Book %s: %v", job.BookHash, err)
		ds.repos.Book.UpdateBookStatus(ctx, job.BookHash, model.BookStatusError, "")
		return fmt.Errorf("failed to download book from Anna's Archive: %w", err)
	}

	ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusDownloading, 90, "")

	title := bookMetadata.Title
	if title == "" || title == "Unknown Title" {
		title = job.BookHash[:8]
	}
	// Sanitize filename - remove invalid characters
	title = sanitizeFilename(title)
	filename := fmt.Sprintf("%s.%s", title, bookMetadata.Format)
	filePath := filepath.Join(ds.downloadDir, filename)

	err = ds.repos.Book.UpdateBookWithMetadata(ctx, job.BookHash, model.BookStatusReady, filePath, bookMetadata)
	if err != nil {
		return fmt.Errorf("failed to update book status: %w", err)
	}

	err = ds.repos.DownloadJob.UpdateJobFilePath(ctx, job.ID, filePath)
	if err != nil {
		log.Printf("Failed to update job file path: %v", err)
	}

	return nil
}

func (ds *DownloadService) cleanupFailedBooks(ctx context.Context) {
	cutoffTime := time.Now().Add(-24 * time.Hour).Unix()

	result, err := ds.repos.Book.DeleteFailedBooks(ctx, cutoffTime)
	if err != nil {
		log.Printf("Failed to cleanup failed books: %v", err)
		return
	}

	if result > 0 {
		log.Printf("Cleaned up %d failed books", result)
	}
}

func (ds *DownloadService) StartService(ctx context.Context) {
	log.Println("Starting download service...")
	ds.ProcessPendingDownloads(ctx)
}

func sanitizeFilename(filename string) string {
	// Remove or replace invalid filename characters for Windows
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := invalidChars.ReplaceAllString(filename, "")

	// Remove control characters
	controlChars := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	sanitized = controlChars.ReplaceAllString(sanitized, "")

	// Trim spaces and dots from the end
	sanitized = strings.TrimRight(sanitized, ". ")

	// Limit length to 200 characters to avoid path too long errors
	if len(sanitized) > 200 {
		sanitized = sanitized[:200]
	}

	if sanitized == "" {
		return "untitled"
	}

	return sanitized
}
