package api

import (
	"encoding/json"
	"fmt"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/usecase"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TrackingSettingAPI interface {
	GetTrackingSetting(w http.ResponseWriter, r *http.Request)
	AddThankYouPage(w http.ResponseWriter, r *http.Request)
}

type trackingSettingAPI struct {
	uc     usecase.UseCase
	config *core.Config
}

type AddThankYouPageRequest struct {
	TrackingSettingID string `json:"tracking_setting_id"`
	URL               string `json:"url"`
	Name              string `json:"name"`
	Point             int    `json:"point"`
}

func (r *AddThankYouPageRequest) GetTrackingSettingID() (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(r.TrackingSettingID)
}

func (r *AddThankYouPageRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name can not be empty")
	}

	if _, err := r.GetTrackingSettingID(); err != nil {
		return fmt.Errorf("tracking_setting_id is not valid")
	}

	if r.URL == "" {
		return fmt.Errorf("url can not be empty")
	}

	_, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return fmt.Errorf("url is not valid")
	}

	return nil
}

func (r *AddThankYouPageRequest) FromReader(rc io.ReadCloser) error {
	defer func() {
		_ = rc.Close()
	}()
	return json.NewDecoder(rc).Decode(r)
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

func (t *trackingSettingAPI) AddThankYouPage(w http.ResponseWriter, r *http.Request) {
	req := &AddThankYouPageRequest{}
	err := req.FromReader(r.Body)
	if err != nil {
		slog.Error("failed to read request", slog.String("error", err.Error()))
		_ = sendError(w, http.StatusBadRequest, fmt.Errorf("failed to read request"))
		return
	}

	err = req.Validate()
	if err != nil {
		slog.Error("request is not valid", slog.String("error", err.Error()))
		_ = sendError(w, http.StatusBadRequest, err)
		return
	}

	trackingSettingID, _ := req.GetTrackingSettingID()
	thankYouPage := &entity.ThankYouPage{
		TrackingSettingID: trackingSettingID,
		URL:               req.URL,
		Point:             req.Point,
		Name:              req.Name,
	}
	err = t.uc.AddThankYouPage(r.Context(), thankYouPage)
	if err != nil {
		slog.Error("failed to add thank you page", slog.String("error", err.Error()))
		sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to add thank you page"))
		return
	}

	response, err := t.uc.GetTrackingSettingByID(r.Context(), trackingSettingID)
	if err != nil {
		slog.Error("failed to get tracking setting", slog.String("error", err.Error()))
		sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to get tracking setting"))
		return
	}

	sendJson(w, http.StatusCreated, response)
}
