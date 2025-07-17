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
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

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

	_, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return fmt.Errorf("url is not valid")
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
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	QueryParam string    `json:"query_param"`
}

type TrackEventRequest struct {
	TrackID     string `json:"track_id"`
	URL         string `json:"url"`
	Fingerprint string `json:"fp"`
	PublishedAt int64  `json:"published_at"`
}

func (t *TrackEventRequest) GetPublishedAt() time.Time {
	res := time.Unix(0, t.PublishedAt*int64(time.Millisecond))
	return res
}

func (t *TrackEventRequest) Validate() error {
	if t.TrackID == "" {
		return fmt.Errorf("track_id can not be empty")
	}

	if t.URL == "" {
		return fmt.Errorf("url can not be empty")
	}

	_, err := url.ParseRequestURI(t.URL)
	if err != nil {
		return fmt.Errorf("url is not valid")
	}

	return nil
}

func (t *TrackEventRequest) FromReader(rc io.ReadCloser) error {
	defer func() {
		_ = rc.Close()
	}()
	return json.NewDecoder(rc).Decode(t)
}

func (t *TrackEventRequest) ToString() string {
	s, _ := json.Marshal(t)
	return string(s)
}

type TrackEventResponse struct{}

func NewTrackAPI(config *core.Config, uc usecase.UseCase) *trackAPI {
	return &trackAPI{config: config, uc: uc}
}

func (t *trackAPI) CreateTrack(w http.ResponseWriter, r *http.Request) {
	req := &CreateTrackRequest{}
	if err := req.FromReader(r.Body); err != nil {
		slog.Error("failed to read request", slog.String("error", err.Error()))
		_ = sendError(w, http.StatusBadRequest, fmt.Errorf("failed to read request"))
		return
	}

	if err := req.Validate(); err != nil {
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
	if err := t.uc.CreateTrack(r.Context(), track); err != nil {
		slog.Error("failed to create new track", slog.String("error", err.Error()))
		sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to create new track"))
		return
	}

	epoch := track.CreatedAt.Unix()
	sendJson(w, http.StatusCreated, CreateTrackResponse{
		ID:         track.ID.Hex(),
		CreatedAt:  track.CreatedAt,
		QueryParam: fmt.Sprintf("ztid=%s&ztts=%d", track.ID.Hex(), epoch),
	})
}

func (t *trackAPI) TrackEvent(w http.ResponseWriter, r *http.Request) {
	req := &TrackEventRequest{}
	if err := req.FromReader(r.Body); err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		slog.Error("request is not valid", slog.String("error", err.Error()))
		_ = sendError(w, http.StatusBadRequest, err)
		return
	}

	slog.Info("event request", slog.String("request", req.ToString()))

	event := &entity.Event{
		TrackID:     req.TrackID,
		UserAgent:   r.UserAgent(),
		Url:         req.URL,
		PublishedAt: req.GetPublishedAt(),
		Fingerprint: req.Fingerprint,
	}
	err := t.uc.ProcessEvent(r.Context(), event)
	if err != nil {
		slog.Error("failed to process event", slog.String("error", err.Error()))
		sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to add new event"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
