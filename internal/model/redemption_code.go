package model

// @Description Redemption code model for admins to create codes that users can redeem
type RedemptionCode struct {
	ID             int64  `db:"id" safe:"true" json:"id,string" example:"123456789"`
	Code           string `db:"code" safe:"true" json:"code" example:"WELCOME2024"`
	Description    string `db:"description" safe:"true" json:"description" example:"Welcome bonus for new users"`
	InviteTokens   int    `db:"invite_tokens" safe:"true" json:"invite_tokens" example:"2"`
	RequestCredits int    `db:"request_credits" safe:"true" json:"request_credits" example:"5"`
	MaxUses        int    `db:"max_uses" safe:"true" json:"max_uses" example:"100"`
	CurrentUses    int    `db:"current_uses" safe:"true" json:"current_uses" example:"45"`
	ExpiresAt      *int64 `db:"expires_at" safe:"true" json:"expires_at,omitempty" example:"1640995200"`
	RevokedAt      *int64 `db:"revoked_at" safe:"true" json:"revoked_at,omitempty"`
	CreatedBy      int64  `db:"created_by" safe:"true" json:"created_by,string" example:"123456789"`
	CreatedAt      int64  `db:"created_at" safe:"true" json:"created_at,string" example:"1640995200"`
}

// @Description Redemption log entry tracking when a user redeemed a code
type RedemptionLog struct {
	ID                    int64 `db:"id" safe:"true" json:"id,string" example:"123456789"`
	CodeID                int64 `db:"code_id" safe:"true" json:"code_id,string" example:"123456789"`
	UserID                int64 `db:"user_id" safe:"true" json:"user_id,string" example:"123456789"`
	RedeemedAt            int64 `db:"redeemed_at" safe:"true" json:"redeemed_at,string" example:"1640995200"`
	InviteTokensGranted   int   `db:"invite_tokens_granted" safe:"true" json:"invite_tokens_granted" example:"2"`
	RequestCreditsGranted int   `db:"request_credits_granted" safe:"true" json:"request_credits_granted" example:"5"`
}

// @Description Request to create a new redemption code
type CreateRedemptionCodeRequest struct {
	Code           string `json:"code" validate:"required,min=3,max=32,alphanum"`
	Description    string `json:"description" validate:"required,min=1,max=500"`
	InviteTokens   int    `json:"invite_tokens" validate:"min=0,max=100"`
	RequestCredits int    `json:"request_credits" validate:"min=0,max=1000"`
	MaxUses        int    `json:"max_uses" validate:"min=1,max=1000000"`
	ExpiresAt      *int64 `json:"expires_at,omitempty"`
}

// @Description Request to redeem a code
type RedeemCodeRequest struct {
	Code string `json:"code" validate:"required,min=3,max=32"`
}

// @Description Response for code redemption
type RedeemCodeResponse struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	InviteTokens      int    `json:"invite_tokens_granted,omitempty"`
	RequestCredits    int    `json:"request_credits_granted,omitempty"`
	NewInviteTokens   int    `json:"new_invite_tokens_total,omitempty"`
	NewRequestCredits int    `json:"new_request_credits_total,omitempty"`
}
