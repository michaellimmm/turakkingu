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
)

type LinkAPI interface {
	CreateLink(w http.ResponseWriter, r *http.Request)
	Redirect(w http.ResponseWriter, r *http.Request)
}

type CreateLinkRequest struct {
	TenantID string `json:"tenant_id"`
	Url      string `json:"url"`
}

func (c *CreateLinkRequest) Validate() error {
	if c.TenantID == "" {
		return fmt.Errorf("tenant_id can not be empty")
	}

	if c.Url == "" {
		return fmt.Errorf("url can not be empty")
	}

	_, err := url.ParseRequestURI(c.Url)
	if err != nil {
		return fmt.Errorf("url is not valid")
	}

	return nil
}

func (c *CreateLinkRequest) FromReader(r io.ReadCloser) error {
	defer func() {
		_ = r.Close()
	}()
	return json.NewDecoder(r).Decode(c)
}

type CreateLinkResponse struct {
	Link string `json:"link"`
}

type linkAPI struct {
	uc     usecase.UseCase
	config *core.Config
}

func NewLinkAPI(config *core.Config, uc usecase.UseCase) LinkAPI {
	return &linkAPI{
		uc:     uc,
		config: config,
	}
}

func (f *linkAPI) CreateLink(w http.ResponseWriter, r *http.Request) {
	req := &CreateLinkRequest{}
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

	link := &entity.Link{
		Url:      req.Url,
		TenantID: req.TenantID,
	}
	err = f.uc.CreateLink(r.Context(), link)
	if err != nil {
		slog.Error("failed to create new link", slog.String("error", err.Error()))
		sendError(w, http.StatusInternalServerError, fmt.Errorf("failed to create new link"))
		return
	}

	sendJson(w, http.StatusCreated, CreateLinkResponse{Link: link.ConstructFixedUrl(f.config.Domain)})
}

func (f *linkAPI) Redirect(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	link, err := f.uc.GetFixedLink(r.Context(), id)
	if err != nil {
		slog.Error("failed to get link", slog.String("error", err.Error()))
		render404(w)
		return
	}

	newParams := r.URL.Query()
	u, _ := url.Parse(link.Url)
	existingParams := u.Query()
	for key, values := range newParams {
		for _, v := range values {
			existingParams.Add(key, v) // Use .Set() if you want to replace instead
		}
	}
	u.RawQuery = existingParams.Encode()

	http.Redirect(w, r, u.String(), http.StatusFound)
}
