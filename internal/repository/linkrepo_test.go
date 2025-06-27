package repository_test

import (
	"context"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"
	"github/michaellimmm/turakkingu/internal/repository"
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TestSuiteLinkRepo struct {
	mongoContainer testcontainers.Container
	client         *mongo.Client
	repo           repository.LinkRepo
}

func setupTestSuiteLinkRepo() (*TestSuiteLinkRepo, error) {
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
	repo := repository.NewLinkRepo(database)

	return &TestSuiteLinkRepo{
		mongoContainer: mongodbContainer,
		client:         client,
		repo:           repo,
	}, nil
}

func (ts *TestSuiteLinkRepo) Cleanup() {
	ctx := context.Background()
	if ts.client != nil {
		ts.client.Disconnect(ctx)
	}
	if ts.mongoContainer != nil {
		ts.mongoContainer.Terminate(ctx)
	}
}

func TestLinkRepo_Create(t *testing.T) {
	suite, err := setupTestSuiteLinkRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()
	t.Run("should create link successfully", func(t *testing.T) {
		link := &entity.Link{
			TenantID: "tenat1",
			Url:      "https://www.github.com",
		}
		err := suite.repo.CreateLink(ctx, link)

		assert.NoError(t, err)
		assert.False(t, link.ID.IsZero(), "ID should be generated")
		assert.NotEmpty(t, link.TenantID, "TenantID should not be empty")
		assert.NotEmpty(t, link.ShortID, "ShortID should not be empty")
		assert.False(t, link.CreatedAt.IsZero(), "CreatedAt should be setted")
		assert.False(t, link.UpdatedAt.IsZero(), "UpdatedAt should be setted")
		assert.Nil(t, link.DeletedAt, "DeletedAt should be nil")
	})
}

func TestLinkRepo_FindByID(t *testing.T) {
	suite, err := setupTestSuiteLinkRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()
	t.Run("should find existing link", func(t *testing.T) {
		link := &entity.Link{
			TenantID: "tenat1",
			Url:      "https://www.github.com",
		}
		err := suite.repo.CreateLink(ctx, link)
		assert.NoError(t, err)

		foundLink, err := suite.repo.FindLinkByID(ctx, link.ID.Hex())

		assert.NoError(t, err)
		assert.Equal(t, foundLink.ID, link.ID)
		assert.Equal(t, foundLink.ShortID, link.ShortID)
		assert.Equal(t, foundLink.TenantID, link.TenantID)
		assert.Equal(t, foundLink.Url, link.Url)
		assert.Nil(t, foundLink.DeletedAt)
	})

	t.Run("should return error for non-existing link", func(t *testing.T) {
		nonExistentID := bson.NewObjectID()

		link, err := suite.repo.FindLinkByID(ctx, nonExistentID.Hex())
		assert.Error(t, err)
		assert.Nil(t, link)
	})
}

func TestLinkRepo_FindByShortID(t *testing.T) {
	suite, err := setupTestSuiteLinkRepo()
	assert.NoError(t, err)
	defer suite.Cleanup()

	ctx := context.Background()
	t.Run("should find existing link", func(t *testing.T) {
		id, err := gonanoid.New()
		assert.NoError(t, err)
		link := &entity.Link{
			TenantID: "tenat1",
			Url:      "https://www.github.com",
			ShortID:  id,
		}
		err = suite.repo.CreateLink(ctx, link)
		assert.NoError(t, err)

		foundLink, err := suite.repo.FindLinkByShortID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, foundLink.ID, link.ID)
		assert.Equal(t, foundLink.ShortID, link.ShortID)
		assert.Equal(t, foundLink.TenantID, link.TenantID)
		assert.Equal(t, foundLink.Url, link.Url)
		assert.Nil(t, foundLink.DeletedAt)
	})

	t.Run("should return error for non-existing link", func(t *testing.T) {
		nonExistentID, err := gonanoid.New()
		assert.NoError(t, err)

		link, err := suite.repo.FindLinkByShortID(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Nil(t, link)
	})
}
