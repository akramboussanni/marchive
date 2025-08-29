package model

import "github.com/akramboussanni/gocode/internal/anna"

// SearchCache represents a cached search result
type SearchCache struct {
	ID           int64  `db:"id" json:"id,string"`
	UserID       int64  `db:"user_id" json:"user_id,string"`
	Query        string `db:"query" json:"query"`
	Results      string `db:"results" json:"-"` // JSON encoded
	TotalResults int    `db:"total_results" json:"total_results"`
	CreatedAt    int64  `db:"created_at" json:"created_at"`
	ExpiresAt    int64  `db:"expires_at" json:"expires_at"`
}

// SearchCacheResult represents a single cached search result with its index
type SearchCacheResult struct {
	SearchID int64      `json:"search_id,string"`
	Index    int        `json:"index"`
	Book     *anna.Book `json:"book"`
}

// CachedSearchResponse includes the search ID with results
type CachedSearchResponse struct {
	SearchID   int64        `json:"search_id,string"`
	Books      []*anna.Book `json:"books"`
	Total      int          `json:"total"`
	Query      string       `json:"query"`
	Pagination Pagination   `json:"pagination"`
	ExpiresAt  int64        `json:"expires_at"`
}

// Pagination represents pagination info
type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Total   int  `json:"total"`
	HasNext bool `json:"has_next"`
}

