package invites

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
)

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
