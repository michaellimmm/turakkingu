package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventName string

const (
	EventNameLandingPage  EventName = "landing_page"
	EventNameThankYouPage EventName = "thank_you_page"
)

type Event struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	TrackID     string        `bson:"track_id"`
	UserAgent   string        `bson:"user_agent"`
	Fingerprint string        `bson:"fingerprint"`
	Url         string        `bson:"url"`
	EventName   EventName     `bson:"event_name"`
	PublishedAt time.Time     `bson:"published_at"`
	BaseEntity  `bson:",inline"`
}

func (t *Event) GetTrackID() (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(t.TrackID)
}

func (t *Event) SetCreatedAt() {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
}

func (t *Event) SetUpdatedAt() {
	t.UpdatedAt = time.Now().UTC()
}
