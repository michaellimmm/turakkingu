package usecase

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TrackingSettingUseCase interface {
	GetTrackingSettingByTenantID(ctx context.Context, tenantID string) (*entity.TrackingSettingWithPages, error)
	GetTrackingSettingByID(ctx context.Context, trackingSettingID bson.ObjectID) (*entity.TrackingSettingWithPages, error)
	AddThankYouPage(ctx context.Context, thankYouPage *entity.ThankYouPage) error
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
		slog.Error("failed to find or create tracking setting with pages by tenant", slog.String("error", err.Error()))
		return nil, err
	}

	return trackingSetting, nil
}

func (uc *trackingSettingUseCase) GetTrackingSettingByID(ctx context.Context, trackingSettingID bson.ObjectID) (*entity.TrackingSettingWithPages, error) {
	trackingSetting, err := uc.repo.FindTrackingSettingWithPagesByID(ctx, trackingSettingID)
	if err != nil {
		slog.Error("failed to find tracking setting with page  by id", slog.String("error", err.Error()))
		return nil, err
	}

	return trackingSetting, nil
}

func (uc *trackingSettingUseCase) AddThankYouPage(ctx context.Context, thankYouPage *entity.ThankYouPage) error {
	thankYouPage.Status = entity.TrackingStatusPending
	err := uc.repo.CreatePage(ctx, thankYouPage)
	if err != nil {
		slog.Error("failed to add thank you page", slog.String("error", err.Error()))
		return err
	}

	return nil
}
