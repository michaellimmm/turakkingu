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
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TrackAPI interface {
	CreateTrack(w http.ResponseWriter, r *http.Request)
	TrackEvent(w http.ResponseWriter, r *http.Request)
}

type trackAPI struct {
	uc     usecase.UseCase
	config *core.Config
}

type CreateTrackRequest struct {
	TenantID          string            `json:"tenant_id"`
	TrackingSettingID string            `json:"tracking_setting_id"`
	URL               string            `json:"url"`        // LP page -> to send click event
	SessionID         string            `json:"session_id"` // should be mandatory ???
	EndUserID         string            `json:"end_user_id"`
	Platform          string            `json:"platform"`
	GeneratedFrom     string            `json:"generated_from"` // source
	Metadata          map[string]string `json:"metadata"`
}

func (r *CreateTrackRequest) GetTrackingSettingID() (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(r.TrackingSettingID)
}

func (r *CreateTrackRequest) Validate() error {
	if r.TenantID == "" {
		return fmt.Errorf("tenant_id can not be empty")
	}

	if _, err := r.GetTrackingSettingID(); err != nil {
		return fmt.Errorf("tracking_setting_id is not valid")
	}

	if r.URL == "" {
		return fmt.Errorf("url can not be empty")
	}

	if r.EndUserID == "" {
		return fmt.Errorf("end_user_id can not be empty")
	}

	return nil
}

func (r *CreateTrackRequest) FromReader(rc io.ReadCloser) error {
	defer func() {
		_ = rc.Close()
	}()
	return json.NewDecoder(rc).Decode(r)
}

type CreateTrackResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type TrackEventRequest struct {
	TrackID   string    `json:"track"`
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
}

func (t *TrackEventRequest) Validate() error {
	if t.TrackID == "" {
		return fmt.Errorf("track_id can not be empty")
	}

	if t.URL == "" {
		return fmt.Errorf("url can not be empty")
	}
	return nil
}

func (t *TrackEventRequest) FromReader(rc io.ReadCloser) error {
	defer func() {
		_ = rc.Close()
	}()
	return json.NewDecoder(rc).Decode(t)
}

type TrackEventResponse struct{}

func NewTrackAPI(config *core.Config, uc usecase.UseCase) TrackAPI {
	return &trackAPI{config: config, uc: uc}
}

func (t *trackAPI) CreateTrack(w http.ResponseWriter, r *http.Request) {
	req := &CreateTrackRequest{}
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
	track := &entity.Track{
		TrackingSettingID: trackingSettingID,
		Url:               req.URL,
		EndUserID:         req.EndUserID,
		SessionID:         req.SessionID,
		Platform:          req.Platform,
		GeneratedFrom:     req.GeneratedFrom,
		Metadata:          req.Metadata,
	}
	t.uc.CreateTrack(r.Context(), track)

	sendJson(w, http.StatusCreated, CreateTrackResponse{
		ID:        track.ID.Hex(),
		CreatedAt: track.CreatedAt,
	})
}

func (t *trackAPI) TrackEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Print the request body
	fmt.Printf("Request Body: %s\n", string(body))
	w.WriteHeader(http.StatusNoContent)
}
