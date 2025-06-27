package api

import (
	"fmt"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"
	"log/slog"
	"net/http"
)

type TrackingSettingAPI interface {
	GetTrackingSetting(w http.ResponseWriter, r *http.Request)
}

type trackingSettingAPI struct {
	uc     usecase.UseCase
	config *core.Config
}

func NewTrackingSettingAPI(config *core.Config, uc usecase.UseCase) TrackingSettingAPI {
	return &trackingSettingAPI{config: config, uc: uc}
}

func (t *trackingSettingAPI) GetTrackingSetting(w http.ResponseWriter, r *http.Request) {
	tenantID := r.PathValue("tenant_id")
	if tenantID == "" {
		slog.Error("tenant id is empty")
		_ = sendError(w, http.StatusBadRequest, fmt.Errorf("tenant_id can"))
		return
	}

	response, err := t.uc.GetTrackingSettingByTenantID(r.Context(), tenantID)
	if err != nil {
		slog.Error("failed to get tracking setting", slog.String("error", err.Error()))
		sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get tracking setting"))
		return
	}

	sendJson(w, http.StatusOK, response)
}
