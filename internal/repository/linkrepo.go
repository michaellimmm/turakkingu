package repository

import (
	"context"
	"errors"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type LinkRepo interface {
	CreateLink(context.Context, *entity.Link) error
	FindLinkByID(context.Context, string) (*entity.Link, error)
	FindLinkByShortID(context.Context, string) (*entity.Link, error)
	FindAllLinkbyTenantID(ctx context.Context, tenantID string) ([]*entity.Link, error)
}

type linkRepo struct {
	collection *mongo.Collection
}

func NewLinkRepo(db *mongo.Database) LinkRepo {
	collection := db.Collection("link")

	return &linkRepo{
		collection: collection,
	}
}

func (r *linkRepo) CreateLink(ctx context.Context, link *entity.Link) error {
	err := link.SetShortID()
	if err != nil {
		return fmt.Errorf("failed to set short id")
	}

	link.SetCreatedAt()
	link.SetUpdatedAt()

	res, err := r.collection.InsertOne(ctx, link)
	if err != nil {
		return err
	}
	link.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *linkRepo) FindLinkByID(ctx context.Context, id string) (*entity.Link, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("id is not valid")
	}

	var link entity.Link
	filter := bson.M{
		"_id":        oid,
		"deleted_at": bson.M{"$exists": false},
	}

	err = r.collection.FindOne(ctx, filter).Decode(&link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *linkRepo) FindLinkByShortID(ctx context.Context, id string) (*entity.Link, error) {
	var link entity.Link
	filter := bson.M{
		"short_id":   id,
		"deleted_at": bson.M{"$exists": false},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *linkRepo) FindAllLinkbyTenantID(ctx context.Context, tenantID string) ([]*entity.Link, error) {
	filter := bson.M{"tenant_id": tenantID, "deleted_at": bson.M{"$exists": false}}

	opts := options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNoEvents
	} else if err != nil {
		return nil, err
	}

	var results []*entity.Link
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}
	return results, nil
}

func (r *linkRepo) SearchLinks() {}
