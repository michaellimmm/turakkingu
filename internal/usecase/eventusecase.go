package usecase

import (
	"context"
	"errors"
	"fmt"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log/slog"
	"net/url"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventUseCase interface {
	ProcessEvent(ctx context.Context, event *entity.Event) error
}

type eventUseCase struct {
	repo   repository.Repo
	config *core.Config
}

func NewEventUseCase(config *core.Config, repo repository.Repo) EventUseCase {
	return &eventUseCase{
		repo:   repo,
		config: config,
	}
}

// WARNING!!!!
// Right now, I’ve just implemented a simple use case where one landing page can have only one conversion.
// Please add more scenario _/\_.
func (uc *eventUseCase) ProcessEvent(ctx context.Context, event *entity.Event) error {
	// TODO: check attribution window
	trackID, err := event.GetTrackID()
	if err != nil {
		// check fingerprint
		return uc.processFingerprint(ctx, event)
	}

	existingEvents, err := uc.repo.FindAllEventByTrackID(ctx, trackID)
	if err != nil && !errors.Is(err, repository.ErrNoEvents) {
		slog.Error("failed to get event by track id", slog.String("error", err.Error()))
		return err
	} else if (err != nil && errors.Is(err, repository.ErrNoEvents)) || len(existingEvents) == 0 { // new event
		track, err := uc.repo.FindTrackByID(ctx, trackID)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			slog.Error("failed to get track by id", slog.String("error", err.Error()))
			return err
		}

		if track != nil {
			// check if event url is equal with track url, it means user just open landing page
			isMatch, err := uc.isQuerySubset(track.Url, event.Url)
			if err != nil {
				return err
			}

			if !isMatch {
				return nil
			}
		}

		event.EventName = entity.EventNameLandingPage
		if err = uc.repo.CreateEvent(ctx, event); err != nil {
			return err
		}

		return nil
	}

	// check if last event is landing page
	if existingEvents[0].EventName == entity.EventNameLandingPage {
		// check if event url is in thank you page
		return uc.checkAndSaveThankYouPageEvent(ctx, trackID, event)
	}

	return nil
}

func (uc *eventUseCase) checkAndSaveThankYouPageEvent(ctx context.Context, trackID bson.ObjectID, event *entity.Event) error {
	trackPages, err := uc.repo.FindTrackByIDWithThankYouPages(ctx, trackID)
	if err != nil {
		return err
	}

	pageID, _ := uc.matchUrlInThankYouPageList(event.Url, trackPages.ThankYouPages)
	if pageID != nil {
		return nil
	}

	event.EventName = entity.EventNameThankYouPage
	if err = uc.repo.CreateEvent(ctx, event); err != nil {
		return err
	}

	// TODO: update thank you page status
	// TODO: publish conversion event

	return nil
}

func (uc *eventUseCase) matchUrlInThankYouPageList(currUrl string, pages []*entity.ThankYouPage) (*bson.ObjectID, error) {
	for _, page := range pages {
		found, err := uc.isQuerySubset(currUrl, page.URL)
		if err != nil {
			return nil, err
		}

		if found {
			return &page.ID, nil
		}
	}
	return nil, nil
}

func (uc *eventUseCase) isQuerySubset(a, b string) (bool, error) {
	parsedA, err := url.Parse(a)
	if err != nil {
		return false, fmt.Errorf("invalid URL: %w", err)
	}
	parsedB, err := url.Parse(b)
	if err != nil {
		return false, fmt.Errorf("invalid URL: %w", err)
	}

	queryA := parsedA.Query()
	queryB := parsedB.Query()

	for key, valuesA := range queryA {
		valuesB, exists := queryB[key]
		if !exists {
			return false, nil
		}
		// Check if all values of A's key are present in B's key
		for _, valA := range valuesA {
			found := false
			for _, valB := range valuesB {
				if valA == valB {
					found = true
					break
				}
			}
			if !found {
				return false, nil
			}
		}
	}

	return true, nil
}

func (uc *eventUseCase) processFingerprint(ctx context.Context, event *entity.Event) error {
	// check last event from fingerprint
	// if no event on this fingerprint return nil
	// if yes then check if url match with thank you page
	lastEvent, err := uc.repo.FindLastEventByFingerprint(ctx, event.Fingerprint)
	if err != nil && errors.Is(err, repository.ErrNoEvents) {
		// not process event when isn't tracked before
		return nil
	} else if err != nil {
		return err
	}

	if lastEvent.EventName == entity.EventNameLandingPage {
		trackID, err := lastEvent.GetTrackID()
		if err != nil {
			return nil
		}

		return uc.checkAndSaveThankYouPageEvent(ctx, trackID, event)
	}

	return nil
}
