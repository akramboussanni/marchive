package model

// @Description Invite model for user invitations
type Invite struct {
	ID              int64   `db:"id" safe:"true" json:"id,string" example:"123456789"`
	Token           string  `db:"token" safe:"true" json:"token" example:"abc123def456"`
	InviterID       int64   `db:"inviter_id" safe:"true" json:"inviter_id,string" example:"123456789"`
	InviteeUsername *string `db:"invitee_username" safe:"true" json:"invitee_username,omitempty" example:"johndoe"`
	InviteeID       *int64  `db:"invitee_id" safe:"true" json:"invitee_id,omitempty" example:"123456789"`
	UsedAt          *int64  `db:"used_at" safe:"true" json:"used_at,omitempty" example:"1640995200"`
	RevokedAt       *int64  `db:"revoked_at" safe:"true" json:"revoked_at,omitempty" example:"1640995200"`
	CreatedAt       int64   `db:"created_at" safe:"true" json:"created_at,string" example:"1640995200"`
}

// @Description Invite creation request
type CreateInviteRequest struct {
	// No fields needed - uses current user's token
}

// @Description Invite usage request
type UseInviteRequest struct {
	Token    string `json:"token" binding:"required" example:"abc123def456"`
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
}

// @Description Invite response
type InviteResponse struct {
	Token     string `json:"token" example:"abc123def456"`
	InviteURL string `json:"invite_url" example:"https://example.com/register?token=abc123def456"`
	CreatedAt int64  `json:"created_at,string" example:"1640995200"`
}

// @Description Invite list response
type InviteListResponse struct {
	Invites []Invite `json:"invites"`
}
