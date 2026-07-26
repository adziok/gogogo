package feature_flags

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"start/internal/auth"
	"start/internal/common"

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

func (h *FeatureFlagHandler) CreateFlag(w http.ResponseWriter, r *http.Request) {
	var createFeatureFlag CreateFeatureFlag

	if err := json.NewDecoder(r.Body).Decode(&createFeatureFlag); err != nil {
		common.RenderErr(w, r, http.StatusBadRequest, err, "Invalid JSON body")
		return
	}

	if err := validate.Struct(createFeatureFlag); err != nil {
		common.RenderErr(w, r, http.StatusUnprocessableEntity, err, "Validation failed")
		return
	}

	userDetails := auth.GetUserDetails(r.Context())
	if err := h.repository.Create(r.Context(), Operation[CreateFeatureFlag]{Data: createFeatureFlag, Tenant: userDetails.OrgID, User: userDetails.UserID}); err != nil {
		common.RenderErr(w, r, http.StatusInternalServerError, err, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createFeatureFlag); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}
}

func (h *FeatureFlagHandler) DisplayFlags(w http.ResponseWriter, r *http.Request) {
	userDetails := auth.GetUserDetails(r.Context())
	data, err := h.repository.GetByTenant(r.Context(), userDetails.OrgID)
	if err != nil {
		common.RenderErr(w, r, http.StatusInternalServerError, err, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}
}

func (h *FeatureFlagHandler) UpdateFlag(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "id")

	if flagID == "" {
		common.RenderErr(w, r, http.StatusBadRequest, nil, "Flag ID is required")
		return
	}

	updateFeatureFlag := UpdateFeatureFlag{ID: flagID}

	if err := json.NewDecoder(r.Body).Decode(&updateFeatureFlag); err != nil {
		common.RenderErr(w, r, http.StatusBadRequest, err, "Invalid JSON body")
		return
	}

	if err := validate.Struct(updateFeatureFlag); err != nil {
		common.RenderErr(w, r, http.StatusUnprocessableEntity, err, "Validation failed")
		return
	}

	userDetails := auth.GetUserDetails(r.Context())
	if err := h.repository.Update(r.Context(), Operation[UpdateFeatureFlag]{Data: updateFeatureFlag, Tenant: userDetails.OrgID, User: userDetails.UserID}); err != nil {
		common.RenderErr(w, r, http.StatusInternalServerError, err, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(updateFeatureFlag); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}
}

func (h *FeatureFlagHandler) DeleteFlag(w http.ResponseWriter, r *http.Request) {
	flagID := chi.URLParam(r, "id")

	deleteFeatureFlag := DeleteFeatureFlag{
		ID: flagID,
	}

	if err := validate.Struct(deleteFeatureFlag); err != nil {
		common.RenderErr(w, r, http.StatusUnprocessableEntity, err, "Validation failed")
		return
	}

	userDetails := auth.GetUserDetails(r.Context())
	if err := h.repository.DeleteById(r.Context(), Operation[DeleteFeatureFlag]{Data: deleteFeatureFlag, Tenant: userDetails.OrgID, User: userDetails.UserID}); err != nil {
		common.RenderErr(w, r, http.StatusInternalServerError, err, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	if err := json.NewEncoder(w).Encode(flagID); err != nil {
		slog.Error("Failed to encode response JSON", "error", err)
	}
}
