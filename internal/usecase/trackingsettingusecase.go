package usecase

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
	"log/slog"
)

type TrackingSettingUseCase interface {
	GetTrackingSettingByTenantID(ctx context.Context, tenantID string) (*entity.TrackingSettingWithPages, error)
}

type trackingSettingUseCase struct {
	repo   repository.Repo
	config *core.Config
}

func NewTrackingSettingUseCase(config *core.Config, repo repository.Repo) TrackingSettingUseCase {
	return &trackingSettingUseCase{
		repo:   repo,
		config: config,
	}
}

func (uc *trackingSettingUseCase) GetTrackingSettingByTenantID(ctx context.Context, tenantID string) (*entity.TrackingSettingWithPages, error) {
	trackingSetting, err := uc.repo.FindOrCreateWithPagesByTenantID(ctx, tenantID)
	if err != nil {
		slog.Error("faield to find or create tracking setting with pages", slog.String("error", err.Error()))
		return nil, err
	}

	return trackingSetting, nil
}
