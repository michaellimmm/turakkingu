package repository

import (
	"context"
	"errors"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TrackRepo interface {
	CreateTrack(ctx context.Context, track *entity.Track) error
}

type trackRepo struct {
	collection          *mongo.Collection
	trackingSettingRepo TrackingSettingRepo
}

func NewTrackRepo(db *mongo.Database, trackingSettingRepo TrackingSettingRepo) TrackRepo {
	collection := db.Collection("track")

	return &trackRepo{
		collection:          collection,
		trackingSettingRepo: trackingSettingRepo,
	}
}

func (r *trackRepo) CreateTrack(ctx context.Context, track *entity.Track) error {
	if err := r.validateTrackingSettingID(ctx, track.TrackingSettingID); err != nil {
		return err
	}

	track.SetCreatedAt()
	track.SetUpdatedAt()

	res, err := r.collection.InsertOne(ctx, track)
	if err != nil {
		return err
	}
	track.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *trackRepo) validateTrackingSettingID(ctx context.Context, trackingSettingID bson.ObjectID) error {
	if trackingSettingID.IsZero() {
		return errors.New("tracking_setting_id is required")
	}

	exists, err := r.trackingSettingRepo.ExistsByID(ctx, trackingSettingID)
	if err != nil {
		return fmt.Errorf("failed to check tracking setting existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("tracking setting with ID %s does not exist", trackingSettingID.Hex())
	}

	return nil
}
