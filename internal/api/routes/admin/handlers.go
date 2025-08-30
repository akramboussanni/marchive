package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
)

func (ar *AdminRouter) HandleSystemStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalUsers, err := ar.UserRepo.CountUsers(ctx)
	if err != nil {
		applog.Error("Failed to count users:", err)
		api.WriteInternalError(w)
		return
	}

	totalBooks, err := ar.BookRepo.CountBooks(ctx)
	if err != nil {
		applog.Error("Failed to count books:", err)
		api.WriteInternalError(w)
		return
	}

	totalDownloads, err := ar.DownloadRequestRepo.CountAllDownloads(ctx)
	if err != nil {
		applog.Error("Failed to count downloads:", err)
		api.WriteInternalError(w)
		return
	}

	activeUsers, err := ar.DownloadRequestRepo.CountActiveUsers(ctx, 24*time.Hour)
	if err != nil {
		applog.Error("Failed to count active users:", err)
		api.WriteInternalError(w)
		return
	}

	recentDownloads, err := ar.DownloadRequestRepo.GetRecentDownloads(ctx, 10)
	if err != nil {
		applog.Error("Failed to get recent downloads:", err)
		api.WriteInternalError(w)
		return
	}

	topBooks, err := ar.DownloadRequestRepo.GetTopDownloadedBooks(ctx, 10)
	if err != nil {
		applog.Error("Failed to get top books:", err)
		api.WriteInternalError(w)
		return
	}

	topBooksStats := make([]BookDownloadStats, 0, len(topBooks))
	for _, book := range topBooks {
		topBooksStats = append(topBooksStats, BookDownloadStats{
			Hash:          book["hash"].(string),
			Title:         book["title"].(string),
			Authors:       book["authors"].(string),
			DownloadCount: book["download_count"].(int),
		})
	}

	response := SystemStatsResponse{
		TotalUsers:      totalUsers,
		TotalBooks:      totalBooks,
		TotalDownloads:  totalDownloads,
		ActiveUsers:     activeUsers,
		RecentDownloads: api.EmptyIfNil(recentDownloads),
		TopBooks:        api.EmptyIfNil(topBooksStats),
	}

	api.WriteJSON(w, http.StatusOK, response)
}

func (ar *AdminRouter) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	users, err := ar.UserRepo.GetUsersWithStats(r.Context(), limit, offset)
	if err != nil {
		applog.Error("Failed to get users:", err)
		api.WriteInternalError(w)
		return
	}

	total, err := ar.UserRepo.CountUsers(r.Context())
	if err != nil {
		applog.Error("Failed to count users:", err)
		api.WriteInternalError(w)
		return
	}

	response := UserListResponse{
		Users: api.EmptyIfNil(users),
		Pagination: Pagination{
			Limit:   limit,
			Offset:  offset,
			Total:   total,
			HasNext: offset+limit < total,
		},
	}

	api.WriteJSON(w, http.StatusOK, response)
}

func (ar *AdminRouter) HandleSearchUsers(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[UserSearchRequest](w, r)
	if err != nil {
		return
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	users, err := ar.UserRepo.SearchUsers(r.Context(), req.Query, req.Role, req.Limit, req.Offset)
	if err != nil {
		applog.Error("Failed to search users:", err)
		api.WriteInternalError(w)
		return
	}

	total, err := ar.UserRepo.CountSearchUsers(r.Context(), req.Query, req.Role)
	if err != nil {
		applog.Error("Failed to count search users:", err)
		api.WriteInternalError(w)
		return
	}

	response := UserListResponse{
		Users: api.EmptyIfNil(users),
		Pagination: Pagination{
			Limit:   req.Limit,
			Offset:  req.Offset,
			Total:   total,
			HasNext: req.Offset+req.Limit < total,
		},
	}

	api.WriteJSON(w, http.StatusOK, response)
}

func (ar *AdminRouter) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[CreateUserRequest](w, r)
	if err != nil {
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	exists, err := ar.UserRepo.DuplicateName(r.Context(), req.Username)
	if err != nil {
		applog.Error("Failed to check username:", err)
		api.WriteInternalError(w)
		return
	}
	if exists {
		api.WriteMessage(w, http.StatusConflict, "error", "username already exists")
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		applog.Error("Failed to hash password:", err)
		api.WriteInternalError(w)
		return
	}

	user := &model.User{
		ID:           utils.GenerateSnowflakeID(),
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         req.Role,
		CreatedAt:    time.Now().Unix(),
		JwtSessionID: utils.GenerateSnowflakeID(),
	}

	err = ar.UserRepo.CreateUser(r.Context(), user)
	if err != nil {
		applog.Error("Failed to create user:", err)
		api.WriteInternalError(w)
		return
	}

	utils.StripUnsafeFields(user)
	api.WriteJSON(w, http.StatusCreated, user)
}

func (ar *AdminRouter) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid user ID")
		return
	}

	user, err := ar.UserRepo.GetUserByIDSafe(r.Context(), userID)
	if err != nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "user not found")
		return
	}

	api.WriteJSON(w, http.StatusOK, user)
}

