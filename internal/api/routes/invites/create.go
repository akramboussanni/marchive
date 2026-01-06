package invites

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
)

// @Summary Create a new invite
// @Description Create a new invite (admins have unlimited, users may be restricted)
// @Tags Invites
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} model.InviteResponse "Invite created successfully"
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
