package repository

import (
	"context"
	"errors"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrNoEvents = errors.New("no events found")
)

type EventRepo interface {
	CreateEvent(ctx context.Context, event *entity.Event) error
	FindAllEventByTenantID(ctx context.Context, tenantID string) ([]*entity.Event, error)
	FindLastEventByFingerprint(ctx context.Context, fingerprint string) (*entity.Event, error)
	FindAllEventByTrackID(ctx context.Context, trackID bson.ObjectID) ([]*entity.Event, error)
}

type eventRepo struct {
	collection *mongo.Collection
	trackRepo  TrackRepo
}

func NewEventRepo(db *mongo.Database, trackRepo TrackRepo) EventRepo {
	return &eventRepo{
		collection: db.Collection("event"),
		trackRepo:  trackRepo,
	}
}

func (r *eventRepo) CreateEvent(ctx context.Context, event *entity.Event) error {
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	res, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return err
	}

	event.ID = res.InsertedID.(bson.ObjectID)

	return nil
}

func (r *eventRepo) FindAllEventByTenantID(ctx context.Context, tenantID string) ([]*entity.Event, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "track"},
			{Key: "localField", Value: "track_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "track"},
		}}},
		bson.D{{Key: "$unwind", Value: "$track"}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "tracking_setting"},
			{Key: "localField", Value: "track.tracking_setting_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "tracking_setting"},
		}}},
		bson.D{{Key: "$unwind", Value: "$tracking_setting"}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "tracking_setting.tenant_id", Value: tenantID},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*entity.Event
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	return results, nil
}

func (r *eventRepo) FindLastEventByFingerprint(ctx context.Context, fingerprint string) (*entity.Event, error) {
	filter := bson.M{"fingerprint": fingerprint}
	opts := options.FindOne().
		SetSort(bson.D{{Key: "published_at", Value: -1}}) // sort descending

	var event entity.Event
	err := r.collection.FindOne(ctx, filter, opts).Decode(&event)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNoEvents
	} else if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepo) FindAllEventByTrackID(ctx context.Context, trackID bson.ObjectID) ([]*entity.Event, error) {
	filter := bson.M{"track_id": trackID}
	opts := options.Find().
		SetSort(bson.D{{Key: "published_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNoEvents
	} else if err != nil {
		return nil, err
	}

	var results []*entity.Event
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}

	return results, nil
}
