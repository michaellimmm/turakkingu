package repository

import (
	"context"
	"github/michaellimmm/turakkingu/internal/core"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repo interface {
	LinkRepo
	TrackingSettingRepo
	TrackRepo
}

type RepoCloser interface {
	Repo
	Close(context.Context) error
}

type repo struct {
	client *mongo.Client
	LinkRepo
	TrackingSettingRepo
	TrackRepo
}

func NewRepo(config *core.Config) (RepoCloser, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(config.MongoDBUri))
	if err != nil {
		return nil, err
	}

	db := client.Database(config.MongoDBName)

	linkRepo := NewLinkRepo(db)
	trackingSettingRepo := NewTrackingSettingRepo(db)
	trackRepo := NewTrackRepo(db, trackingSettingRepo)

	return &repo{
		client:              client,
		LinkRepo:            linkRepo,
		TrackingSettingRepo: trackingSettingRepo,
		TrackRepo:           trackRepo,
	}, nil
}

func (r *repo) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
