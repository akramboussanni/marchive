package routes

import (
	"net/http"
	"time"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/middleware"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
)

type RedemptionCodeRouter struct {
	RedemptionCodeRepo *repo.RedemptionCodeRepo
	UserRepo           *repo.UserRepo
}

func NewRedemptionCodeRouter(repos *repo.Repos) http.Handler {
	rcr := &RedemptionCodeRouter{
		RedemptionCodeRepo: repos.RedemptionCode,
		UserRepo:           repos.User,
	}
	r := chi.NewRouter()

	r.Use(middleware.MaxBytesMiddleware(1 << 20))

	// Rate limit: 5 attempts per hour for code redemption
	r.Group(func(r chi.Router) {
		middleware.AddRatelimit(r, 5, 1*time.Hour)
		middleware.AddAuth(r, repos.User, repos.Token)
		r.Post("/redeem", rcr.HandleRedeemCode)
	})

	return r
}

// HandleRedeemCode handles user redemption of a code
func (rcr *RedemptionCodeRouter) HandleRedeemCode(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[model.RedeemCodeRequest](w, r)
	if err != nil {
		return
	}

	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	// Get the redemption code
	code, err := rcr.RedemptionCodeRepo.GetRedemptionCodeByCode(r.Context(), req.Code)
	if err != nil {
		applog.Error("Failed to get redemption code:", err)
		api.WriteInternalError(w)
		return
	}

	if code == nil {
		api.WriteMessage(w, http.StatusNotFound, "error", "invalid or expired redemption code")
		return
	}

	// Check if code is revoked
	if code.RevokedAt != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "this redemption code has been revoked")
		return
	}

	// Check if code has expired
	if code.ExpiresAt != nil && *code.ExpiresAt < time.Now().Unix() {
		api.WriteMessage(w, http.StatusBadRequest, "error", "this redemption code has expired")
		return
	}

	// Check if code has reached max uses
	if code.CurrentUses >= code.MaxUses {
		api.WriteMessage(w, http.StatusBadRequest, "error", "this redemption code has reached its maximum uses")
		return
	}

	// Check if user has already redeemed this code
	alreadyRedeemed, err := rcr.RedemptionCodeRepo.IsCodeRedeemedByUser(r.Context(), code.ID, user.ID)
	if err != nil {
		applog.Error("Failed to check if code was redeemed:", err)
		api.WriteInternalError(w)
		return
	}

	if alreadyRedeemed {
		api.WriteMessage(w, http.StatusBadRequest, "error", "you have already redeemed this code")
		return
	}

	// Redeem the code
	err = rcr.RedemptionCodeRepo.RedeemCode(r.Context(), code.ID, user.ID, code.InviteTokens, code.RequestCredits)
	if err != nil {
		applog.Error("Failed to redeem code:", err)
		api.WriteInternalError(w)
		return
	}

	// Get updated user info to return new totals
	updatedUser, err := rcr.UserRepo.GetUserByIDSafe(r.Context(), user.ID)
	if err != nil {
		applog.Error("Failed to get updated user info:", err)
		// Don't fail the redemption, just return what we have
	}

	response := model.RedeemCodeResponse{
		Success:           true,
		Message:           "Code redeemed successfully!",
		InviteTokens:      code.InviteTokens,
		RequestCredits:    code.RequestCredits,
		NewInviteTokens:   updatedUser.InviteTokens,
		NewRequestCredits: updatedUser.RequestCredits,
	}

	api.WriteJSON(w, http.StatusOK, response)
}
