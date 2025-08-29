package model

type Favorite struct {
	ID        int64  `db:"id" json:"id,string"`
	UserID    int64  `db:"user_id" json:"user_id,string"`
	BookHash  string `db:"book_hash" json:"book_hash"`
	CreatedAt int64  `db:"created_at" json:"created_at,string"`
}
