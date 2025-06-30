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

type EventName string

const (
	EventNameLandingPage  EventName = "landing_page"
	EventNameThankYouPage EventName = "thank_you_page"
)

type Event struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	TrackID     bson.ObjectID `bson:"track_id"`
	UserAgent   string        `bson:"user_agent"`
	Fingerprint string        `bson:"fingerprint"`
	Url         string        `bson:"url"`
	EventName   EventName     `bson:"event_name"`
	PublishedAt time.Time     `bson:"published_at"`
	BaseEntity  `bson:",inline"`
}

func (t *Event) SetCreatedAt() {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
}

func (t *Event) SetUpdatedAt() {
	t.UpdatedAt = time.Now().UTC()
}
