package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Track struct {
	ID                bson.ObjectID     `bson:"_id,omitempty" json:"id"`                        // id
	TrackingSettingID bson.ObjectID     `bson:"tracking_setting_id" json:"tracking_setting_id"` // put tracking setting id
	Url               string            `bson:"url" json:"url"`                                 // LP page url
	EndUserID         string            `bson:"end_user_id" json:"end_user_id"`
	SessionID         string            `bson:"session_id" json:"session_id"`
	Platform          string            `bson:"platform" json:"platform"`
	GeneratedFrom     string            `bson:"generated_from" json:"generated_from"`
	Metadata          map[string]string `bson:"metadata" json:"metadata"`
	BaseEntity        `bson:",inline"`
}

func (t *Track) SetCreatedAt() {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
}

func (t *Track) SetUpdatedAt() {
	t.UpdatedAt = time.Now().UTC()
}

type TrackWithThankYouPages struct {
	Track           `bson:",inline"`
	ThankYouPages   []*ThankYouPage  `bson:"thank_you_pages" json:"thank_you_pages"`
	TrackingSetting *TrackingSetting `bson:"tracking_setting,omitempty" json:"tracking_setting,omitempty"`
}
