package externalapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"start/internal/auth"
	"start/internal/common"

	"github.com/go-chi/chi/v5"
)

type ExternalApiHandler struct {
	repository FeatureFlagExternalReposiotory
}

func CreateExternalApiHandler(repo FeatureFlagExternalReposiotory) *ExternalApiHandler {
	return &ExternalApiHandler{
		repository: repo,
	}
}

func (h *ExternalApiHandler) GetByTenantAndName(w http.ResponseWriter, r *http.Request) {
	flagName := chi.URLParam(r, "id")
	userDetails := auth.GetUserDetails(r.Context())

	data, err := h.repository.GetByTenantAndName(r.Context(), userDetails.OrgID, flagName)
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
