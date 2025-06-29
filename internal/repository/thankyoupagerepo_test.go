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

type TestSuiteThankYouPageRepo struct {
	mongoContainer      testcontainers.Container
	client              *mongo.Client
	thankYouPageRepo    repository.ThankYouPageRepo
	trackingSettingRepo repository.TrackingSettingRepo
}

func setupTestSuiteThankYouPageRepo() (*TestSuiteThankYouPageRepo, error) {
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
	repo := repository.NewThankYouPageRepo(database, trackingSettingRepo)

	return &TestSuiteThankYouPageRepo{
		mongoContainer:      mongodbContainer,
		client:              client,
		thankYouPageRepo:    repo,
		trackingSettingRepo: trackingSettingRepo,
	}, nil
}

func (ts *TestSuiteThankYouPageRepo) Cleanup() {
	ctx := context.Background()
	if ts.client != nil {
		ts.client.Disconnect(ctx)
	}
	if ts.mongoContainer != nil {
		ts.mongoContainer.Terminate(ctx)
	}
}

func TestThankYouPageRepo_CreatePage(t *testing.T) {
	suite, err := setupTestSuiteThankYouPageRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()
	t.Run("should create thankyoupage successfully", func(t *testing.T) {
		trackingSetting, err := suite.trackingSettingRepo.FindOrCreateWithPagesByTenantID(ctx, "tenant1")
		assert.NoError(t, err)

		thankYouPage := &entity.ThankYouPage{
			TrackingSettingID: trackingSetting.ID,
			URL:               "http://example.com/thank_you",
			Name:              "thank you page 1",
		}
		err = suite.thankYouPageRepo.CreatePage(ctx, thankYouPage)

		assert.NoError(t, err)
		assert.False(t, thankYouPage.ID.IsZero(), "ID should be generated")
		assert.False(t, thankYouPage.CreatedAt.IsZero(), "CreatedAt should be setted")
		assert.False(t, thankYouPage.UpdatedAt.IsZero(), "UpdatedAt should be setted")
	})

	t.Run("should return error when tracking setting id is not found", func(t *testing.T) {
		thankYouPage := &entity.ThankYouPage{
			TrackingSettingID: bson.NewObjectID(),
			URL:               "http://example.com/thank_you",
			Name:              "thank you page 1",
		}
		err = suite.thankYouPageRepo.CreatePage(ctx, thankYouPage)

		assert.Error(t, err)
	})

	t.Run("should return error when insert duplicate url", func(t *testing.T) {
		trackingSetting, err := suite.trackingSettingRepo.FindOrCreateWithPagesByTenantID(ctx, "tenant1")
		assert.NoError(t, err)

		thankYouPage := &entity.ThankYouPage{
			TrackingSettingID: trackingSetting.ID,
			URL:               "http://example.com/thank_you",
			Name:              "thank you page 1",
		}
		err = suite.thankYouPageRepo.CreatePage(ctx, thankYouPage)

		assert.Error(t, err)
	})
}
