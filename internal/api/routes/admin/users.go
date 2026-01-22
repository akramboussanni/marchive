package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/akramboussanni/marchive/internal/services"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
)

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

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		applog.Error("Failed to hash password:", err)
		api.WriteInternalError(w)
		return
	}

	user, err := ar.UserService.CreateUser(r.Context(), services.CreateUserParams{
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         req.Role,
	})
	if err != nil {
		if err == repo.ErrUsernameTaken {
			api.WriteMessage(w, http.StatusConflict, "error", "username already exists")
			return
		}
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

func (ar *AdminRouter) HandleSetDailyLimit(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[SetDailyLimitRequest](w, r)
	if err != nil {
		return
	}

	// Validate daily limit
	if req.DailyLimit < 0 {
		api.WriteMessage(w, http.StatusBadRequest, "error", "daily limit must be non-negative")
		return
	}

	// Check if user exists
	user, err := ar.UserRepo.GetUserByIDSafe(r.Context(), req.UserID)
	if err != nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "user not found")
		return
	}

	// Update daily download limit
	err = ar.UserRepo.UpdateDailyDownloadLimit(r.Context(), req.UserID, req.DailyLimit)
	if err != nil {
		applog.Error("Failed to update daily download limit:", err)
		api.WriteInternalError(w)
		return
	}

	response := map[string]interface{}{
		"user_id":             req.UserID,
		"username":            user.Username,
		"daily_download_limit": req.DailyLimit,
		"message":             fmt.Sprintf("Successfully set daily download limit to %d for %s", req.DailyLimit, user.Username),
	}

	api.WriteJSON(w, http.StatusOK, response)
}

