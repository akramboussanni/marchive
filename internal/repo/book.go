package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akramboussanni/marchive/internal/anna"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/jmoiron/sqlx"
)

type BookRepo struct {
	Columns
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) *BookRepo {
	repo := &BookRepo{db: db}
	repo.Columns = ExtractColumns[model.SavedBook]()
	return repo
}

func (r *BookRepo) CreateBook(ctx context.Context, book *model.SavedBook) error {
	book.ID = utils.GenerateSnowflakeID()
	book.CreatedAt = time.Now().Unix()
	book.UpdatedAt = time.Now().Unix()

	query := fmt.Sprintf(
		"INSERT INTO savedbooks (%s) VALUES (%s)",
		r.AllRaw,
		r.AllPrefixed,
	)
	_, err := r.db.NamedExecContext(ctx, query, book)
	return err
}

func (r *BookRepo) GetBookByHash(ctx context.Context, hash string) (*model.SavedBook, error) {
	var book model.SavedBook
	query := fmt.Sprintf("SELECT %s FROM savedbooks WHERE hash = $1", r.AllRaw)
	err := r.db.GetContext(ctx, &book, query, hash)
	return &book, err
}

func (r *BookRepo) GetBookByHashForUser(ctx context.Context, hash string, userID int64, isAdmin bool) (*model.SavedBook, error) {
	var book model.SavedBook
	query := fmt.Sprintf("SELECT %s FROM savedbooks WHERE hash = $1", r.AllRaw)

	if !isAdmin {
		// Non-admin users can only see non-ghost books or their own ghost books
		query += " AND (is_ghost = false OR (is_ghost = true AND requested_by IS NOT NULL AND requested_by = $2))"
		err := r.db.GetContext(ctx, &book, query, hash, userID)
		return &book, err
	}

	err := r.db.GetContext(ctx, &book, query, hash)
	return &book, err
}

