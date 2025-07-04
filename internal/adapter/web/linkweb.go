package web

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/usecase"
	"net/http"
	"strings"

	webui "github/michaellimmm/turakkingu/web"

	"github.com/a-h/templ"
)

type linkWeb struct {
	uc     usecase.UseCase
	config *core.Config
}

func NewLinkWeb(config *core.Config, uc usecase.UseCase) *linkWeb {
	return &linkWeb{
		uc:     uc,
		config: config,
	}
}

func (l *linkWeb) Index(w http.ResponseWriter, r *http.Request) {
	var component templ.Component

	// TODO: fix this
	links, err := l.uc.GetAllLinks(r.Context(), "tenant1")
	if err != nil {
		component = webui.LandingPagesContent([]webui.LandingPage{})
	} else {
		res := []webui.LandingPage{}
		for _, link := range links {
			res = append(res, webui.LandingPage{
				ID:              link.ID.Hex(),
				LandingPageName: link.Name,
				LandingPageURL:  link.Url,
				FixedURL:        link.ConstructFixedUrl(l.config.Domain),
			})
		}
		component = webui.LandingPagesContent(res)
	}

	component.Render(context.Background(), w)
}

func (l *linkWeb) Search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")

	filtered := []webui.LandingPage{}
	for _, lp := range landingPageTracker.LandingPages {
		if strings.Contains(strings.ToLower(lp.LandingPageName), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(lp.LandingPageURL), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(lp.FixedURL), strings.ToLower(query)) {
			filtered = append(filtered, lp)
		}
	}

	component := webui.LandingPagesTable(filtered)
	component.Render(context.Background(), w)
}

func (l *linkWeb) Create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("landing_page_name")
	url := r.FormValue("landing_page_url")

	link := &entity.Link{
		Name:     name,
		Url:      url,
		TenantID: "tenant1",
	}
	_ = l.uc.CreateLink(r.Context(), link)

	var component templ.Component

	// TODO: fix this
	links, err := l.uc.GetAllLinks(r.Context(), "tenant1")
	if err != nil {
		component = webui.LandingPagesTable([]webui.LandingPage{})
	} else {
		res := []webui.LandingPage{}
		for _, link := range links {
			res = append(res, webui.LandingPage{
				ID:              link.ID.Hex(),
				LandingPageName: link.Name,
				LandingPageURL:  link.Url,
				FixedURL:        link.ConstructFixedUrl(l.config.Domain),
			})
		}
		component = webui.LandingPagesTable(res)
	}

	component.Render(context.Background(), w)
}

func (l *linkWeb) Edit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	fixedURL := r.FormValue("fixed_url")
	landingPageName := r.FormValue("landing_page_name")
	landingPageURL := r.FormValue("landing_page_url")

	for i, lp := range landingPageTracker.LandingPages {
		if lp.ID == idStr {
			landingPageTracker.LandingPages[i].FixedURL = fixedURL
			landingPageTracker.LandingPages[i].LandingPageName = landingPageName
			landingPageTracker.LandingPages[i].LandingPageURL = landingPageURL
			break
		}
	}

	component := webui.LandingPagesTable(landingPageTracker.LandingPages)
	component.Render(context.Background(), w)
}

func (l *linkWeb) BulkEdit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	selectedIDs := r.Form["selected"]

	// Get the field values
	fixedURL := r.FormValue("fixed_url")
	landingPageName := r.FormValue("landing_page_name")
	landingPageURL := r.FormValue("landing_page_url")

	// Check which fields should be updated (based on checkboxes)
	updateFixedURL := r.FormValue("update_fixed_url") == "on"
	updateLandingPageName := r.FormValue("update_landing_page_name") == "on"
	updateLandingPageURL := r.FormValue("update_landing_page_url") == "on"

	// Update selected landing pages
	for _, idStr := range selectedIDs {

		for i, lp := range landingPageTracker.LandingPages {
			if lp.ID == idStr {
				if updateFixedURL && fixedURL != "" {
					landingPageTracker.LandingPages[i].FixedURL = fixedURL
				}
				if updateLandingPageName && landingPageName != "" {
					landingPageTracker.LandingPages[i].LandingPageName = landingPageName
				}
				if updateLandingPageURL && landingPageURL != "" {
					landingPageTracker.LandingPages[i].LandingPageURL = landingPageURL
				}

				break
			}
		}
	}

	component := webui.LandingPagesTable(landingPageTracker.LandingPages)
	component.Render(context.Background(), w)
}
