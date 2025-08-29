package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/akramboussanni/gocode/internal/anna"
	"github.com/akramboussanni/gocode/internal/model"
	"github.com/akramboussanni/gocode/internal/utils"
	"github.com/jmoiron/sqlx"
)

type SearchCacheRepo struct {
	db   *sqlx.DB
	name string
}

func NewSearchCacheRepo(db *sqlx.DB) *SearchCacheRepo {
	return &SearchCacheRepo{
		db:   db,
		name: "search_cache",
	}
}

// StoreSearchResults stores search results in cache and returns the cache ID
func (r *SearchCacheRepo) StoreSearchResults(ctx context.Context, userID int64, query string, books []*anna.Book, total int) (*model.SearchCache, error) {
	// Serialize books to JSON
	resultsJSON, err := json.Marshal(books)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search results: %w", err)
	}

	// Cache expires in 1 hour
	now := time.Now().Unix()
	expiresAt := now + 3600 // 1 hour

	searchCache := &model.SearchCache{
		ID:           utils.GenerateSnowflakeID(),
		UserID:       userID,
		Query:        query,
		Results:      string(resultsJSON),
		TotalResults: total,
		CreatedAt:    now,
		ExpiresAt:    expiresAt,
	}

	query_sql := fmt.Sprintf(`
		INSERT INTO search_cache (id, user_id, query, results, total_results, created_at, expires_at) 
		VALUES (:id, :user_id, :query, :results, :total_results, :created_at, :expires_at)
	`)

	_, err = r.db.NamedExecContext(ctx, query_sql, searchCache)
	if err != nil {
		return nil, fmt.Errorf("failed to store search cache: %w", err)
	}

	return searchCache, nil
}

// GetSearchResult retrieves a specific search result by cache ID and index
func (r *SearchCacheRepo) GetSearchResult(ctx context.Context, userID, searchID int64, index int) (*anna.Book, error) {
	var cache model.SearchCache
	query := `SELECT id, user_id, query, results, total_results, created_at, expires_at 
	          FROM search_cache 
	          WHERE id = $1 AND user_id = $2 AND expires_at > $3`

	err := r.db.GetContext(ctx, &cache, query, searchID, userID, time.Now().Unix())
	if err != nil {
		return nil, fmt.Errorf("search cache not found or expired: %w", err)
	}

	// Deserialize results
	var books []*anna.Book
	err = json.Unmarshal([]byte(cache.Results), &books)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	// Check if index is valid
	if index < 0 || index >= len(books) {
		return nil, fmt.Errorf("invalid result index: %d", index)
	}

	return books[index], nil
}

// GetSearchCache retrieves full search cache by ID
func (r *SearchCacheRepo) GetSearchCache(ctx context.Context, userID, searchID int64) (*model.SearchCache, []*anna.Book, error) {
	var cache model.SearchCache
	query := `SELECT id, user_id, query, results, total_results, created_at, expires_at 
	          FROM search_cache 
	          WHERE id = $1 AND user_id = $2 AND expires_at > $3`

	err := r.db.GetContext(ctx, &cache, query, searchID, userID, time.Now().Unix())
	if err != nil {
		return nil, nil, fmt.Errorf("search cache not found or expired: %w", err)
	}

	// Deserialize results
	var books []*anna.Book
	err = json.Unmarshal([]byte(cache.Results), &books)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return &cache, books, nil
}

// CleanupExpiredCache removes expired search cache entries
func (r *SearchCacheRepo) CleanupExpiredCache(ctx context.Context) error {
	query := "DELETE FROM search_cache WHERE expires_at < $1"
	result, err := r.db.ExecContext(ctx, query, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to cleanup expired search cache: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("Cleaned up %d expired search cache entries\n", rowsAffected)
	}

	return nil
}

// GetUserCacheCount returns the number of active cache entries for a user
func (r *SearchCacheRepo) GetUserCacheCount(ctx context.Context, userID int64) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM search_cache WHERE user_id = $1 AND expires_at > $2"
	err := r.db.GetContext(ctx, &count, query, userID, time.Now().Unix())
	return count, err
}

// DeleteUserOldestCache deletes the oldest cache entry for a user (for cache size limiting)
func (r *SearchCacheRepo) DeleteUserOldestCache(ctx context.Context, userID int64) error {
	query := `DELETE FROM search_cache 
	          WHERE id = (
	              SELECT id FROM search_cache 
	              WHERE user_id = $1 
	              ORDER BY created_at ASC 
	              LIMIT 1
	          )`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}


