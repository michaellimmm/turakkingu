package repository

import (
	"context"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LinkRepo interface {
	CreateLink(context.Context, *entity.Link) error
	FindLinkByID(context.Context, string) (*entity.Link, error)
	FindLinkByShortID(context.Context, string) (*entity.Link, error)
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
		"_id":       oid,
		"deletedAt": bson.M{"$exists": false},
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
		"short_id":  id,
		"deletedAt": bson.M{"$exists": false},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}
