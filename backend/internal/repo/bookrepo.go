package repo

import (
	"context"
	"fmt"
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
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`, r.AllRaw)
	err := r.db.SelectContext(ctx, &books, query, limit, offset)
	return books, err
}

func (r *BookRepo) SearchBooks(ctx context.Context, searchQuery string, limit, offset int) ([]model.SavedBook, error) {
	var books []model.SavedBook
	query := fmt.Sprintf(`
		SELECT %s FROM savedbooks 
		WHERE title ILIKE $1 OR authors ILIKE $1 OR publisher ILIKE $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`, r.AllRaw)
	searchPattern := "%" + searchQuery + "%"
	err := r.db.SelectContext(ctx, &books, query, searchPattern, limit, offset)
	return books, err
}

func (r *BookRepo) CountBooks(ctx context.Context) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM savedbooks")
	return count, err
}

func (r *BookRepo) CountSearchBooks(ctx context.Context, searchQuery string) (int, error) {
	var count int
	searchPattern := "%" + searchQuery + "%"
	err := r.db.GetContext(ctx, &count,
		"SELECT COUNT(*) FROM savedbooks WHERE title ILIKE $1 OR authors ILIKE $1 OR publisher ILIKE $1",
		searchPattern)
	return count, err
}

func (r *BookRepo) GetBooksWithDownloadCount(ctx context.Context, limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT sb.*, COALESCE(dl.download_count, 0) as download_count
		FROM savedbooks sb
		LEFT JOIN (
			SELECT md5, COUNT(*) as download_count
			FROM downloadrequests
			GROUP BY md5
		) dl ON sb.hash = dl.md5
		ORDER BY sb.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []map[string]interface{}
	for rows.Next() {
		var book model.SavedBook
		var downloadCount int

		err := rows.Scan(
			&book.ID, &book.Hash, &book.Title, &book.Authors, &book.Publisher,
			&book.Language, &book.Format, &book.Size, &book.CoverURL, &book.CoverData,
			&book.FilePath, &book.Status, &book.CreatedAt, &book.UpdatedAt, &downloadCount,
		)
		if err != nil {
			return nil, err
		}

		bookMap := map[string]interface{}{
			"hash":           book.Hash,
			"title":          book.Title,
			"authors":        book.Authors,
			"publisher":      book.Publisher,
			"language":       book.Language,
			"format":         book.Format,
			"size":           book.Size,
			"cover_url":      book.CoverURL,
			"cover_data":     book.CoverData,
			"status":         book.Status,
			"download_count": downloadCount,
			"created_at":     book.CreatedAt,
		}
		books = append(books, bookMap)
	}

	return books, nil
}
