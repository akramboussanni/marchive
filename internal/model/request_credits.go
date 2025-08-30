package model

// RequestCreditsLog tracks changes to user request credits
type RequestCreditsLog struct {
	ID          int64  `db:"id" safe:"true" json:"id,string" example:"123456789"`
	UserID      int64  `db:"user_id" safe:"true" json:"user_id,string" example:"123456789"`
	Action      string `db:"action" safe:"true" json:"action" example:"granted"`
	Amount      int    `db:"amount" safe:"true" json:"amount" example:"5"`
	Reason      string `db:"reason" safe:"true" json:"reason" example:"Bonus for active user"`
	AdminUserID *int64 `db:"admin_user_id" safe:"true" json:"admin_user_id,omitempty" example:"123456789"`
	CreatedAt   int64  `db:"created_at" safe:"true" json:"created_at,string" example:"1640995200"`
}

// RequestCreditAction constants
const (
	RequestCreditActionGranted = "granted"
	RequestCreditActionUsed    = "used"
	RequestCreditActionExpired = "expired"
)

// RequestCreditsUpdate represents a request to update user credits
type RequestCreditsUpdate struct {
	UserID int64  `json:"user_id" binding:"required" example:"123456789"`
	Amount int    `json:"amount" binding:"required" example:"5"`
	Reason string `json:"reason" example:"Bonus for active user"`
}

// RequestCreditsResponse represents the response for request credits operations
type RequestCreditsResponse struct {
	UserID         int64  `json:"user_id" example:"123456789"`
	RequestCredits int    `json:"request_credits" example:"5"`
	Message        string `json:"message" example:"Credits updated successfully"`
}
