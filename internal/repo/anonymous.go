package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type AnonymousDownloadRepo struct {
	Columns
	db *sqlx.DB
}

func NewAnonymousDownloadRepo(db *sqlx.DB) *AnonymousDownloadRepo {
	repo := &AnonymousDownloadRepo{db: db}
	repo.Columns = ExtractColumns[model.AnonymousDownload]()
	return repo
}

func (r *AnonymousDownloadRepo) CreateAnonymousDownload(ctx context.Context, ipAddress, md5, title string) error {
	download := model.AnonymousDownload{
		ID:        utils.GenerateSnowflakeID(),
		IPAddress: ipAddress,
		MD5:       md5,
		Title:     title,
		CreatedAt: time.Now().Unix(),
	}

	query := fmt.Sprintf(
		"INSERT INTO anonymous_downloads (%s) VALUES (%s)",
		r.AllRaw,
		r.AllPrefixed,
	)
	_, err := r.db.NamedExecContext(ctx, query, download)
	return err
}

func (r *AnonymousDownloadRepo) GetDailyDownloadCountByIP(ctx context.Context, ipAddress string) (int, error) {
	// Convert today's date to start and end of day Unix timestamps
	startOfDay := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(), time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC).Unix()
	endOfDay := startOfDay + 86400 // 24 hours in seconds

	var count int
	query := `
		SELECT COUNT(*) 
		FROM anonymous_downloads 
		WHERE ip_address = $1 
		AND created_at >= $2 AND created_at < $3
	`

	err := r.db.GetContext(ctx, &count, query, ipAddress, startOfDay, endOfDay)
	return count, err
}

func (r *AnonymousDownloadRepo) HasIPRequestedBook(ctx context.Context, ipAddress, md5 string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM anonymous_downloads WHERE ip_address = $1 AND md5 = $2`
	err := r.db.GetContext(ctx, &count, query, ipAddress, md5)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AnonymousDownloadRepo) CheckAndCreateAnonymousDownload(ctx context.Context, ipAddress, md5, title string) (bool, error) {
	// Check if IP has already requested this book
	hasRequested, err := r.HasIPRequestedBook(ctx, ipAddress, md5)
	if err != nil {
		return false, fmt.Errorf("failed to check if book already requested: %w", err)
	}

	if hasRequested {
		return false, nil // Already requested, don't allow duplicate
	}

	// Check daily limit (hardcoded 10 for anonymous users)
	count, err := r.GetDailyDownloadCountByIP(ctx, ipAddress)
	if err != nil {
		return false, fmt.Errorf("failed to check download limit: %w", err)
	}

	if count >= 10 {
		return false, nil // Limit reached
	}

	// Create download record
	err = r.CreateAnonymousDownload(ctx, ipAddress, md5, title)
	if err != nil {
		return false, fmt.Errorf("failed to create anonymous download: %w", err)
	}

	return true, nil
}
