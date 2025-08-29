package model

type JwtBlacklist struct {
	TokenID   string `db:"jti"`
	UserID    int64  `db:"user_id"`
	ExpiresAt int64  `db:"expires_at"`
}
