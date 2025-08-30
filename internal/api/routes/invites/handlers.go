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
	"github.com/go-chi/chi/v5"
)

type InviteRouter struct {
	InviteRepo *repo.InviteRepo
	UserRepo   *repo.UserRepo
}

func NewInviteRouter(inviteRepo *repo.InviteRepo, userRepo *repo.UserRepo) *InviteRouter {
	return &InviteRouter{
		InviteRepo: inviteRepo,
		UserRepo:   userRepo,
	}
}

// @Summary Create a new invite
// @Description Create a new invite using one of the user's invite tokens
// @Tags Invites
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.InviteResponse "Invite created successfully"
// @Failure 400 {object} api.ErrorResponse "No invite tokens available"
// @Failure 401 {object} api.ErrorResponse "Unauthorized"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /invites [post]
func (ir *InviteRouter) HandleCreateInvite(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	invite, err := ir.InviteRepo.CreateInvite(r.Context(), user.ID)
	if err != nil {
		if err == repo.ErrNoInviteTokens {
			api.WriteMessage(w, http.StatusBadRequest, "error", "No invite tokens available")
			return
		}
		applog.Error("Failed to create invite:", err)
		api.WriteInternalError(w)
		return
	}

	// Generate invite URL (using relative path for now)
	inviteURL := "/register?token=" + invite.Token

	response := model.InviteResponse{
		Token:     invite.Token,
		InviteURL: inviteURL,
		CreatedAt: invite.CreatedAt,
	}

	applog.Info("Invite created", "userID:", user.ID, "token:", invite.Token)
	api.WriteJSON(w, http.StatusOK, response)
}

// @Summary List user's invites
// @Description Get all invites created by the current user
// @Tags Invites
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.InviteListResponse "Invites retrieved successfully"
// @Failure 401 {object} api.ErrorResponse "Unauthorized"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /invites [get]
func (ir *InviteRouter) HandleListInvites(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	invites, err := ir.InviteRepo.GetUserInvites(r.Context(), user.ID)
	if err != nil {
		applog.Error("Failed to get user invites:", err)
		api.WriteInternalError(w)
		return
	}

	tokens, err := ir.InviteRepo.GetUserInviteTokens(r.Context(), user.ID)
	if err != nil {
		applog.Error("Failed to get user invite tokens:", err)
		api.WriteInternalError(w)
		return
	}

	response := model.InviteListResponse{
		Invites: invites,
		Tokens:  tokens,
	}

	api.WriteJSON(w, http.StatusOK, response)
}

// @Summary Revoke an invite
// @Description Revoke an unused invite and return the token to the user
// @Tags Invites
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param token path string true "Invite token to revoke"
// @Success 200 {object} api.SuccessResponse "Invite revoked successfully"
// @Failure 400 {object} api.ErrorResponse "Invalid invite or already used/revoked"
// @Failure 401 {object} api.ErrorResponse "Unauthorized"
// @Failure 500 {object} api.ErrorResponse "Internal server error"
// @Router /invites/{token}/revoke [post]
func (ir *InviteRouter) HandleRevokeInvite(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	token := chi.URLParam(r, "token")
	if token == "" {
		api.WriteMessage(w, http.StatusBadRequest, "error", "Token parameter required")
		return
	}

	err := ir.InviteRepo.RevokeInvite(r.Context(), token, user.ID)
	if err != nil {
		applog.Error("Failed to revoke invite:", err)
		api.WriteMessage(w, http.StatusBadRequest, "error", "Invalid invite or already used/revoked")
		return
	}

	applog.Info("Invite revoked", "userID:", user.ID, "token:", token)
	api.WriteJSON(w, http.StatusOK, api.SuccessResponse{Message: "Invite revoked successfully"})
}

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
