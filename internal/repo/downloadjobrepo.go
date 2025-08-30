package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type DownloadJobRepo struct {
	Columns
	db *sqlx.DB
}

func NewDownloadJobRepo(db *sqlx.DB) *DownloadJobRepo {
	repo := &DownloadJobRepo{db: db}
	repo.Columns = ExtractColumns[model.DownloadJob]()
	return repo
}

func (r *DownloadJobRepo) CreateJob(ctx context.Context, userID int64, bookHash string) (*model.DownloadJob, error) {
	job := &model.DownloadJob{
		ID:        utils.GenerateSnowflakeID(),
		UserID:    userID,
		BookHash:  bookHash,
		Status:    model.DownloadStatusPending,
		Progress:  0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	query := fmt.Sprintf(
		"INSERT INTO downloadjobs (%s) VALUES (%s)",
		r.AllRaw,
		r.AllPrefixed,
	)
	_, err := r.db.NamedExecContext(ctx, query, job)
	return job, err
}

func (r *DownloadJobRepo) GetJobByID(ctx context.Context, jobID int64) (*model.DownloadJob, error) {
	var job model.DownloadJob
	query := fmt.Sprintf("SELECT %s FROM downloadjobs WHERE id = $1", r.AllRaw)
	err := r.db.GetContext(ctx, &job, query, jobID)
	return &job, err
}

func (r *DownloadJobRepo) GetUserJobs(ctx context.Context, userID int64, limit, offset int) ([]model.DownloadJobWithMetadata, error) {
	var jobs []model.DownloadJobWithMetadata
	query := `
		SELECT 
			dj.id, dj.user_id, dj.book_hash, dj.status, dj.progress, 
			dj.error_msg, dj.file_path, dj.created_at, dj.updated_at,
			COALESCE(sb.title, '') as title,
			COALESCE(sb.authors, '') as authors,
			COALESCE(sb.publisher, '') as publisher,
			COALESCE(sb.language, '') as language,
			COALESCE(sb.format, '') as format,
			COALESCE(sb.size, '') as size,
			COALESCE(sb.cover_url, '') as cover_url,
			COALESCE(sb.cover_data, '') as cover_data
		FROM downloadjobs dj
		LEFT JOIN savedbooks sb ON dj.book_hash = sb.hash
		WHERE dj.user_id = $1 
		ORDER BY dj.created_at DESC 
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(ctx, &jobs, query, userID, limit, offset)
	return jobs, err
}

func (r *DownloadJobRepo) UpdateJobStatus(ctx context.Context, jobID int64, status string, progress int, errorMsg string) error {
	query := `UPDATE downloadjobs SET status = $1, progress = $2, error_msg = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, status, progress, errorMsg, time.Now().Unix(), jobID)
	return err
}

func (r *DownloadJobRepo) UpdateJobFilePath(ctx context.Context, jobID int64, filePath string) error {
	query := `UPDATE downloadjobs SET file_path = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, filePath, time.Now().Unix(), jobID)
	return err
}

func (r *DownloadJobRepo) GetPendingJobs(ctx context.Context, limit int) ([]model.DownloadJob, error) {
	var jobs []model.DownloadJob
	query := fmt.Sprintf(`
		SELECT %s FROM downloadjobs 
		WHERE status = $1 
		ORDER BY created_at ASC 
		LIMIT $2
	`, r.AllRaw)
	err := r.db.SelectContext(ctx, &jobs, query, model.DownloadStatusPending, limit)
	return jobs, err
}

func (r *DownloadJobRepo) GetJobByUserAndBook(ctx context.Context, userID int64, bookHash string) (*model.DownloadJob, error) {
	var job model.DownloadJob
	query := fmt.Sprintf(`
		SELECT %s FROM downloadjobs 
		WHERE user_id = $1 AND book_hash = $2 
		ORDER BY created_at DESC 
		LIMIT 1
	`, r.AllRaw)
	err := r.db.GetContext(ctx, &job, query, userID, bookHash)
	return &job, err
}

func (r *DownloadJobRepo) CountUserJobs(ctx context.Context, userID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM downloadjobs WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}
