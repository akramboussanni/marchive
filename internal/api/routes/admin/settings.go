package admin

import (
	"net/http"

	"github.com/akramboussanni/marchive/internal/api"
	"github.com/akramboussanni/marchive/internal/applog"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
)

type SettingsHandler struct {
	settingsRepo *repo.SettingsRepo
}

func NewSettingsHandler(settingsRepo *repo.SettingsRepo) *SettingsHandler {
	return &SettingsHandler{settingsRepo: settingsRepo}
}

// HandleGetSettings returns all application settings
func (h *SettingsHandler) HandleGetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingsRepo.GetAllSettings(r.Context())
	if err != nil {
		applog.Error("Failed to get settings:", err)
		api.WriteMessage(w, http.StatusInternalServerError, "error", "failed to get settings")
		return
	}

	// Convert to map for easier frontend consumption
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	api.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":   "success",
		"settings": settingsMap,
	})
}

type UpdateSettingRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HandleUpdateSetting updates a single setting
func (h *SettingsHandler) HandleUpdateSetting(w http.ResponseWriter, r *http.Request) {
	req, err := api.DecodeJSON[UpdateSettingRequest](w, r)
	if err != nil {
		return
	}

	if req.Key == "" {
		api.WriteMessage(w, http.StatusBadRequest, "error", "key is required")
		return
	}

	// Validate known settings
	validSettings := map[string]bool{
		model.SettingAnonymousAccessEnabled: true,
	}

	if !validSettings[req.Key] {
		api.WriteMessage(w, http.StatusBadRequest, "error", "invalid setting key")
		return
	}

	// For boolean settings, validate value
	if req.Key == model.SettingAnonymousAccessEnabled {
		if req.Value != "true" && req.Value != "false" {
			api.WriteMessage(w, http.StatusBadRequest, "error", "value must be 'true' or 'false'")
			return
		}
	}

	err = h.settingsRepo.SetSetting(r.Context(), req.Key, req.Value)
	if err != nil {
		applog.Error("Failed to update setting:", err)
		api.WriteMessage(w, http.StatusInternalServerError, "error", "failed to update setting")
		return
	}

	api.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "setting updated",
		"key":     req.Key,
		"value":   req.Value,
	})
}