func (r *BookRepo) IncrementDownloadCount(ctx context.Context, hash string) error {
	query := `UPDATE savedbooks SET download_count = download_count + 1, updated_at = $1 WHERE hash = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now().Unix(), hash)
	return err
}

func (r *BookRepo) UpdateRequestedBy(ctx context.Context, hash string, userID *int64) error {
	query := `UPDATE savedbooks SET requested_by = $1, updated_at = $2 WHERE hash = $3`
	_, err := r.db.ExecContext(ctx, query, userID, time.Now().Unix(), hash)
	return err
}

func (r *BookRepo) UpdateBookStatus(ctx context.Context, hash, status, filePath string) error {
	query := `UPDATE savedbooks SET status = $1, file_path = $2, updated_at = $3 WHERE hash = $4`
	_, err := r.db.ExecContext(ctx, query, status, filePath, time.Now().Unix(), hash)
	return err
}

func (r *BookRepo) UpdateBookWithMetadata(ctx context.Context, hash, status, filePath string, book *anna.Book) error {
	query := `UPDATE savedbooks SET 
		title = $1, authors = $2, publisher = $3, language = $4, format = $5, size = $6,
		cover_url = $7, cover_data = $8, status = $9, file_path = $10, updated_at = $11 
		WHERE hash = $12`
	_, err := r.db.ExecContext(ctx, query,
		book.Title, book.Authors, book.Publisher, book.Language, book.Format, book.Size,
		book.CoverURL, book.CoverData, status, filePath, time.Now().Unix(), hash)
	return err
}

func (r *BookRepo) GetBooks(ctx context.Context, limit, offset int) ([]model.SavedBook, error) {
	var books []model.SavedBook
	query := fmt.Sprintf(`
		SELECT %s FROM savedbooks 
		WHERE is_ghost = false
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`, r.AllRaw)
	err := r.db.SelectContext(ctx, &books, query, limit, offset)
	return books, err
}

func (r *BookRepo) GetBooksForUser(ctx context.Context, userID int64, isAdmin bool, limit, offset int) ([]model.SavedBook, error) {
	var books []model.SavedBook
	var query string

	if isAdmin {
		// Admins see all books
		query = fmt.Sprintf(`
			SELECT %s FROM savedbooks 
			ORDER BY created_at DESC 
			LIMIT $1 OFFSET $2
		`, r.AllRaw)
		err := r.db.SelectContext(ctx, &books, query, limit, offset)
		return books, err
	}

	// Regular users see only non-ghost books or their own ghost books
	query = fmt.Sprintf(`
		SELECT %s FROM savedbooks 
		WHERE (is_ghost = false OR (is_ghost = true AND requested_by IS NOT NULL AND requested_by = $1))
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`, r.AllRaw)
	err := r.db.SelectContext(ctx, &books, query, userID, limit, offset)
	return books, err
}

func (r *BookRepo) SearchBooks(ctx context.Context, searchQuery string, limit, offset int) ([]model.SavedBook, error) {
	var books []model.SavedBook
	query := fmt.Sprintf(`
		SELECT %s FROM savedbooks 
		WHERE (LOWER(title) LIKE LOWER($1) OR LOWER(authors) LIKE LOWER($1) OR LOWER(publisher) LIKE LOWER($1))
		AND is_ghost = false
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`, r.AllRaw)
	searchPattern := "%" + searchQuery + "%"
	err := r.db.SelectContext(ctx, &books, query, searchPattern, limit, offset)
	return books, err
}

func (r *BookRepo) SearchBooksForUser(ctx context.Context, userID int64, isAdmin bool, searchQuery string, limit, offset int) ([]model.SavedBook, error) {
	var books []model.SavedBook
	searchPattern := "%" + searchQuery + "%"

	if isAdmin {
		// Admins see all books
		query := fmt.Sprintf(`
			SELECT %s FROM savedbooks 
			WHERE LOWER(title) LIKE LOWER($1) OR LOWER(authors) LIKE LOWER($1) OR LOWER(publisher) LIKE LOWER($1)
			ORDER BY created_at DESC 
			LIMIT $2 OFFSET $3
		`, r.AllRaw)
		err := r.db.SelectContext(ctx, &books, query, searchPattern, limit, offset)
		return books, err
	}

	// Regular users see only non-ghost books or their own ghost books
	query := fmt.Sprintf(`
		SELECT %s FROM savedbooks 
		WHERE (LOWER(title) LIKE LOWER($1) OR LOWER(authors) LIKE LOWER($1) OR LOWER(publisher) LIKE LOWER($1))
		AND (is_ghost = false OR (is_ghost = true AND requested_by IS NOT NULL AND requested_by = $2))
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4
	`, r.AllRaw)
	err := r.db.SelectContext(ctx, &books, query, searchPattern, userID, limit, offset)
	return books, err
}

func (r *BookRepo) CountBooks(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM savedbooks WHERE is_ghost = false")
	return count, err
}

func (r *BookRepo) CountBooksForUser(ctx context.Context, userID int64, isAdmin bool) (int, error) {
	var count int

	if isAdmin {
		err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM savedbooks")
		return count, err
	}

	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM savedbooks WHERE (is_ghost = false OR (is_ghost = true AND requested_by IS NOT NULL AND requested_by = $1))", userID)
	return count, err
}

func (r *BookRepo) CountSearchBooks(ctx context.Context, searchQuery string) (int, error) {
	var count int
	searchPattern := "%" + searchQuery + "%"
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM savedbooks WHERE LOWER(title) LIKE LOWER($1) OR LOWER(authors) LIKE LOWER($1) OR LOWER(publisher) LIKE LOWER($1)",
		searchPattern)
	return count, err
}

func (r *BookRepo) GetBooksWithDownloadCount(ctx context.Context, limit, offset int) ([]model.SavedBook, error) {
	// Select all columns from savedbooks, but override download_count with the actual count from downloadrequests
	query := `
		SELECT 
			sb.id, sb.hash, sb.title, sb.authors, sb.publisher, sb.language, 
			sb.format, sb.size, sb.cover_url, sb.cover_data, sb.file_path, 
			sb.status, COALESCE(dl.download_count, 0) as download_count, 
			sb.created_at, sb.updated_at
		FROM savedbooks sb
		LEFT JOIN (
			SELECT md5, COUNT(*) as download_count
			FROM downloadrequests
			GROUP BY md5
		) dl ON sb.hash = dl.md5
		ORDER BY sb.created_at DESC
		LIMIT $1 OFFSET $2
	`

	var books []model.SavedBook
	err := r.db.SelectContext(ctx, &books, query, limit, offset)
	return books, err
}

func (r *BookRepo) GetBooksAvailabilityByHashes(ctx context.Context, hashes []string) ([]struct {
	Hash     string `db:"hash"`
	Status   string `db:"status"`
	FilePath string `db:"file_path"`
}, error) {
	if len(hashes) == 0 {
		return []struct {
			Hash     string `db:"hash"`
			Status   string `db:"status"`
			FilePath string `db:"file_path"`
		}{}, nil
	}

	// Build the query with proper placeholders for PostgreSQL
	placeholders := make([]string, len(hashes))
	args := make([]interface{}, len(hashes))
	for i, hash := range hashes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = hash
	}

	query := fmt.Sprintf(`SELECT hash, status, file_path FROM savedbooks WHERE hash IN (%s)`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Hash     string `db:"hash"`
		Status   string `db:"status"`
		FilePath string `db:"file_path"`
	}

	for rows.Next() {
		var result struct {
			Hash     string `db:"hash"`
			Status   string `db:"status"`
			FilePath string `db:"file_path"`
		}
		err := rows.Scan(&result.Hash, &result.Status, &result.FilePath)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (r *BookRepo) DeleteFailedBooks(ctx context.Context, cutoffTime int64) (int64, error) {
	query := `DELETE FROM savedbooks WHERE status = $1 AND created_at < $2`
	result, err := r.db.ExecContext(ctx, query, model.BookStatusError, cutoffTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *BookRepo) UpdateGhostMode(ctx context.Context, hash string, isGhost bool) error {
	query := `UPDATE savedbooks SET is_ghost = $1, updated_at = $2 WHERE hash = $3`
	_, err := r.db.ExecContext(ctx, query, isGhost, time.Now().Unix(), hash)
	return err
}

func (r *BookRepo) DeleteBook(ctx context.Context, hash string) error {
	query := `DELETE FROM savedbooks WHERE hash = $1`
	_, err := r.db.ExecContext(ctx, query, hash)
	return err
}

func (r *BookRepo) UpdateBookMetadata(ctx context.Context, hash string, title, authors, publisher string) error {
	query := `UPDATE savedbooks SET title = $1, authors = $2, publisher = $3, updated_at = $4 WHERE hash = $5`
	_, err := r.db.ExecContext(ctx, query, title, authors, publisher, time.Now().Unix(), hash)
	return err
}
