package auth

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/utils"
)

// @Summary Get current user profile
// @Description Retrieve the current authenticated user's profile information. Returns safe user data (excluding sensitive fields like password hash).
// @Tags Account
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} model.User "User profile information (safe fields only)"
// @Failure 401 {object} api.ErrorResponse "Unauthorized - invalid or missing session cookie"
// @Failure 429 {object} api.ErrorResponse "Rate limit exceeded (30 requests per minute)"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /auth/me [get]
func (ar *AuthRouter) HandleProfile(w http.ResponseWriter, r *http.Request) {
	applog.Info("HandleProfile called")
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		applog.Error("Failed to get user from context")
		api.WriteInternalError(w)
		return
	}

	utils.StripUnsafeFields(user)
	applog.Info("Profile retrieved", "userID:", user.ID)
	api.WriteJSON(w, 200, user)
}

// @Summary Get current user's request credits
// @Description Retrieve the current authenticated user's request credits information.
// @Tags Account
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {object} map[string]interface{} "User request credits information"
// @Failure 401 {object} api.ErrorResponse "Unauthorized - invalid or missing session cookie"
// @Failure 429 {object} api.ErrorResponse "Rate limit exceeded (30 requests per minute)"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /auth/me/credits [get]
func (ar *AuthRouter) HandleGetMyCredits(w http.ResponseWriter, r *http.Request) {
	applog.Info("HandleGetMyCredits called")
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		applog.Error("Failed to get user from context")
		api.WriteInternalError(w)
		return
	}

	// Get user's request credits
	credits, err := ar.RequestCreditsRepo.GetUserRequestCredits(r.Context(), user.ID)
	if err != nil {
		applog.Error("Failed to get user request credits:", err)
		api.WriteInternalError(w)
		return
	}

	response := map[string]interface{}{
		"user_id":         user.ID,
		"username":        user.Username,
		"request_credits": credits,
	}

	applog.Info("Credits retrieved", "userID:", user.ID, "credits:", credits)
	api.WriteJSON(w, 200, response)
}
