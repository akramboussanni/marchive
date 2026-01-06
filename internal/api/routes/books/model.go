package books

import (
	"github.com/akramboussanni/marchive/internal/anna"
	"github.com/akramboussanni/marchive/internal/model"
)

type SearchRequest struct {
	Query  string `json:"query" binding:"required" example:"programming golang"`
	Limit  int    `json:"limit,omitempty" example:"20"`
	Offset int    `json:"offset,omitempty" example:"0"`
}

// BookWithStatus extends anna.Book with availability status
type BookWithStatus struct {
	*anna.Book
	Status string `json:"status"` // "available" or "not_available"
}

type SearchResponse struct {
	Books      []*BookWithStatus `json:"books"`
	Total      int               `json:"total"`
	Query      string            `json:"query"`
	Pagination Pagination        `json:"pagination"`
}

type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Total   int  `json:"total"`
	HasNext bool `json:"has_next"`
}

type DownloadRequest struct {
	Hash      string `json:"hash" binding:"required" example:"abc123def456"`
	Title     string `json:"title" binding:"required" example:"Programming in Go"`
	Authors   string `json:"authors,omitempty" example:"John Doe"`
	Publisher string `json:"publisher,omitempty" example:"O'Reilly"`
	Language  string `json:"language,omitempty" example:"English"`
	Format    string `json:"format,omitempty" example:"PDF"`
	Size      string `json:"size,omitempty" example:"5.2 MB"`
	CoverURL  string `json:"cover_url,omitempty" example:"https://example.com/cover.jpg"`
	CoverData string `json:"cover_data,omitempty"`
}

type DownloadResponse struct {
	JobID   int64  `json:"job_id,string"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type BookListResponse struct {
	Books      []BookWithStats `json:"books"`
	Pagination Pagination      `json:"pagination"`
}

type BookWithStats struct {
	Hash          string `json:"hash"`
	Title         string `json:"title"`
	Authors       string `json:"authors"`
	Publisher     string `json:"publisher"`
	Language      string `json:"language"`
	Format        string `json:"format"`
	Size          string `json:"size"`
	CoverURL      string `json:"cover_url"`
	CoverData     string `json:"cover_data"`
	Status        string `json:"status"`
	DownloadCount int    `json:"download_count"`
	CreatedAt     int64  `json:"created_at,string"`
}

type JobStatusResponse struct {
	JobID     int64  `json:"job_id,string"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	ErrorMsg  string `json:"error_msg,omitempty"`
	BookHash  string `json:"book_hash"`
	Available bool   `json:"available"`
}

type ToggleFavoriteRequest struct {
	BookHash string `json:"book_hash" binding:"required" example:"abc123def456"`
}

type ToggleFavoriteResponse struct {
	IsFavorited bool   `json:"is_favorited"`
	Message     string `json:"message"`
}

type FavoritesResponse struct {
	Books      []BookWithStats `json:"books"`
	Pagination Pagination      `json:"pagination"`
}

type DownloadsResponse struct {
	Jobs       []model.DownloadJobWithMetadata `json:"jobs"`
	Pagination Pagination                      `json:"pagination"`
}
