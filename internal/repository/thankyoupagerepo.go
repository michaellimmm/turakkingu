package repository

import (
	"context"
	"errors"
	"fmt"
	"github/michaellimmm/turakkingu/internal/entity"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ThankYouPageRepo interface {
	CreatePage(context.Context, *entity.ThankYouPage) error
	UpdatePageFieldsAndReturn(context.Context, bson.ObjectID, *entity.ThankYouPage) (*entity.ThankYouPage, error)
}

type thankYouPageRepo struct {
	collection *mongo.Collection
}

func NewThankYouPageRepo(db *mongo.Database) ThankYouPageRepo {
	collection := db.Collection("thank_you_page")
	return &thankYouPageRepo{
		collection: collection,
	}
}

func (r *thankYouPageRepo) CreatePage(ctx context.Context, page *entity.ThankYouPage) error {
	now := time.Now()
	page.BaseEntity.CreatedAt = now
	page.BaseEntity.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, page)
	if err != nil {
		return fmt.Errorf("failed to create thank you page: %w", err)
	}

	page.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *thankYouPageRepo) UpdatePageFieldsAndReturn(ctx context.Context, id bson.ObjectID,
	page *entity.ThankYouPage) (*entity.ThankYouPage, error) {
	filter := bson.M{"_id": id}

	updateDoc := bson.M{
		"updated_at": time.Now(),
	}

	// Only update non-zero/non-empty fields
	if page.URL != "" {
		updateDoc["url"] = page.URL
	}
	if page.Point != 0 {
		updateDoc["point"] = page.Point
	}
	if page.Name != "" {
		updateDoc["name"] = page.Name
	}
	if page.Status != 0 {
		updateDoc["status"] = page.Status
	}

	update := bson.M{"$set": updateDoc}

	var updatedPage entity.ThankYouPage
	err := r.collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedPage)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("thank you page not found")
		}
		return nil, fmt.Errorf("failed to update and return thank you page: %w", err)

	}

	return &updatedPage, nil
}
