package usecase

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
)

type EventUseCase interface {
	CreateEvent(ctx context.Context, event *entity.Event) error
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

func (uc *eventUseCase) CreateEvent(ctx context.Context, event *entity.Event) error {
	// update logic
	// check the thankyou page
	// if events has lp page and new event has thank you page,
	// then check the attribution window
	// if attribution window is valid then publish conversion
	return nil
}
