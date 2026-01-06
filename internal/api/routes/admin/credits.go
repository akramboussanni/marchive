package admin

import (
	"fmt"
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/utils"
)

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
