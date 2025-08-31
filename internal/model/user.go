package model

// @Description User model with profile information
type User struct {
	ID             int64  `db:"id" safe:"true" json:"id,string" example:"123456789"`
	Username       string `db:"username" safe:"true" json:"username" example:"johndoe"`
	PasswordHash   string `db:"password_hash" json:"-"`
	CreatedAt      int64  `db:"created_at" safe:"true" json:"created_at,string" example:"1640995200"`
	Role           string `db:"user_role" safe:"true" json:"role" example:"user"`
	JwtSessionID   int64  `db:"jwt_session_id" json:"-"`
	RequestCredits int    `db:"request_credits" safe:"true" json:"request_credits" example:"0"`
	InviteTokens   int    `db:"invite_tokens" safe:"true" json:"invite_tokens" example:"1"`
}
