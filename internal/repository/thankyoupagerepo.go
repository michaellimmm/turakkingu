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
	collection          *mongo.Collection
	trackingSettingRepo TrackingSettingRepo
}

func NewThankYouPageRepo(db *mongo.Database, trackingSettingRepo TrackingSettingRepo) ThankYouPageRepo {
	collection := db.Collection("thank_you_page")
	return &thankYouPageRepo{
		collection:          collection,
		trackingSettingRepo: trackingSettingRepo,
	}
}

func (r *thankYouPageRepo) CreatePage(ctx context.Context, page *entity.ThankYouPage) error {
	now := time.Now()
	page.CreatedAt = now
	page.UpdatedAt = now

	err := r.validateTrackingSettingID(ctx, page.TrackingSettingID)
	if err != nil {
		return err
	}

	isExist, err := r.existByUrls(ctx, page.TrackingSettingID, page.URL)
	if err != nil {
		return fmt.Errorf("failed to validate url: %w", err)
	}

	if isExist {
		return fmt.Errorf("url %s does exist", page.URL)
	}

	result, err := r.collection.InsertOne(ctx, page)
	if err != nil {
		return fmt.Errorf("failed to create thank you page: %w", err)
	}

	page.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *thankYouPageRepo) validateTrackingSettingID(ctx context.Context, trackingSettingID bson.ObjectID) error {
	if trackingSettingID.IsZero() {
		return errors.New("tracking_setting_id is required")
	}

	exists, err := r.trackingSettingRepo.IsTrackingSettingIDExist(ctx, trackingSettingID)
	if err != nil {
		return fmt.Errorf("failed to check tracking setting existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("tracking setting with ID %s does not exist", trackingSettingID.Hex())
	}

	return nil
}

func (r *thankYouPageRepo) existByUrls(ctx context.Context, trackingSettingID bson.ObjectID, url string) (bool, error) {
	filter := bson.M{
		"tracking_setting_id": trackingSettingID,
		"url":                 url,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	return count > 0, err
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
