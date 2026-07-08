package feature_flags

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type FeatureFlagHandler struct {
	repository FeatureFlagRepository
}

func NewFeatureFlagHandler(repo FeatureFlagRepository) *FeatureFlagHandler {
	return &FeatureFlagHandler{
		repository: repo,
	}
}

func (h FeatureFlagHandler) CreateFlag(w http.ResponseWriter, r *http.Request) {
	var createFeatureFlag CreateFeatureFlag

	if err := json.NewDecoder(r.Body).Decode(&createFeatureFlag); err != nil {
		slog.Warn("Failed to decode create flag request", "error", err)
		http.Error(w, `{"error": "Invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	if createFeatureFlag.Name == "" {
		http.Error(w, `{"error": "Field 'name' is required"}`, http.StatusBadRequest)
		return
	}

	if err := h.repository.Create(r.Context(), createFeatureFlag, "test-tenant", "user-1"); err != nil {
		slog.Error("Failed to save flag to database", "error", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createFeatureFlag); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}
}
