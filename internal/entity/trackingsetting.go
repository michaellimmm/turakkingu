package entity

import "go.mongodb.org/mongo-driver/v2/bson"

// add migration
// db.thankyoupages.createIndex({"tracking_setting_id": 1})
// db.thankyoupages.createIndex({"tenant_id": 1, "status": 1})

type TrackingSetting struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	TenantID   string        `bson:"tenant_id"`
	BaseEntity `bson:",inline"`
}

type ThankYouPage struct {
	ID                bson.ObjectID  `bson:"_id,omitempty"`
	TrackingSettingID bson.ObjectID  `bson:"tracking_setting_id" json:"tracking_setting_id"`
	URL               string         `bson:"url" json:"url"`
	Point             int            `bson:"point" json:"point"`
	Name              string         `bson:"name" json:"name"`
	Status            TrackingStatus `bson:"tracking_status" json:"tracking_status"`
	BaseEntity        `bson:",inline"`
}

type TrackingStatus int

const (
	TrackingStatusUnknown TrackingStatus = iota
	TrackingStatusPending
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

type TrackingSettingWithPages struct {
	ID            bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	TenantID      string         `bson:"tenant_id"`
	ThankYouPages []ThankYouPage `bson:"thank_you_pages"`
	BaseEntity    `bson:",inline"`
}
