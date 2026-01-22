package admin

import "github.com/akramboussanni/marchive/internal/model"

type UserListResponse struct {
	Users      []UserWithStats `json:"users"`
	Pagination Pagination      `json:"pagination"`
}

type UserWithStats = model.UserWithStats

type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Total   int  `json:"total"`
	HasNext bool `json:"has_next"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required" example:"newuser"`
	Password string `json:"password" binding:"required" example:"SecurePass123!"`
	Role     string `json:"role" example:"user"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" example:"updateduser"`
	Role     *string `json:"role,omitempty" example:"admin"`
}

type ChangeUserPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required" example:"NewSecurePass123!"`
}

type SystemStatsResponse struct {
	TotalUsers      int                     `json:"total_users"`
	TotalBooks      int                     `json:"total_books"`
	TotalDownloads  int                     `json:"total_downloads"`
	ActiveUsers     int                     `json:"active_users_24h"`
	RecentDownloads []model.DownloadRequest `json:"recent_downloads"`
	TopBooks        []BookDownloadStats     `json:"top_books"`
}

type BookDownloadStats struct {
	Hash          string `json:"hash"`
	Title         string `json:"title"`
	Authors       string `json:"authors"`
	DownloadCount int    `json:"download_count"`
}

type UserSearchRequest struct {
	Query  string `json:"query,omitempty" example:"john"`
	Role   string `json:"role,omitempty" example:"user"`
	Limit  int    `json:"limit,omitempty" example:"20"`
	Offset int    `json:"offset,omitempty" example:"0"`
}

type RequestCreditsUpdate struct {
	UserID int64  `json:"user_id" binding:"required" example:"123456789"`
	Amount int    `json:"amount" binding:"required" example:"5"`
	Reason string `json:"reason" example:"Bonus for active user"`
}

type RequestCreditsResponse struct {
	UserID         int64  `json:"user_id" example:"123456789"`
	RequestCredits int    `json:"request_credits" example:"5"`
	Message        string `json:"message" example:"Credits updated successfully"`
}

type SetDailyLimitRequest struct {
	UserID     int64 `json:"user_id" binding:"required" example:"123456789"`
	DailyLimit int   `json:"daily_limit" binding:"required" example:"10"`
}

