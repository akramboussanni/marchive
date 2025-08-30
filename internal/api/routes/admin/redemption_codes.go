package admin

import (
	"net/http"
	"strconv"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/utils"
	"github.com/go-chi/chi/v5"
)

// HandleCreateRedemptionCode creates a new redemption code
func (ar *AdminRouter) HandleCreateRedemptionCode(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[model.CreateRedemptionCodeRequest](w, r)
	if err != nil {
		return
	}

	// Validate that at least one reward is specified
	if req.InviteTokens == 0 && req.RequestCredits == 0 {
		api.WriteMessage(w, http.StatusBadRequest, "error", "must specify at least one reward (invite tokens or request credits)")
		return
	}

	user, ok := utils.UserFromContext(r.Context())
	if !ok {
		api.WriteInvalidCredentials(w)
		return
	}

	code, err := ar.RedemptionCodeRepo.CreateRedemptionCode(r.Context(), &req, user.ID)
	if err != nil {
		applog.Error("Failed to create redemption code:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteJSON(w, http.StatusCreated, code)
}

// HandleListRedemptionCodes lists all redemption codes with pagination
func (ar *AdminRouter) HandleListRedemptionCodes(w http.ResponseWriter, r *http.Request) {
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

	codes, err := ar.RedemptionCodeRepo.ListRedemptionCodes(r.Context(), limit, offset)
	if err != nil {
		applog.Error("Failed to list redemption codes:", err)
		api.WriteInternalError(w)
		return
	}

	total, err := ar.RedemptionCodeRepo.CountRedemptionCodes(r.Context())
	if err != nil {
		applog.Error("Failed to count redemption codes:", err)
		api.WriteInternalError(w)
		return
	}

	response := RedemptionCodeListResponse{
		Codes: api.EmptyIfNil(codes),
		Pagination: Pagination{
			Limit:   limit,
			Offset:  offset,
			Total:   total,
			HasNext: offset+limit < total,
		},
	}

	api.WriteJSON(w, http.StatusOK, response)
}

// HandleRevokeRedemptionCode revokes a redemption code
func (ar *AdminRouter) HandleRevokeRedemptionCode(w http.ResponseWriter, r *http.Request) {
	codeIDStr := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseInt(codeIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid code ID")
		return
	}

	err = ar.RedemptionCodeRepo.RevokeRedemptionCode(r.Context(), codeID)
	if err != nil {
		applog.Error("Failed to revoke redemption code:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "redemption code revoked successfully")
}

// HandleDeleteRedemptionCode permanently deletes a redemption code
func (ar *AdminRouter) HandleDeleteRedemptionCode(w http.ResponseWriter, r *http.Request) {
	codeIDStr := chi.URLParam(r, "codeID")
	codeID, err := strconv.ParseInt(codeIDStr, 10, 64)
	if err != nil {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid code ID")
		return
	}

	err = ar.RedemptionCodeRepo.DeleteRedemptionCode(r.Context(), codeID)
	if err != nil {
		applog.Error("Failed to delete redemption code:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "redemption code deleted successfully")
}

// HandleGiveEveryoneInvite gives 1 invite token to all users
func (ar *AdminRouter) HandleGiveEveryoneInvite(w http.ResponseWriter, r *http.Request) {
	err := ar.UserRepo.GiveEveryoneInviteToken(r.Context())
	if err != nil {
		applog.Error("Failed to give everyone invite token:", err)
		api.WriteInternalError(w)
		return
	}

	api.WriteMessage(w, http.StatusOK, "success", "1 invite token given to all users successfully")
}

// RedemptionCodeListResponse represents the response for listing redemption codes
type RedemptionCodeListResponse struct {
	Codes      []*model.RedemptionCode `json:"codes"`
	Pagination Pagination               `json:"pagination"`
}
