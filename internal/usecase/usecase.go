package usecase

import (
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/repository"
)

type UseCase interface {
	LinkUseCase
	TrackingSettingUseCase
	TrackUseCase
	EventUseCase
}

type usecase struct {
	LinkUseCase
	TrackingSettingUseCase
	TrackUseCase
	EventUseCase
}

func NewUseCase(config *core.Config, repo repository.Repo) UseCase {
	linkUseCase := NewLinkUseCase(config, repo)
	trackingSettingUseCase := NewTrackingSettingUseCase(config, repo)
	trackUseCase := NewTrackUseCase(config, repo)
	eventUseCase := NewEventUseCase(config, repo)

	return &usecase{
		LinkUseCase:            linkUseCase,
		TrackingSettingUseCase: trackingSettingUseCase,
		TrackUseCase:           trackUseCase,
		EventUseCase:           eventUseCase,
	}
}
