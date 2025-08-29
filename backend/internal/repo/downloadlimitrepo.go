package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type DownloadRequestRepo struct {
	Columns
	db *sqlx.DB
}

func NewDownloadRequestRepo(db *sqlx.DB) *DownloadRequestRepo {
	repo := &DownloadRequestRepo{db: db}
	repo.Columns = ExtractColumns[model.DownloadRequest]()
	return repo
}

func (r *DownloadRequestRepo) CreateDownloadRequest(ctx context.Context, userID int64, md5, title string) error {
	request := model.DownloadRequest{
		ID:        utils.GenerateSnowflakeID(),
		UserID:    userID,
		MD5:       md5,
		Title:     title,
		CreatedAt: time.Now().Unix(),
	}

	query := fmt.Sprintf(
		"INSERT INTO downloadrequests (%s) VALUES (%s)",
		r.AllRaw,
		r.AllPrefixed,
	)
	_, err := r.db.NamedExecContext(ctx, query, request)
	return err
}

func (r *DownloadRequestRepo) GetDailyDownloadCount(ctx context.Context, userID int64) (int, error) {
	today := time.Now().UTC().Format("2006-01-02")

	var count int
	query := `
		SELECT COUNT(*) 
		FROM downloadrequests 
		WHERE user_id = $1 
		AND DATE(created_at, 'unixepoch') = $2
	`

	err := r.db.GetContext(ctx, &count, query, userID, today)
	return count, err
}

func (r *DownloadRequestRepo) CanDownload(ctx context.Context, userID int64) (bool, error) {
	count, err := r.GetDailyDownloadCount(ctx, userID)
	if err != nil {
		return false, err
	}

	return count < 10, nil
}

func (r *DownloadRequestRepo) GetRemainingDownloads(ctx context.Context, userID int64) (int, error) {
	count, err := r.GetDailyDownloadCount(ctx, userID)
	if err != nil {
		return 0, err
	}

	remaining := 10 - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

func (r *DownloadRequestRepo) HasUserRequestedBook(ctx context.Context, userID int64, md5 string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM downloadrequests WHERE user_id = $1 AND md5 = $2`
	err := r.db.GetContext(ctx, &count, query, userID, md5)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DownloadRequestRepo) CheckAndCreateDownload(ctx context.Context, userID int64, md5, title string) (bool, error) {
	// Check if user has already requested this book
	hasRequested, err := r.HasUserRequestedBook(ctx, userID, md5)
	if err != nil {
		return false, fmt.Errorf("failed to check if book already requested: %w", err)
	}

	if hasRequested {
		return false, nil // Already requested, don't allow duplicate
	}

	canDownload, err := r.CanDownload(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check download limit: %w", err)
	}

	if !canDownload {
		return false, nil
	}

	err = r.CreateDownloadRequest(ctx, userID, md5, title)
	if err != nil {
		return false, fmt.Errorf("failed to create download request: %w", err)
	}

	return true, nil
}

func (r *DownloadRequestRepo) GetUserDownloadHistory(ctx context.Context, userID int64, limit int, offset int) ([]model.DownloadRequest, error) {
	var requests []model.DownloadRequest
	query := fmt.Sprintf(`
		SELECT %s 
		FROM downloadrequests 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`, r.AllRaw)

	err := r.db.SelectContext(ctx, &requests, query, userID, limit, offset)
	return requests, err
}

func (r *DownloadRequestRepo) GetDownloadsByMD5(ctx context.Context, md5 string) ([]model.DownloadRequest, error) {
	var requests []model.DownloadRequest
	query := fmt.Sprintf(`
		SELECT %s 
		FROM downloadrequests 
		WHERE md5 = $1 
		ORDER BY created_at DESC
	`, r.AllRaw)

	err := r.db.SelectContext(ctx, &requests, query, md5)
	return requests, err
}

func (r *DownloadRequestRepo) GetDailyDownloadStats(ctx context.Context, userID int64, days int) (map[string]int, error) {
	stats := make(map[string]int)

	query := `
		SELECT DATE(created_at, 'unixepoch') as download_date, COUNT(*) as count
		FROM downloadrequests 
		WHERE user_id = $1 
		AND created_at >= $2
		GROUP BY DATE(created_at, 'unixepoch')
		ORDER BY download_date DESC
	`

	cutoffTime := time.Now().UTC().AddDate(0, 0, -days).Unix()

	rows, err := r.db.QueryContext(ctx, query, userID, cutoffTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			return nil, err
		}
		stats[date] = count
	}

	return stats, nil
}

func (r *DownloadRequestRepo) CleanupOldRequests(ctx context.Context, daysOld int) error {
	cutoffTime := time.Now().UTC().AddDate(0, 0, -daysOld).Unix()

	query := "DELETE FROM downloadrequests WHERE created_at < $1"
	result, err := r.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		return fmt.Errorf("failed to cleanup old requests: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected > 0 {
		fmt.Printf("Cleaned up %d old download requests\n", rowsAffected)
	}

	return nil
}

func (r *DownloadRequestRepo) CountAllDownloads(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM downloadrequests")
	return count, err
}

func (r *DownloadRequestRepo) CountActiveUsers(ctx context.Context, duration time.Duration) (int, error) {
	cutoffTime := time.Now().UTC().Add(-duration).Unix()
	var count int
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(DISTINCT user_id) FROM downloadrequests WHERE created_at >= $1",
		cutoffTime)
	return count, err
}

func (r *DownloadRequestRepo) GetRecentDownloads(ctx context.Context, limit int) ([]model.DownloadRequest, error) {
	var requests []model.DownloadRequest
	query := fmt.Sprintf(`
		SELECT %s FROM downloadrequests 
		ORDER BY created_at DESC 
		LIMIT $1
	`, r.AllRaw)

	err := r.db.SelectContext(ctx, &requests, query, limit)
	return requests, err
}

func (r *DownloadRequestRepo) GetTopDownloadedBooks(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT md5, title, COUNT(*) as download_count
		FROM downloadrequests
		GROUP BY md5, title
		ORDER BY download_count DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []map[string]interface{}
	for rows.Next() {
		var hash, title string
		var count int

		err := rows.Scan(&hash, &title, &count)
		if err != nil {
			return nil, err
		}

		book := map[string]interface{}{
			"hash":           hash,
			"title":          title,
			"authors":        "",
			"download_count": count,
		}
		books = append(books, book)
	}

	return books, nil
}
