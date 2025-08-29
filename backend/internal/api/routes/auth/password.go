// this file contains translations
package auth

import (
	"context"
	"net/http"

	"github.com/akramboussanni/gocode/internal/api"
	"github.com/akramboussanni/gocode/internal/applog"
	"github.com/akramboussanni/gocode/internal/model"
	"github.com/akramboussanni/gocode/internal/utils"
)

// shared helper for password change logic
func (ar *AuthRouter) changeUserPassword(ctx context.Context, w http.ResponseWriter, user *model.User, newPassword, ip string) bool {
	if !utils.IsValidPassword(newPassword) {
		applog.Warn("Invalid new password format", "userID:", user.ID)
		api.WriteMessage(w, 400, "error", "invalid password")
		return false
	}
	if utils.ComparePassword(user.PasswordHash, newPassword) {
		applog.Error("Same password")
		api.WriteMessage(w, 400, "error", "same password")
		return false
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		applog.Error("Failed to hash new password:", err)
		api.WriteInternalError(w)
		return false
	}
	if err := ar.UserRepo.ChangeUserPassword(ctx, hash, user.ID); err != nil {
		applog.Error("Failed to change user password:", err)
		api.WriteInternalError(w)
		return false
	}
	if err := ar.UserRepo.ChangeJwtSessionID(ctx, user.ID, utils.GenerateSnowflakeID()); err != nil {
		applog.Error("Failed to revoke all sessions:", err)
		api.WriteInternalError(w)
		return false
	}
	if err := ar.LockoutRepo.UnlockAccount(ctx, user.ID, ip); err != nil {
		applog.Error("Failed to revoke all sessions:", err)
		api.WriteInternalError(w)
		return false
	}
	applog.Info("Password changed successfully", "userID:", user.ID)
	return true
}

// @Summary Change password (authenticated)
// @Description Change user password while authenticated. Requires current password verification and new password must meet security requirements.
// @Tags Password Management
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param request body PasswordChangeRequest true "Current password and new password"
// @Success 200 {string} string "Password changed successfully"
// @Failure 400 {object} api.ErrorResponse "Invalid password format or requirements not met"
// @Failure 401 {object} api.ErrorResponse "Unauthorized or incorrect current password"
// @Failure 429 {object} api.ErrorResponse "Rate limit exceeded (5 requests per minute)"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /auth/change-password [post]
func (ar *AuthRouter) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	applog.Info("HandleChangePassword called")
	req, err := api.DecodeJSON[PasswordChangeRequest](w, r)
	if err != nil {
		applog.Error("Failed to decode change password request:", err)
		return
	}

	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		applog.Error("Failed to get user from context")
		return
	}

	if !utils.ComparePassword(user.PasswordHash, req.OldPassword) {
		applog.Warn("Incorrect current password", "userID:", user.ID)
		api.WriteInvalidCredentials(w)
		return
	}

	if !ar.changeUserPassword(r.Context(), w, user, req.NewPassword, utils.GetClientIP(r)) {
		return
	}

	w.WriteHeader(http.StatusOK)
}
