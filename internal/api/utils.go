package api

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteMessage(w http.ResponseWriter, status int, msgType, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{msgType: msg})
}

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var data T
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return data, err
	}
	return data, nil
}

func WriteInternalError(w http.ResponseWriter) {
	WriteMessage(w, http.StatusInternalServerError, "error", "server error")
}

func WriteInvalidCredentials(w http.ResponseWriter) {
	WriteMessage(w, http.StatusUnauthorized, "error", "invalid credentials")
}

// EmptyIfNil returns an empty slice if the input slice is nil, otherwise returns the original slice.
// This ensures JSON responses always return [] instead of null for arrays.
func EmptyIfNil[T any](slice []T) []T {
	if slice == nil {
		return []T{}
	}
	return slice
}
