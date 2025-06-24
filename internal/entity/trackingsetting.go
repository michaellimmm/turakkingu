package entity

import "go.mongodb.org/mongo-driver/v2/bson"

type TrackingSetting struct {
	ID             bson.ObjectID
	ThanksYouPages []ThankYouPage `bson:"thank_you_pages"`
	BaseEntity     `bson:",inline"`
}

type ThankYouPage struct {
	ID         bson.ObjectID
	URL        string `bson:"url"`
	Point      int    `bson:"point"`
	Name       string `bson:"name"`
	Status     string
	BaseEntity `bson:",inline"`
}

type TrackingStatus int

const (
	TrackingStatusPending TrackingStatus = iota
	TrackingStatusCollected
)

func (ts TrackingStatus) String() string {
	switch ts {
	case TrackingStatusPending:
		return "Pending"
	case TrackingStatusCollected:
		return "Collected"
	default:
		return "Unknown"
	}
}
