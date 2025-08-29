package model

// DownloadRequest tracks individual download requests by users
type DownloadRequest struct {
	ID        int64  `db:"id" safe:"true" json:"id,string" example:"123456789"`
	UserID    int64  `db:"user_id" safe:"true" json:"user_id,string" example:"123456789"`
	MD5       string `db:"md5" safe:"true" json:"md5" example:"d41d8cd98f00b204e9800998ecf8427e"`
	Title     string `db:"title" safe:"true" json:"title" example:"example_file.txt"`
	CreatedAt int64  `db:"created_at" safe:"true" json:"created_at,string" example:"1640995200"`
}
