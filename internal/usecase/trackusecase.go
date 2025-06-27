package usecase

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
	"log/slog"
)

type TrackUseCase interface {
	CreateTrack(ctx context.Context, track *entity.Track) error
}

type trackUseCase struct {
	repo   repository.Repo
	config *core.Config
}

func NewTrackUseCase(config *core.Config, repo repository.Repo) TrackUseCase {
	return &trackUseCase{
		repo:   repo,
		config: config,
	}
}

func (uc *trackUseCase) CreateTrack(ctx context.Context, track *entity.Track) error {
	err := uc.repo.CreateTrack(ctx, track)
	if err != nil {
		slog.Error("failed to create track", slog.String("error", err.Error()))
		return err
	}

	return nil
}
