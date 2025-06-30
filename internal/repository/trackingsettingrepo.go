package repository

import (
	"context"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// TODO: put unique index on tenant_id

// get tracking settings
type TrackingSettingRepo interface {
	FindOrCreateWithPagesByTenantID(ctx context.Context, tenantID string) (*entity.TrackingSettingWithPages, error)
	FindTrackingSettingByID(ctx context.Context, id bson.ObjectID) (*entity.TrackingSetting, error)
	FindTrackingSettingWithPagesByID(ctx context.Context, trackingSettingID bson.ObjectID) (*entity.TrackingSettingWithPages, error)
	IsTrackingSettingIDExist(ctx context.Context, id bson.ObjectID) (bool, error)
}

type trackingSettingRepo struct {
	collection *mongo.Collection
}

func NewTrackingSettingRepo(db *mongo.Database) TrackingSettingRepo {
	collection := db.Collection("tracking_setting")

	return &trackingSettingRepo{
		collection: collection,
	}
}

func (r *trackingSettingRepo) FindOrCreateWithPagesByTenantID(ctx context.Context,
	tenantID string) (*entity.TrackingSettingWithPages, error) {
	existing, err := r.FindTrakingSettingWithPagesByTenantID(ctx, tenantID)
	if err == nil {
		return existing, nil
	}
	if err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to check existing tracking setting: %w", err)
	}

	now := time.Now()
	setting := &entity.TrackingSetting{
		TenantID: tenantID,
		BaseEntity: entity.BaseEntity{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	insertResult, err := r.collection.InsertOne(ctx, setting)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracking setting: %w", err)
	}
	setting.ID = insertResult.InsertedID.(bson.ObjectID)

	result := &entity.TrackingSettingWithPages{
		ID:            setting.ID,
		TenantID:      setting.TenantID,
		BaseEntity:    setting.BaseEntity,
		ThankYouPages: []entity.ThankYouPage{},
	}

	return result, nil
}

func (r *trackingSettingRepo) FindTrackingSettingWithPagesByID(ctx context.Context,
	trackingSettingID bson.ObjectID) (*entity.TrackingSettingWithPages, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": trackingSettingID},
		},
		{
			"$lookup": bson.M{
				"from":         "thank_you_page",
				"localField":   "_id",
				"foreignField": "tracking_setting_id",
				"as":           "thank_you_pages",
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		slog.Error("failed to aggregate tracking setting", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to aggregate tracking setting: %w", err)
	}
	defer cursor.Close(ctx)

	var results []entity.TrackingSettingWithPages
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &results[0], nil
}

func (r *trackingSettingRepo) FindTrakingSettingWithPagesByTenantID(ctx context.Context,
	tenantID string) (*entity.TrackingSettingWithPages, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"tenant_id": tenantID},
		},
		{
			"$lookup": bson.M{
				"from":         "thank_you_page",
				"localField":   "_id",
				"foreignField": "tracking_setting_id",
				"as":           "thank_you_pages",
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		slog.Error("failed to aggregate tracking setting", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to aggregate tracking setting: %w", err)
	}
	defer cursor.Close(ctx)

	var results []entity.TrackingSettingWithPages
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &results[0], nil
}

// TODO: update fields (PATCH)
func (r *trackingSettingRepo) UpdateFieldsAndReturn() {

}

func (r *trackingSettingRepo) IsTrackingSettingIDExist(ctx context.Context, id bson.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
}

func (r *trackingSettingRepo) FindTrackingSettingByID(ctx context.Context, id bson.ObjectID) (*entity.TrackingSetting, error) {
	filter := bson.M{"_id": id}

	var setting entity.TrackingSetting
	err := r.collection.FindOne(ctx, filter).Decode(&setting)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}
