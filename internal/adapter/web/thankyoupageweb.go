package web

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/usecase"
	webui "github/michaellimmm/turakkingu/web"
	"net/http"
	"strings"
)

type thankYouPageWeb struct {
	uc     usecase.UseCase
	config *core.Config
}

func NewThankYouPageWeb(config *core.Config, uc usecase.UseCase) *thankYouPageWeb {
	return &thankYouPageWeb{
		uc:     uc,
		config: config,
	}
}

// TODO: fix this
func (t *thankYouPageWeb) Index(w http.ResponseWriter, r *http.Request) {
	trackingSetting, _ := t.uc.GetTrackingSettingByTenantID(r.Context(), "tenant1")

	res := &webui.ConversionTracker{}
	res.ConversionPoints = []webui.ConversionPoint{}
	for _, page := range trackingSetting.ThankYouPages {
		res.ConversionPoints = append(res.ConversionPoints, webui.ConversionPoint{
			ID:     page.ID.Hex(),
			Name:   page.Name,
			URL:    page.URL,
			Status: "Draft",
		})
	}

	component := webui.IndexPage(res)
	component.Render(context.Background(), w)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")

	filtered := []webui.ConversionPoint{}
	for _, cp := range conversionTracker.ConversionPoints {
		if strings.Contains(strings.ToLower(cp.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(cp.URL), strings.ToLower(query)) {
			filtered = append(filtered, cp)
		}
	}

	component := webui.ConversionPointsTable(filtered)
	component.Render(context.Background(), w)
}

func handleFilter(w http.ResponseWriter, r *http.Request) {
	status := r.FormValue("status")

	filtered := []webui.ConversionPoint{}
	for _, cp := range conversionTracker.ConversionPoints {
		if status == "" || status == "all" || cp.Status == status {
			filtered = append(filtered, cp)
		}
	}

	component := webui.ConversionPointsTable(filtered)
	component.Render(context.Background(), w)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	url := r.FormValue("url")

	newCP := webui.ConversionPoint{
		ID:     string(len(conversionTracker.ConversionPoints) + 1),
		Name:   name,
		URL:    url,
		Status: "Draft",
	}

	conversionTracker.ConversionPoints = append(conversionTracker.ConversionPoints, newCP)

	component := webui.ConversionPointsTable(conversionTracker.ConversionPoints)
	component.Render(context.Background(), w)
}

func handleStartTracking(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	selectedIDs := r.Form["selected"]

	for _, idStr := range selectedIDs {

		for i, cp := range conversionTracker.ConversionPoints {
			if cp.ID == idStr {
				conversionTracker.ConversionPoints[i].Status = "Active"
				break
			}
		}
	}

	component := webui.ConversionPointsTable(conversionTracker.ConversionPoints)
	component.Render(context.Background(), w)
}