func (ar *AdminRouter) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid user ID")
		return
	}

	req, err := api.DecodeJSON[UpdateUserRequest](w, r)
	if err != nil {
		return
	}

	user, err := ar.UserRepo.GetUserByID(r.Context(), userID)
	if err != nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "user not found")
		return
	}

	if req.Username != nil {
		exists, err := ar.UserRepo.DuplicateName(r.Context(), *req.Username)
		if err != nil {
			applog.Error("Failed to check username:", err)
			api.WriteInternalError(w)
			return
		}
		if exists && user.Username != *req.Username {
			api.WriteMessage(w, http.StatusConflict, "error", "username already exists")
			return
		}
		user.Username = *req.Username
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	err = ar.UserRepo.UpdateUser(r.Context(), user)
	if err != nil {
		applog.Error("Failed to update user:", err)
		api.WriteInternalError(w)
		return
	}

	utils.StripUnsafeFields(user)
	api.WriteJSON(w, http.StatusOK, user)
}

func (ar *AdminRouter) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid user ID")
		return
	}

	err = ar.UserRepo.DeleteUser(r.Context(), userID)
	if err != nil {
		applog.Error("Failed to delete user:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "user deleted successfully")
}

func (ar *AdminRouter) HandleChangeUserPassword(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid user ID")
		return
	}

	req, err := api.DecodeJSON[ChangeUserPasswordRequest](w, r)
	if err != nil {
		return
	}

	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		applog.Error("Failed to hash password:", err)
		api.WriteInternalError(w)
		return
	}

	err = ar.UserRepo.ChangeUserPassword(r.Context(), passwordHash, userID)
	if err != nil {
		applog.Error("Failed to change password:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "password changed successfully")
}

func (ar *AdminRouter) HandleInvalidateUserSessions(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid user ID")
		return
	}

	newSessionID := utils.GenerateSnowflakeID()
	err = ar.UserRepo.ChangeJwtSessionID(r.Context(), userID, newSessionID)
	if err != nil {
		applog.Error("Failed to invalidate sessions:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "user sessions invalidated")
}

func (ar *AdminRouter) HandleGrantRequestCredits(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[RequestCreditsUpdate](w, r)
	if err != nil {
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		api.WriteMessage(w, http.StatusBadRequest, "error", "amount must be positive")
		return
	}

	// Check if user exists
	user, err := ar.UserRepo.GetUserByIDSafe(r.Context(), req.UserID)
	if err != nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "user not found")
		return
	}

	// Get admin user from context
	adminUser, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	// Grant credits
	err = ar.RequestCreditsRepo.GrantCredits(r.Context(), req.UserID, req.Amount, req.Reason, adminUser.ID)
	if err != nil {
		applog.Error("Failed to grant request credits:", err)
		api.WriteInternalError(w)
		return
	}

	// Get updated credits
	updatedCredits, err := ar.RequestCreditsRepo.GetUserRequestCredits(r.Context(), req.UserID)
	if err != nil {
		applog.Error("Failed to get updated credits:", err)
		api.WriteInternalError(w)
		return
	}

	response := RequestCreditsResponse{
		UserID:         req.UserID,
		RequestCredits: updatedCredits,
		Message:        fmt.Sprintf("Successfully granted %d request credits to %s", req.Amount, user.Username),
	}

	api.WriteJSON(w, http.StatusOK, response)
}

func (ar *AdminRouter) HandleListUsersWithCredits(w http.ResponseWriter, r *http.Request) {
	users, err := ar.UserRepo.GetUsersWithStats(r.Context(), 1000, 0) // Get all users
	if err != nil {
		applog.Error("Failed to get users:", err)
		api.WriteInternalError(w)
		return
	}

	// Get credits for all users
	allCredits, err := ar.RequestCreditsRepo.GetAllUsersCredits(r.Context())
	if err != nil {
		applog.Error("Failed to get users credits:", err)
		api.WriteInternalError(w)
		return
	}

	// Build response
	usersWithCredits := make([]UserWithCredits, 0, len(users))
	for _, user := range users {
		credits := allCredits[user.ID]
		usersWithCredits = append(usersWithCredits, UserWithCredits{
			ID:             user.ID,
			Username:       user.Username,
			Role:           user.Role,
			RequestCredits: credits,
		})
	}

	response := UserCreditsListResponse{
		Users: usersWithCredits,
	}

	api.WriteJSON(w, http.StatusOK, response)
}

func (ar *AdminRouter) HandleGetUserCredits(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid user ID")
		return
	}

	// Check if user exists
	user, err := ar.UserRepo.GetUserByIDSafe(r.Context(), userID)
	if err != nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "user not found")
		return
	}

	// Get user credits
	credits, err := ar.RequestCreditsRepo.GetUserRequestCredits(r.Context(), userID)
	if err != nil {
		applog.Error("Failed to get user credits:", err)
		api.WriteInternalError(w)
		return
	}

	response := UserWithCredits{
		ID:             user.ID,
		Username:       user.Username,
		Role:           user.Role,
		RequestCredits: credits,
	}

	api.WriteJSON(w, http.StatusOK, response)
}
