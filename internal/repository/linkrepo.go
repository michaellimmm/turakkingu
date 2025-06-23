package repository

import (
	"context"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LinkRepo interface {
	Create(context.Context, *entity.Link) error
	FindByID(context.Context, string) (*entity.Link, error)
	FindByShortID(context.Context, string) (*entity.Link, error)
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

func (l *linkRepo) Create(ctx context.Context, link *entity.Link) error {
	err := link.SetShortID()
	if err != nil {
		return fmt.Errorf("failed to set short id")
	}

	link.SetCreatedAt()
	link.SetUpdatedAt()

	res, err := l.collection.InsertOne(ctx, link)
	if err != nil {
		return err
	}
	link.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (l *linkRepo) FindByID(ctx context.Context, id string) (*entity.Link, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("id is not valid")
	}

	var link entity.Link
	filter := bson.M{
		"_id":       oid,
		"deletedAt": bson.M{"$exists": false},
	}

	err = l.collection.FindOne(ctx, filter).Decode(&link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (l *linkRepo) FindByShortID(ctx context.Context, id string) (*entity.Link, error) {
	var link entity.Link
	filter := bson.M{
		"short_id":  id,
		"deletedAt": bson.M{"$exists": false},
	}

	err := l.collection.FindOne(ctx, filter).Decode(&link)
	if err != nil {
		return nil, err
	}
	return &link, nil
}
