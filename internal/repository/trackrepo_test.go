package repository_test

import (
	"context"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TestSuiteTrackRepo struct {
	mongoContainer      testcontainers.Container
	client              *mongo.Client
	trackRepo           repository.TrackRepo
	trackingSettingRepo repository.TrackingSettingRepo
}

func setupTestSuiteTrackRepo() (*TestSuiteTrackRepo, error) {
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
	repo := repository.NewTrackRepo(database, trackingSettingRepo)
	return &TestSuiteTrackRepo{
		mongoContainer:      mongodbContainer,
		client:              client,
		trackRepo:           repo,
		trackingSettingRepo: trackingSettingRepo,
	}, nil
}

func (ts *TestSuiteTrackRepo) Cleanup() {
	ctx := context.Background()
	if ts.client != nil {
		ts.client.Disconnect(ctx)
	}
	if ts.mongoContainer != nil {
		ts.mongoContainer.Terminate(ctx)
	}
}

func TestTrackRepo_Create(t *testing.T) {
	suite, err := setupTestSuiteTrackRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()
	t.Run("should create track successfully", func(t *testing.T) {
		trackingSetting, err := suite.trackingSettingRepo.FindOrCreateWithPagesByTenantID(ctx, "tenant1")
		assert.NoError(t, err)

		track := &entity.Track{
			TrackingSettingID: trackingSetting.ID,
			Url:               "https://www.google.com",
			EndUserID:         "EndUserID12345",
			SessionID:         "SessionID12345",
			Platform:          "line",
		}
		err = suite.trackRepo.CreateTrack(ctx, track)

		assert.NoError(t, err)
		assert.False(t, track.ID.IsZero(), "ID should be generated")
		assert.False(t, track.CreatedAt.IsZero(), "CreatedAt should be setted")
		assert.False(t, track.UpdatedAt.IsZero(), "UpdatedAt should be setted")
		assert.Nil(t, track.DeletedAt, "DeletedAt should be nil")
	})

	t.Run("should return error when tracking setting id is not found", func(t *testing.T) {
		track := &entity.Track{
			TrackingSettingID: bson.NewObjectID(),
			Url:               "https://www.google.com",
			EndUserID:         "EndUserID12345",
			SessionID:         "SessionID12345",
			Platform:          "line",
		}
		err = suite.trackRepo.CreateTrack(ctx, track)

		assert.Error(t, err)
	})
}
