package invites

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
)

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
