package model

// @Description Anonymous download tracking for unauthenticated users
type AnonymousDownload struct {
	ID        int64  `db:"id" json:"id,string"`
	IPAddress string `db:"ip_address" json:"ip_address"`
	MD5       string `db:"md5" json:"md5"`
	Title     string `db:"title" json:"title"`
	CreatedAt int64  `db:"created_at" json:"created_at,string"`
}
