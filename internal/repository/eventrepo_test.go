package repository_test

import (
	"context"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TestSuiteEventRepo struct {
	mongoContainer      testcontainers.Container
	client              *mongo.Client
	trackingSettingRepo repository.TrackingSettingRepo
	eventRepo           repository.EventRepo
	trackRepo           repository.TrackRepo
}

func setupTestSuiteEventRepo() (*TestSuiteEventRepo, error) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:8")
	if err != nil {
		return nil, err
	}

	endpoint, err := mongodbContainer.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", endpoint)))
	if err != nil {
		return nil, err
	}

	database := client.Database("test")
	trackingSettingRepo := repository.NewTrackingSettingRepo(database)
	trackRepo := repository.NewTrackRepo(database, trackingSettingRepo)
	eventRepo := repository.NewEventRepo(database, trackRepo)

	return &TestSuiteEventRepo{
		mongoContainer:      mongodbContainer,
		client:              client,
		trackingSettingRepo: trackingSettingRepo,
		trackRepo:           trackRepo,
		eventRepo:           eventRepo,
	}, nil
}

func (ts TestSuiteEventRepo) Cleanup() {
	ctx := context.Background()
	if ts.client != nil {
		ts.client.Disconnect(ctx)
	}
	if ts.mongoContainer != nil {
		ts.mongoContainer.Terminate(ctx)
	}
}

func TestEventRepo_CreateEvent(t *testing.T) {
	suite, err := setupTestSuiteEventRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()
	t.Run("should create event successfully", func(t *testing.T) {
		trackingSetting, err := suite.trackingSettingRepo.FindOrCreateWithPagesByTenantID(ctx, "tenant1")
		assert.NoError(t, err)
		track := &entity.Track{
			TrackingSettingID: trackingSetting.ID,
			Url:               "http://www.example.com",
			EndUserID:         "EndUserID12345",
			SessionID:         "SessionID12345",
			Platform:          "line",
		}
		err = suite.trackRepo.CreateTrack(ctx, track)
		assert.NoError(t, err)

		event := &entity.Event{
			TrackID:     track.ID,
			UserAgent:   "Mozilla/5.0",
			Fingerprint: "fingerprint123456",
			Url:         "http://www.example.com",
			EventName:   entity.EventNameLandingPage,
			PublishedAt: time.Now(),
		}
		err = suite.eventRepo.CreateEvent(ctx, event)

		assert.NoError(t, err)
		assert.False(t, event.ID.IsZero(), "ID should be generated")
		assert.False(t, event.CreatedAt.IsZero(), "CreatedAt should be setted")
		assert.False(t, event.UpdatedAt.IsZero(), "UpdatedAt should be setted")
		assert.Nil(t, event.DeletedAt, "DeletedAt should be nil")
	})

	t.Run("should return error when track id is not found", func(t *testing.T) {
		event := &entity.Event{
			TrackID:     bson.NewObjectID(),
			UserAgent:   "Mozilla/5.0",
			Fingerprint: "fingerprint123456",
			Url:         "http://www.example.com",
			EventName:   entity.EventNameLandingPage,
			PublishedAt: time.Now(),
		}
		err = suite.eventRepo.CreateEvent(ctx, event)

		assert.Error(t, err)
		assert.True(t, event.ID.IsZero(), "ID should be generated")
	})
}

func TestEventRepo_FindAllEventByTenantID(t *testing.T) {
	suite, err := setupTestSuiteEventRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()
	t.Run("should return list of events", func(t *testing.T) {
		trackingSetting, err := suite.trackingSettingRepo.FindOrCreateWithPagesByTenantID(ctx, "tenant1")
		assert.NoError(t, err)
		track := &entity.Track{
			TrackingSettingID: trackingSetting.ID,
			Url:               "http://www.example.com",
			EndUserID:         "EndUserID12345",
			SessionID:         "SessionID12345",
			Platform:          "line",
		}
		err = suite.trackRepo.CreateTrack(ctx, track)
		assert.NoError(t, err)

		event := &entity.Event{
			TrackID:     track.ID,
			UserAgent:   "Mozilla/5.0",
			Fingerprint: "fingerprint123456",
			Url:         "http://www.example.com",
			EventName:   entity.EventNameLandingPage,
			PublishedAt: time.Now(),
		}
		err = suite.eventRepo.CreateEvent(ctx, event)
		assert.NoError(t, err)
		event2 := &entity.Event{
			TrackID:     track.ID,
			UserAgent:   "Mozilla/5.0",
			Fingerprint: "fingerprint123456",
			Url:         "http://www.example.com/2",
			EventName:   entity.EventNameLandingPage,
			PublishedAt: time.Now(),
		}
		err = suite.eventRepo.CreateEvent(ctx, event2)
		assert.NoError(t, err)

		events, err := suite.eventRepo.FindAllEventByTenantID(ctx, "tenant1")

		assert.NoError(t, err)
		assert.Equal(t, len(events), 2)
	})
}

func TestEventRepo_FindLastEventByFingerprint(t *testing.T) {
	suite, err := setupTestSuiteEventRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()
	t.Run("should return event", func(t *testing.T) {
		trackingSetting, err := suite.trackingSettingRepo.FindOrCreateWithPagesByTenantID(ctx, "tenant1")
		assert.NoError(t, err)
		track := &entity.Track{
			TrackingSettingID: trackingSetting.ID,
			Url:               "http://www.example.com",
			EndUserID:         "EndUserID12345",
			SessionID:         "SessionID12345",
			Platform:          "line",
		}
		err = suite.trackRepo.CreateTrack(ctx, track)
		assert.NoError(t, err)

		event := &entity.Event{
			TrackID:     track.ID,
			UserAgent:   "Mozilla/5.0",
			Fingerprint: "fingerprint123456",
			Url:         "http://www.example.com",
			EventName:   entity.EventNameLandingPage,
			PublishedAt: time.Now(),
		}
		err = suite.eventRepo.CreateEvent(ctx, event)
		assert.NoError(t, err)

		actual, err := suite.eventRepo.FindLastEventByFingerprint(ctx, "fingerprint123456")

		assert.NoError(t, err)
		assert.Equal(t, event.ID, actual.ID)
	})
}
