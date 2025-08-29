package model

type SavedBook struct {
	ID        int64  `db:"id" safe:"true" json:"id,string"`
	Hash      string `db:"hash" safe:"true" json:"hash"`
	Title     string `db:"title" safe:"true" json:"title"`
	Authors   string `db:"authors" safe:"true" json:"authors"`
	Publisher string `db:"publisher" safe:"true" json:"publisher"`
	Language  string `db:"language" safe:"true" json:"language"`
	Format    string `db:"format" safe:"true" json:"format"`
	Size      string `db:"size" safe:"true" json:"size"`
	CoverURL  string `db:"cover_url" safe:"true" json:"cover_url"`
	CoverData string `db:"cover_data" safe:"true" json:"cover_data"`
	FilePath  string `db:"file_path" json:"-"`
	Status    string `db:"status" safe:"true" json:"status"`
	CreatedAt int64  `db:"created_at" safe:"true" json:"created_at,string"`
	UpdatedAt int64  `db:"updated_at" safe:"true" json:"updated_at,string"`
}

type DownloadJob struct {
	ID        int64  `db:"id" safe:"true" json:"id,string"`
	UserID    int64  `db:"user_id" safe:"true" json:"user_id,string"`
	BookHash  string `db:"book_hash" safe:"true" json:"book_hash"`
	Status    string `db:"status" safe:"true" json:"status"`
	Progress  int    `db:"progress" safe:"true" json:"progress"`
	ErrorMsg  string `db:"error_msg" safe:"true" json:"error_msg"`
	FilePath  string `db:"file_path" json:"-"`
	CreatedAt int64  `db:"created_at" safe:"true" json:"created_at,string"`
	UpdatedAt int64  `db:"updated_at" safe:"true" json:"updated_at,string"`
}

type DownloadJobWithMetadata struct {
	ID        int64  `db:"id" safe:"true" json:"id,string"`
	UserID    int64  `db:"user_id" safe:"true" json:"user_id,string"`
	BookHash  string `db:"book_hash" safe:"true" json:"book_hash"`
	Status    string `db:"status" safe:"true" json:"status"`
	Progress  int    `db:"progress" safe:"true" json:"progress"`
	ErrorMsg  string `db:"error_msg" safe:"true" json:"error_msg"`
	FilePath  string `db:"file_path" json:"-"`
	CreatedAt int64  `db:"created_at" safe:"true" json:"created_at,string"`
	UpdatedAt int64  `db:"updated_at" safe:"true" json:"updated_at,string"`
	// Book metadata
	Title     string `db:"title" safe:"true" json:"title"`
	Authors   string `db:"authors" safe:"true" json:"authors"`
	Publisher string `db:"publisher" safe:"true" json:"publisher"`
	Language  string `db:"language" safe:"true" json:"language"`
	Format    string `db:"format" safe:"true" json:"format"`
	Size      string `db:"size" safe:"true" json:"size"`
	CoverURL  string `db:"cover_url" safe:"true" json:"cover_url"`
	CoverData string `db:"cover_data" safe:"true" json:"cover_data"`
}

const (
	BookStatusProcessing = "processing"
	BookStatusReady      = "ready"
	BookStatusError      = "error"

	DownloadStatusPending     = "pending"
	DownloadStatusDownloading = "downloading"
	DownloadStatusCompleted   = "completed"
	DownloadStatusFailed      = "failed"
)
