package feature_flags

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	if err := h.repository.Create(r.Context(), Operation[CreateFeatureFlag]{Data: createFeatureFlag, Tenant: "test-tenant", User: "user-1"}); err != nil {
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

func (h FeatureFlagHandler) DisplayFlags(w http.ResponseWriter, r *http.Request) {
	data, err := h.repository.GetByTenant(r.Context(), "test-tenant")
	if err != nil {
		slog.Error("Failed to save flag to database", "error", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}
}

func (h FeatureFlagHandler) UpdateFlag(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "id")

	if flagID == "" {
		http.Error(w, `{"error": "flagID is required"}`, http.StatusBadRequest)
		return
	}

	updateFeatureFlag := UpdateFeatureFlag{ID: flagID}

	if err := json.NewDecoder(r.Body).Decode(&updateFeatureFlag); err != nil {
		slog.Warn("Failed to decode create flag request", "error", err)
		http.Error(w, `{"error": "Invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	if err := h.repository.Update(r.Context(), Operation[UpdateFeatureFlag]{Data: updateFeatureFlag, Tenant: "test-tenant", User: "user-1"}); err != nil {
		slog.Error("Failed to update flag to database", "error", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(updateFeatureFlag); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}
}

func (h FeatureFlagHandler) DeleteFlag(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "id")

	if err := h.repository.DeleteById(r.Context(), Operation[DeleteFeatureFlag]{Data: DeleteFeatureFlag{
		ID: flagID,
	}, Tenant: "test-tenant", User: "user-1"}); err != nil {
		slog.Error("Failed to delete flag to database", "error", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	if err := json.NewEncoder(w).Encode(flagID); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}
}
