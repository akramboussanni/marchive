package model

type UserWithStats struct {
	ID            int64  `json:"id,string"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	CreatedAt     int64  `json:"created_at,string"`
	DownloadCount int    `json:"download_count"`
	LastActive    int64  `json:"last_active,string,omitempty"`
}
