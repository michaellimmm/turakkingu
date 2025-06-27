package repository_test

import (
	"context"
	"fmt"
	"github/michaellimmm/turakkingu/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TestSuiteTrackingSettingRepo struct {
	mongoContainer testcontainers.Container
	client         *mongo.Client
	repo           repository.TrackingSettingRepo
}

func setupTestSuiteTrackingSettingRepo() (*TestSuiteTrackingSettingRepo, error) {
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
	repo := repository.NewTrackingSettingRepo(database)

	return &TestSuiteTrackingSettingRepo{
		mongoContainer: mongodbContainer,
		client:         client,
		repo:           repo,
	}, nil
}

func (ts *TestSuiteTrackingSettingRepo) Cleanup() {
	ctx := context.Background()
	if ts.client != nil {
		ts.client.Disconnect(ctx)
	}
	if ts.mongoContainer != nil {
		ts.mongoContainer.Terminate(ctx)
	}
}

func TestTrackingSettingRepo_FindOrCreateWithPages(t *testing.T) {
	suite, err := setupTestSuiteTrackingSettingRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()

	t.Run("should create tracking setting successfully", func(t *testing.T) {
		tracking, err := suite.repo.FindOrCreateWithPagesByTenantID(ctx, "tenant1")

		assert.NoError(t, err)
		assert.False(t, tracking.ID.IsZero(), "ID should be generated")
		assert.Equal(t, tracking.TenantID, "tenant1", "TenantID should be not empty")
		assert.Empty(t, tracking.ThankYouPages, "Thank you page should be empty")
		assert.False(t, tracking.CreatedAt.IsZero(), "CreatedAt should be setted")
		assert.False(t, tracking.UpdatedAt.IsZero(), "UpdatedAt should be setted")
	})
}
