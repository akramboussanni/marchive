package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
		err = ds.processNewBook(ctx, job)
		if err != nil {
			log.Printf("Failed to process new book: %v", err)
			ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusFailed, 0, err.Error())
			return
		}
	} else if book.Status == model.BookStatusReady {
		ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusCompleted, 100, "")
		ds.repos.DownloadJob.UpdateJobFilePath(ctx, job.ID, book.FilePath)
		log.Printf("Book %s already available, job %d completed", job.BookHash, job.ID)
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
			Hash:      job.BookHash,
			Title:     bookMetadata.Title,
			Authors:   bookMetadata.Authors,
			Publisher: bookMetadata.Publisher,
			Language:  bookMetadata.Language,
			Format:    bookMetadata.Format,
			Size:      bookMetadata.Size,
			CoverURL:  bookMetadata.CoverURL,
			CoverData: bookMetadata.CoverData,
			Status:    model.BookStatusProcessing,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
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
		ds.repos.Book.UpdateBookStatus(ctx, job.BookHash, model.BookStatusError, "")
		return fmt.Errorf("failed to download book: %w", err)
	}

	ds.repos.DownloadJob.UpdateJobStatus(ctx, job.ID, model.DownloadStatusDownloading, 90, "")

	title := bookMetadata.Title
	if title == "" || title == "Unknown Title" {
		title = job.BookHash[:8]
	}
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

func (ds *DownloadService) StartService(ctx context.Context) {
	log.Println("Starting download service...")
	ds.ProcessPendingDownloads(ctx)
}
