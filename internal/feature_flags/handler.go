package feature_flags

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"start/internal/auth"

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

	if err := validate.Struct(createFeatureFlag); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		w.Write([]byte(`{"error": "Validation failed", "details": "` + err.Error() + `"}`))
		return
	}

	userDetails := auth.GetUserDetails(r.Context())
	if err := h.repository.Create(r.Context(), Operation[CreateFeatureFlag]{Data: createFeatureFlag, Tenant: userDetails.OrgID, User: userDetails.UserID}); err != nil {
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
	userDetails := auth.GetUserDetails(r.Context())
	slog.Info(userDetails.UserID, userDetails.OrgID)
	data, err := h.repository.GetByTenant(r.Context(), userDetails.OrgID)
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

	if err := validate.Struct(updateFeatureFlag); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		w.Write([]byte(`{"error": "Validation failed", "details": "` + err.Error() + `"}`))
		return
	}

	userDetails := auth.GetUserDetails(r.Context())
	if err := h.repository.Update(r.Context(), Operation[UpdateFeatureFlag]{Data: updateFeatureFlag, Tenant: userDetails.OrgID, User: userDetails.UserID}); err != nil {
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

	deleteFeatureFlag := DeleteFeatureFlag{
		ID: flagID,
	}

	if err := validate.Struct(deleteFeatureFlag); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		w.Write([]byte(`{"error": "Validation failed", "details": "` + err.Error() + `"}`))
		return
	}

	userDetails := auth.GetUserDetails(r.Context())
	if err := h.repository.DeleteById(r.Context(), Operation[DeleteFeatureFlag]{Data: deleteFeatureFlag, Tenant: userDetails.OrgID, User: userDetails.UserID}); err != nil {
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
