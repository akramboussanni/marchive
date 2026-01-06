package invites

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/api/routes/auth"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/jwt"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/akramboussanni/marchive/internal/utils"
)

// @Summary Use an invite to register
// @Description Use an invite token to create a new account and automatically sign in
// @Tags Invites
// @Accept json
// @Produce json
// @Param request body model.UseInviteRequest true "Registration details"
// @Success 200 {object} api.SuccessResponse "Account created and signed in successfully"
// @Failure 400 {object} api.ErrorResponse "Invalid invite or username taken"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /invites/use [post]
func (ir *InviteRouter) HandleUseInvite(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[model.UseInviteRequest](w, r)
	if err != nil {
		applog.Error("Failed to decode invite usage request:", err)
		return
	}

	// Validate request
	if req.Token == "" || req.Username == "" || req.Password == "" {
		api.WriteMessage(w, http.StatusBadRequest, "error", "Token, username, and password are required")
		return
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		applog.Error("Failed to hash password:", err)
		api.WriteInternalError(w)
		return
	}

	// Use invite to create user
	err = ir.InviteRepo.UseInvite(r.Context(), req.Token, req.Username, passwordHash)
	if err != nil {
		if err == repo.ErrUsernameTaken {
			api.WriteMessage(w, http.StatusBadRequest, "error", "Username already taken")
			return
		}
		applog.Error("Failed to use invite:", err)
		api.WriteMessage(w, http.StatusBadRequest, "error", "Invalid or expired invite")
		return
	}

	// Get the created user
	user, err := ir.UserRepo.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		applog.Error("Failed to get created user:", err)
		api.WriteInternalError(w)
		return
	}

	// Generate login tokens
	loginTokens := auth.GenerateLogin(jwt.CreateJwtFromUser(user))

	// Set cookies
	utils.ClearAllCookies(w)
	utils.SetSessionCookie(w, loginTokens.Session)
	utils.SetRefreshCookie(w, loginTokens.Refresh)

	applog.Info("Account created via invite", "userID:", user.ID, "username:", req.Username)
	api.WriteJSON(w, http.StatusOK, api.SuccessResponse{Message: "Account created and signed in successfully"})
}
